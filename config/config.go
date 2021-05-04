package config

//Propperties Configurations properties on env
type Properties struct {
	Port           string `env:"APP_PORT" env-default:"8080"`
	Host           string `env:"APP_HOST" env-default:"localhost"`
	DBHost         string `env:"DB_HOST" env-default:"localhost"`
	DBPort         string `env:"DB_PORT" env-default:"27017"`
	DBName         string `env:"DB_NAME" env-default:"tronics"`
	DBUser         string `env:"DB_USER" env-default:"root"`
	DBPass         string `env:"DB_PASS" env-default:"123456"`
	CollectionName string `env:"COLLECTION_NAME" env-default:"products"`
}
