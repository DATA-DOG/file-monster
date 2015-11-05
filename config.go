package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config main config file structure
type Config struct {
	Port           int    `yaml:"port"`
	ClientID       string `yaml:"client_id"`
	ClientSecret   string `yaml:"client_secret"`
	DbUsername     string `yaml:"db_username"`
	DbPassword     string `yaml:"db_password"`
	DbDatabase     string `yaml:"db_database"`
	SlackWebookURL string `yaml:"slack_webhook_url"`
}

var config Config

func loadConfig() {
	if len(os.Args) < 2 {
		log.Fatal("Missing required config parameter")
	}

	contents, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(contents, &config)

	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	log.Print("Config successfully loaded")
}
