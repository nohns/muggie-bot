package main

import (
	"fmt"
	"os"
	"strconv"
)

type config struct {
	token   string
	permInt int
}

// Read configuration of bot from environment variables
func readConfFromEnv() (*config, error) {
	token, err := stringEnvVar("DISCORD_BOT_TOKEN")
	if err != nil {
		return nil, err
	}

	permInt, err := intEnvVar("PERMISSION_INTEGER")
	if err != nil {
		return nil, err
	}

	return &config{
		token:   token,
		permInt: permInt,
	}, nil
}

// Read environment variable as string. Error returned when var is not found
func stringEnvVar(envname string) (string, error) {
	val, ok := os.LookupEnv(envname)
	if !ok {
		return "", fmt.Errorf("missing env var '%s' string value", envname)
	}

	return val, nil
}

// Read environment variable as int. Error returned when var could not be parsed as int
func intEnvVar(envname string) (int, error) {
	val, err := strconv.Atoi(os.Getenv(envname))
	if err != nil {
		return 0, fmt.Errorf("could not read env var '%s' int value: %v", envname, err)
	}

	return val, nil
}
