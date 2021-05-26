package config

//Propperties Configurations properties on env
type Properties struct {
	Port              string `env:"APP_PORT" env-default:"8080"`
	Host              string `env:"APP_HOST" env-default:"localhost"`
	DBHost            string `env:"DB_HOST" env-default:"localhost"`
	DBPort            string `env:"DB_PORT" env-default:"27017"`
	DBName            string `env:"DB_NAME" env-default:"tronics"`
	DBTestName        string `env:"DB_TEST_NAME" env-default:"tronics_test"`
	DBUser            string `env:"DB_USER" env-default:"root"`
	DBPass            string `env:"DB_PASS" env-default:"123456"`
	ProductCollection string `env:"PRODUCTS_COL_NAME" env-default:"products"`
	UserCollection    string `env:"USERS_COL_NAME" env-default:"users"`
	JwtTokenSecret    string `env:"JWT_TOKEN_SECRET" env-default:"27rmzGBbdlNCpHltxogK69DSqjdhPttQSBkm_F3SvJG4XTes6Kjts5LhENBZDVhtSRda_FJKwcRVR0V4iSS30t2cDvD5ZF_ZIK4aLjOOP8jPCZjmamhidAhHP8MIZ-0MtRlDf4HVqRDWTRXoIoXKDBEC2VntRC7nhKL1OqfplZk"`
}
