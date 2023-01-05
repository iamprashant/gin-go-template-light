package configs

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type BasicAuth struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type LogConfig struct {
}

type PostgresConfig struct {
	Host               string    `mapstructure:"host" validate:"required"`
	Port               int       `mapstructure:"port"`
	Auth               BasicAuth `mapstructure:"auth"`
	DBName             string    `mapstructure:"db_name" validate:"required"`
	MaxIdealConnection int       `mapstructure:"max_ideal_connection" validate:"required"`
	MaxOpenConnection  int       `mapstructure:"max_open_connection" validate:"required"`
	SslMode            string    `mapstructure:"ssl_mode" validate:"required"`
}

// Application config structure
type AppConfig struct {
	Name           string         `mapstructure:"name" validate:"required"`
	LogLevel       string         `mapstructure:"log_level" validate:"required"`
	PostgresConfig PostgresConfig `mapstructure:"postgres" validate:"required"`
	Log            LogConfig      `mapstructure:"log"`
}

// reading config and intializing configs for application
func InitConfig() (*viper.Viper, error) {
	vConfig := viper.NewWithOptions(viper.KeyDelimiter("__"))

	vConfig.AddConfigPath(".")
	vConfig.SetConfigName(".env")
	vConfig.SetConfigType("env")
	vConfig.AutomaticEnv()
	vConfig.ReadInConfig()
	//
	setDefault(vConfig)
	if err := vConfig.ReadInConfig(); err != nil && !os.IsNotExist(err) {
		log.Printf("Reading from env varaibles.")
	}

	return vConfig, nil
}

func setDefault(v *viper.Viper) {
	// setting all default values
	// keeping watch on https://github.com/spf13/viper/issues/188

	v.SetDefault("NAME", "slack-app")
	v.SetDefault("LOG_LEVEL", "debug")

	v.SetDefault("POSTGRES__HOST", "localhost")
	v.SetDefault("POSTGRES__PORT", 5432)
	v.SetDefault("POSTGRES__DB_NAME", "<>")
	v.SetDefault("POSTGRES__AUTH__USER", "<>")
	v.SetDefault("POSTGRES__AUTH__PASSWORD", "<>")
	v.SetDefault("POSTGRES__MAX_OPEN_CONNECTION", 10)
	v.SetDefault("POSTGRES__MAX_IDEAL_CONNECTION", 10)
	v.SetDefault("POSTGRES__SSL_MODE", "disable")
}

// Getting application config from viper
func GetApplicationConfig(v *viper.Viper) (*AppConfig, error) {
	var config AppConfig
	err := v.Unmarshal(&config)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	// valdating the app config
	validate := validator.New()
	err = validate.Struct(&config)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return &config, nil
}
