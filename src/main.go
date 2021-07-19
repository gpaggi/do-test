package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	S "echoapi/server"
)

func main() {
	// Basic configuration to the logger to output in JSON format
	log.SetFormatter(&log.JSONFormatter{})

	var (
		listen       string
		htpasswdPath string
	)

	// Get the bind address for the webserver from the environment
	la, ok := os.LookupEnv("LISTEN_ADDR")
	if !ok {
		listen = ":9090"
	} else {
		listen = la
	}

	// Get the log level from the environment and configure the logger accordingly
	lvl := os.Getenv("LOG_LEVEL")
	if lvl != "" {
		lvl, err := log.ParseLevel(lvl)
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(lvl)
	}

	// Get the path to the htpasswd file from the environment
	ht, ok := os.LookupEnv("HTPASSWD_PATH")
	if !ok {
		htpasswdPath = "./conf/htpasswd"
	} else {
		htpasswdPath = ht
	}

	// Start the webserver or fail hard
	log.Fatal(S.Start(listen, htpasswdPath))
}
