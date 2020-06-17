/*
Package chat-service is a simple chat service
*/
package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ChatMessage is a structure for storing message data
type ChatMessage struct {
	msg string
	id  int
}

// User is a structure for storing user data
type User struct {
	id   int
	name string
}

// Config is a structure for storing service configuration data
type Configuration struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	ConnType string `json:"connType"`
	LogPath  string `json:"logFilePath"`
}

type HttpAPI struct {
	broker				*Broker
	configuration	Configuration
}

var state = []User{}								// Array of current users in the chat
var userNameList = map[string]int{} // Map userNames to user ID
var msgHistory []string             // Keep last 30 messages in memory

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	cfg := configureService()

	// system user used for notifications
	createNewUser("system")

	// create and start our primary chatroom
	b := CreateMessageBroker()
	go b.Start()

	// Listen for incoming tcp connections, i.e. telnet
	ln, err := net.Listen(cfg.ConnType, cfg.Host+":"+cfg.Port)

	if err != nil {
		log.Fatal().Err(err).Msg("Error listening")
		os.Exit(1)
	}

	defer ln.Close();

	// handle new incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal().Err(err).Msg("Error handling connection")
		}

		userId := createNewUser("")

		go handleConnection(conn, userId, cfg)
	}
}

// read configuration from file
func configureService() Configuration {
	byteVal, _ := ioutil.ReadFile("./config.json")

	var cfg Configuration
	if err := json.Unmarshal(byteVal, &cfg); err != nil {
		log.Fatal().
			Err(err).
			Str("service", "chat-service").
			Msgf("Cannot start %s", "chat-service")

		panic(err)
	}
	return cfg
}

func handleConnection(conn net.Conn, userId int, cfg Configuration) {
	remoteAddress := conn.RemoteAddr().String()
	userName = getUserNameById(userId)
	fmt.Printf("Client id: %d with username [%s] connected from %s\n", userId, , remoteAddress)

	msg : fmt.Sprintf("%s has joined", userName)
	
}

// createNewUser creates new user and returns its id
func createNewUser(userName string) int {
	newUser := User {
		name: userName, 
		id: len(state) 
	}
	state = append(state, newUser)

	userNameList[newUser.name] = newUser.id

	return newUser.id
}

// look up username
func getUserNameById(id int) string {
	userName := state[id].name

	if userName != "" {
		return userName
	}

	userName = "Anonymous" + strconv.Itoa(id)
	return userName
}
