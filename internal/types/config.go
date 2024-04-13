package types

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Config struct {
	Port             string `mapstructure:"PORT"`
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresPort     string `mapstructure:"POSTGRES_PORT"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	TGToken          string `mapstructure:"TG_TOKEN"`
	JWTSecret        string `mapstructure:"JWT_SECRET"`
}

type ServerContext struct {
	Config Config
	Log    *logrus.Logger
	DB     *gorm.DB
}

func InitConfig() (Config, error) {
	var config Config
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	//viper.AutomaticEnv()
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	fmt.Println(config)
	return config, nil
}
