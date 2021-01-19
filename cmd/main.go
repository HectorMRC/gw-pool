package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/HectorMRC/gw-pool/pool"

	"github.com/HectorMRC/gw-pool/location"
	"github.com/joho/godotenv"
)

const (
	infoStart = "The server is being started on %s%s"
	infoDone  = "The service has finished successfully"

	errServicePort  = "Service port must be set."
	errServiceNetw  = "Service network must be set."
	errDotenvConfig = "Service has failed setting up dotenv: %s"
	errListenFailed = "Service has failed listening: %s"
	errServeFailed  = "Service has failed serving: %s"

	envPortKey = "SERVICE_PORT"
	envNetwKey = "SERVICE_NETW"
)

// Single instance of a Gateway
var gwpool = pool.NewDatapool()

func requestHandler(w http.ResponseWriter, r *http.Request) {
	var coord location.Coordinates

	err := json.NewDecoder(r.Body).Decode(&coord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != gwpool.Insert(&coord) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func main() {
	// setting up environment
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		log.Fatalf(err.Error())
	}

	servicePort, exists := os.LookupEnv(envPortKey)
	if !exists {
		log.Panicf(errServicePort)
	}

	serviceNetw, exists := os.LookupEnv(envNetwKey)
	if !exists {
		log.Panicf(errServiceNetw)
	}

	// starting http service
	address := ":" + servicePort
	log.Printf(infoStart, serviceNetw, address)

	lis, err := net.Listen(serviceNetw, address)
	if err != nil {
		log.Panicf(errListenFailed, err)
	}

	defer lis.Close()
	http.HandleFunc("/", requestHandler)
	if err := http.Serve(lis, nil); err != nil {
		log.Panicf(errServeFailed, err)
	}

	// once finishing
	log.Print(infoDone)
}
