package config

import (
	"encoding/json"
	"os"
	"strings"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	var config Config
	configPath, err := getConfigFilePath()
	if err != nil {
		return config, err
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return config, err
	}

	defer configFile.Close()
	buffer := make([]byte, 1024)
	if _, err = configFile.Read(buffer); err != nil {
		return config, err
	}

	cleanData := strings.ReplaceAll(string(buffer), "\x00", "")
	if err = json.Unmarshal([]byte(cleanData), &config); err != nil {
		return config, err
	}

	return config, nil
}

func (c *Config) SetUser(user string) error {
	c.Current_user_name = user
	if err := write(*c); err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/" + configFileName, nil
}

func write(cfg Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	configFile, err := os.Create(configPath)
	if err != nil {
		return err
	}

	defer configFile.Close()
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	_, err = configFile.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
