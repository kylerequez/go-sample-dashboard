package utils

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	log.Println("::: Loading Env Variables...")
	if err := godotenv.Load(); err != nil {
		return err
	}
	log.Println("::: Successfully Loaded Env Variables")
	return nil
}

func GetEnv(name string) (*string, error) {
	env, isExists := os.LookupEnv(name)
	if !isExists {
		return nil, errors.New("env variable does not exists")
	}
	if env == "" {
		return nil, errors.New("env variable is empty")
	}

	return &env, nil
}
