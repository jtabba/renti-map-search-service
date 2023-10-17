package envHelper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVar(key string) string {
	err := godotenv.Load(".env")

	if (err != nil) {
		log.Fatalf("Error while reading config file %s", err)
	}

	return os.Getenv(key)
}