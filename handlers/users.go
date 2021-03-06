package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Mersock/golang-echo-mongodb-restful-api/config"
	"github.com/Mersock/golang-echo-mongodb-restful-api/dbiface"
	"github.com/dgrijalva/jwt-go"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	cf config.Properties
)

type User struct {
	Email    string `json:"username" bson:"username" validate:"required,email"`
	Password string `json:"password,omitempty" bson:"password" validate:"required,min=8,max=300"`
	IsAdmin  bool   `json:"isAdmin,omitempty" bson:"isAdmin"`
}

type UsersHandler struct {
	Col dbiface.CollectionAPI
}

type errorMessage struct {
	Message string `json:"message"`
}

func insertUser(ctx context.Context, user User, collection dbiface.CollectionAPI) (interface{}, *echo.HTTPError) {
	var newUser User
	res := collection.FindOne(ctx, bson.M{"username": user.Email})
	err := res.Decode(&newUser)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Errorf("Unable to decode retrived user :%v", err)
		return newUser, echo.NewHTTPError(http.StatusBadRequest, errorMessage{Message: "Unable to decode retrived user"})
	}
	if newUser.Email != "" {
		log.Errorf("User by %s already exists", user.Email)
		return newUser, echo.NewHTTPError(http.StatusBadRequest, errorMessage{Message: "User already exists"})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		log.Errorf("Unable to hash the password: %v", err)
		return newUser, echo.NewHTTPError(http.StatusInternalServerError, errorMessage{Message: "Unable to process the password"})
	}
	user.Password = string(hashPassword)
	insertRes, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Errorf("Unable to insert the user :%+v", err)
		return newUser, echo.NewHTTPError(http.StatusBadRequest, errorMessage{Message: "Unable to create the user"})
	}
	return insertRes.InsertedID, nil
}

func (h *UsersHandler) CreateUser(c echo.Context) error {
	var user User
	c.Echo().Validator = &userValidator{validator: v}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Unable to bind user struct : %+v", err)
		return c.JSON(http.StatusBadRequest, errorMessage{Message: "Unable to parse the request payload"})
	}
	if err := c.Validate(user); err != nil {
		log.Errorf("Unable to validate the requested body : %+v", err)
		return c.JSON(http.StatusBadRequest, errorMessage{Message: "Unable to validate request body"})
	}
	insertedUserID, insertErr := insertUser(context.Background(), user, h.Col)
	if insertErr != nil {
		return c.JSON(insertErr.Code, insertErr.Message)
	}
	token, err := user.createToken()
	if err != nil {
		log.Errorf("Unable to generate the token")
		return echo.NewHTTPError(http.StatusInternalServerError, errorMessage{Message: "Unable to generate the token"})
	}
	c.Response().Header().Set("x-auth-token", token)
	return c.JSON(http.StatusCreated, insertedUserID)
}

func isCredValid(givenPwd, storedPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(storedPwd), []byte(givenPwd)); err != nil {
		return false
	}
	return true
}

func authenticateUser(ctx context.Context, reqUser User, collection dbiface.CollectionAPI) (User, *echo.HTTPError) {
	var storedUser User
	res := collection.FindOne(ctx, bson.M{"username": reqUser.Email})
	err := res.Decode(&storedUser)
	if err != nil {
		log.Errorf("Unable to decode retrieved user: %v", err)
		return storedUser, echo.NewHTTPError(http.StatusUnauthorized, errorMessage{Message: "Credentials invalid"})
	}
	if err == mongo.ErrNoDocuments {
		log.Errorf("user %s does not exist", reqUser.Email)
		return storedUser, echo.NewHTTPError(http.StatusUnauthorized, errorMessage{Message: "Credentials invalid"})
	}
	if !isCredValid(reqUser.Password, storedUser.Password) {
		return storedUser, echo.NewHTTPError(http.StatusUnauthorized, errorMessage{Message: "Credentials invalid"})
	}
	return storedUser, nil
}

func (u User) createToken() (string, *echo.HTTPError) {
	if err := cleanenv.ReadEnv(&cf); err != nil {
		log.Fatalf("Configuration cannot be read :%v", err)
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = u.IsAdmin
	claims["user_id"] = u.Email
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(cf.JwtTokenSecret))
	if err != nil {
		log.Errorf("Unable to generate the token :%v", err)
		return "", echo.NewHTTPError(http.StatusInternalServerError, errorMessage{Message: "Unable to generate the token"})
	}
	return token, nil
}

func (h *UsersHandler) AuthUser(c echo.Context) error {
	var user User
	c.Echo().Validator = &userValidator{validator: v}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Unable to bind to user struct: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, errorMessage{Message: "Unable to parse the request payload"})
	}
	if err := c.Validate(user); err != nil {
		log.Errorf("Unable to validate the request body: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, errorMessage{Message: "Unable to validate request body"})
	}
	user, err := authenticateUser(context.Background(), user, h.Col)
	if err != nil {
		log.Errorf("Unable to authenticate to database: %v", err)
		return c.JSON(err.Code, err.Message)
	}
	token, err := user.createToken()
	if err != nil {
		log.Errorf("Unable to genarate the token: %v", err)
		return c.JSON(err.Code, err.Message)
	}
	c.Response().Header().Set("x-auth-token", "Bearer "+token)
	return c.JSON(http.StatusOK, User{Email: user.Email})
}
