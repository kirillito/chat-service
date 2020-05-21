package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ChatMessage struct {
	msg string
	id  int
}

type User struct {
	id   int
	name string
}

type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	ConnType string `json:"connType"`
	LogPath  string `json:"logFilePath"`
}

// Global vars
var state = []User{}
var userNameList = map[string]int{} // map userNames to user ID
var msgHistory []string             // keep last 30 messages in mem

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg := configureService()

	// Listen for incoming tcp connections ie: telnet
	ln, err := net.Listen(cfg.ConnType, cfg.Host+":"+cfg.Port)

	if err != nil {
		log.Fatal().Err(err).Msg("Error listening")
		os.Exit(1)
	}

	defer ln.Close()

	// handle new incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal().Err(err).Msg("Error handling connection")
		}

		go handleConnection(conn, cfg)
	}
}

// read configuration from file
func configureService() Config {
	byteVal, _ := ioutil.ReadFile("./config.json")

	var cfg Config
	if err := json.Unmarshal(byteVal, &cfg); err != nil {
		log.Fatal().
			Err(err).
			Str("service", "chat-service").
			Msgf("Cannot start %s", "chat-service")

		panic(err)
	}
	return cfg
}

func handleConnection(conn net.Conn, cfg Config) {

}
