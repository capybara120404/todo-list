package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Addr     string
	PathToDB string
)

func init() {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatalf("an error occurred while importing configuration files: %v", err)
	}
	Addr = fmt.Sprintf(":%s", os.Getenv("TODO_PORT"))
	PathToDB = fmt.Sprintf(":%s", os.Getenv("TODO_DBFILE"))
}
