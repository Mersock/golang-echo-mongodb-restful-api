package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Mersock/golang-echo-mongodb-restful-api/config"
	"github.com/Mersock/golang-echo-mongodb-restful-api/dbiface"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
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
}

type UsersHandler struct {
	Col dbiface.CollectionAPI
}

type userValidator struct {
	validator *validator.Validate
}

func (u *userValidator) Validate(i interface{}) error {
	return u.validator.Struct(i)
}

func insertUser(ctx context.Context, user User, collection dbiface.CollectionAPI) (interface{}, *echo.HTTPError) {
	var newUser User
	res := collection.FindOne(ctx, bson.M{"username": user.Email})
	err := res.Decode(&newUser)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Errorf("Unable to decode retrived user :%v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Unable to decode retrived user")
	}
	if newUser.Email != "" {
		log.Errorf("User by %s already exists", user.Email)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "User already exists")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		log.Errorf("Unable to hash the password: %v", err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Unable to process the password")
	}
	user.Password = string(hashPassword)
	insertRes, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Errorf("Unable to insert the user :%+v", err)
		return nil, echo.NewHTTPError(http.StatusBadRequest, "Unable to create the user")
	}
	return insertRes.InsertedID, nil
}

func (h *UsersHandler) CreateUser(c echo.Context) error {
	var user User
	c.Echo().Validator = &userValidator{validator: v}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Unable to bind user struct : %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to parse the request payload")
	}
	if err := c.Validate(user); err != nil {
		log.Errorf("Unable to validate the requested body : %+v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to validate request payload.")
	}
	insertedUserID, err := insertUser(context.Background(), user, h.Col)
	if err != nil {
		return err
	}
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
		return storedUser, echo.NewHTTPError(http.StatusBadRequest, "Unable to decode retreved user")
	}
	if err == mongo.ErrNoDocuments {
		log.Errorf("user %s does not exist", reqUser.Email)
		return storedUser, echo.NewHTTPError(http.StatusNotFound, "User does not exist")
	}
	if !isCredValid(reqUser.Password, storedUser.Password) {
		return storedUser, echo.NewHTTPError(http.StatusUnauthorized, "Credentials invalid")
	}
	return User{Email: storedUser.Email}, nil
}

func createToken(username string) (string, *echo.HTTPError) {
	if err := cleanenv.ReadEnv(&cf); err != nil {
		log.Fatalf("Configuration cannot be read :%v", err)
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = username
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(cf.JwtTokenSecret))
	if err != nil {
		log.Errorf("Unable to generate the token :%v", err)
		return "", echo.NewHTTPError(http.StatusInternalServerError, "Unable to generate the token")
	}
	return token, nil
}

func (h *UsersHandler) AuthUser(c echo.Context) error {
	var user User
	c.Echo().Validator = &userValidator{validator: v}
	if err := c.Bind(&user); err != nil {
		log.Errorf("Unable to bind to user struct: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to parse the request payload")
	}
	if err := c.Validate(user); err != nil {
		log.Errorf("Unable to validate the request body: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to validate request payload")
	}
	user, err := authenticateUser(context.Background(), user, h.Col)
	if err != nil {
		log.Errorf("Unable to authenticate to database")
		return err
	}
	token, err := createToken(user.Email)
	if err != nil {
		log.Errorf("Unable to genarate the token")
		return err
	}
	c.Response().Header().Set("x-auth-token", token)
	return c.JSON(http.StatusOK, User{Email: user.Email})
}
