package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/HectorMRC/gw-pool/pool"
	"github.com/joho/godotenv"

	// required by postgres connections
	_ "github.com/lib/pq"
)

const (
	infoStart = "The server is being started on %s%s"
	infoDone  = "The service has finished successfully"

	errServicePort  = "Service port must be set."
	errServiceNetw  = "Service network must be set."
	errDotenvConfig = "Service has failed setting up dotenv: %s"
	errListenFailed = "Service has failed listening: %s"
	errServeFailed  = "Service has failed serving: %s"
	errPostgresDNS  = "Database DNS must be set."

	envPortKey  = "SERVICE_PORT"
	envNetwKey  = "SERVICE_NETW"
	envSleepKey = "SLEEP_SEC"
	envDNSKey   = "DATABASE_DNS"
)

// Single instance of a Pool
var datapool pool.Pool

func requestHandler(w http.ResponseWriter, r *http.Request) {
	var coord pool.Coordinates
	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&coord); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Got a location update from driver %v", coord.GetDriverID())
	datapool.Insert(&coord)
}

func main() {
	// setting up environment
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		log.Fatalf(errDotenvConfig, err.Error())
	}

	servicePort, exists := os.LookupEnv(envPortKey)
	if !exists {
		log.Panicf(errServicePort)
	}

	serviceNetw, exists := os.LookupEnv(envNetwKey)
	if !exists {
		log.Panicf(errServiceNetw)
	}

	dns, exists := os.LookupEnv(envDNSKey)
	if !exists {
		log.Panicf(errPostgresDNS)
	}

	sleepEnv := os.Getenv(envSleepKey)
	sleep := time.Second
	if secs, _ := strconv.Atoi(sleepEnv); secs > 0 {
		sleep = time.Duration(secs) * time.Second
	}

	// initializing data pool
	open := func() (pool.Conn, error) { return sql.Open("postgres", dns) }
	datapool = pool.NewDatapool(open, sleep)
	datapool.Reset()
	defer datapool.Stop()

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

	// on finishing
	log.Print(infoDone)
}
