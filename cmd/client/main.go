package main

import (
	"log"
	"os"
	"strconv"

	"github.com/HectorMRC/gw-pool/client"
	"github.com/joho/godotenv"
)

const (
	infoDone        = "Got response 200 - Ok "
	errDotenvConfig = "Client has failed setting up dotenv: %s"
	errClientFailed = "Client request got error %s"
	envHostKey      = "SERVICE_HOST"
	envPortKey      = "SERVICE_PORT"
	envDotKey       = "DOTENV_PATH"
)

var (
	// dotenv file path, if empty dotenv disabled
	dotenvPath = os.Getenv(envDotKey)
	defValues  = []int{123, 123, 1}
)

func main() {
	// setting up environment
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(dotenvPath) > 0 {
		if err := godotenv.Load(); err != nil {
			log.Fatalf(errDotenvConfig, err.Error())
		}
	}

	serverHost, exists := os.LookupEnv(envHostKey)
	if !exists {
		serverHost = "localhost"
	}

	serverPort, exists := os.LookupEnv(envPortKey)
	if !exists {
		serverPort = "8080"
	}

	client.SetURL("http://" + serverHost + ":" + serverPort)

	args := os.Args[1:]
	for index, value := range args {
		if index >= len(defValues) {
			break
		}

		if custom, err := strconv.Atoi(value); err == nil {
			defValues[index] = custom
		}
	}

	if err := client.PostRequest(
		defValues[0],
		defValues[1],
		defValues[2]); err != nil {
		log.Printf(errClientFailed, err.Error())
		os.Exit(1)
	}

	log.Println(infoDone)
}
