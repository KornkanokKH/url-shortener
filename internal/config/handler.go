package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	environment = "DMPENV"
)

var (
	consoleLog zerolog.Logger
)

func ReadConfigFile(configModel interface{}) error {
	function := "internal.config.ReadFile()"

	env := os.Getenv(environment)

	consoleLog.Info().Msgf("Loading configs env : %v.\n", env)
	currentDir, _ := os.Getwd()
	rViper := viper.New()
	rViper.SetConfigType("yml")
	rViper.AddConfigPath(currentDir)
	rViper.AddConfigPath(currentDir + "/configs")
	rViper.AddConfigPath(currentDir + "/../configs")
	rViper.SetConfigName("config") //configs file name

	//Find and Read configs file
	if err := rViper.ReadInConfig(); err != nil {
		return fmt.Errorf("%v: cannot read configuration file, %v", function, err)
	}

	var allEnvConfig map[string]interface{}
	if err := rViper.Unmarshal(&allEnvConfig); err != nil {
		return fmt.Errorf("%v: cannot unmarshal configuration, %v", function, err)
	}

	// get configuration by environment(local, dev, or etc.) and marshal to binary
	envConfig, err := yaml.Marshal(allEnvConfig[env])
	if err != nil {
		return fmt.Errorf("%s: unable to marshal configuration data at env=%s, %s", function, env, err.Error())
	}

	err = yaml.Unmarshal(envConfig, configModel)
	if err != nil {
		return fmt.Errorf("%s: unable to unmarshal configuration data to Config model, %s", function, err.Error())
	}
	return nil
}
