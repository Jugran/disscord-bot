package models

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Token       string
	Port        uint64
	DefaultRole string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func NewConfig() Config {

	port, err := strconv.ParseUint(os.Getenv("PORT"), 10, 16)

	if err != nil {
		log.Fatal("Error: Invalid PORT value")
	}

	config := Config{
		Token:       os.Getenv("TOKEN"),
		DefaultRole: os.Getenv("DEFAULT_ROLE"),
		Port:        port,
	}

	return config
}
