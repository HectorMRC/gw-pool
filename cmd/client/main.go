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
	envURLKey       = "SERVICE_URL"
)

var (
	serverURL = os.Getenv(envURLKey)
	defValues = []int{123, 123, 1}
)

func main() {
	// setting up environment
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		log.Fatalf(errDotenvConfig, err.Error())
	}

	if len(serverURL) > 0 {
		client.SetURL(serverURL)
	}

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
