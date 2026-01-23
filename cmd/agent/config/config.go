package config

import (
	"os"

	"gitlab.com/MikeTTh/env"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"github.com/creasty/defaults"
)

const ConfigDefaultPath = "/etc/gg_agent/config.yaml"

func LoadConfig(logger *zap.Logger, pathOverride string) (AgentConfig, error) {
	var err error

	// open file
	configPath := ConfigDefaultPath
	if pathOverride != "" {
		configPath = pathOverride
	}
	configPath = env.String("CONFIG", configPath)
	logger.Info("Loading config", zap.String("configPath", configPath))

	var configFile *os.File
	configFile, err = os.Open(configPath) // #nosec G304
	if err != nil {
		logger.Error("Error while opening config file", zap.Error(err))
		return AgentConfig{}, err
	}

	defer func(configFile *os.File) {
		e := configFile.Close()
		if e != nil {
			logger.Warn("Could not close config file", zap.Error(e))
		}
	}(configFile)

	// parse it
	var newConfig AgentConfig

	// read defaults
	err = defaults.Set(&newConfig)
	if err != nil {
		return AgentConfig{}, err
	}

	// decode file
	decoder := yaml.NewDecoder(configFile)
	decoder.KnownFields(true)
	err = decoder.Decode(&newConfig)
	if err != nil {
		return AgentConfig{}, err
	}

	// Validate config
	err = newConfig.Validate()
	if err != nil {
		logger.Error("Config file failed validation!", zap.Error(err))
		return AgentConfig{}, err
	}

	logger.Debug("Config successfully loaded", zap.Any("config", newConfig))
	// very good
	return newConfig, nil
}
