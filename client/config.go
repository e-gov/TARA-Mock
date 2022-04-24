package main

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

// Config on TARA-Mock kliendi seadistuse tüüp.
type Config struct {
	// klientrakenduse hostinimi
	AppHost string `json:"appHost"`
	// klientrakenduse HTTPS serveri port
	AppPort string `json:"appPort"`
	// klientrakenduse HTTPS sert.
	AppCert string `json:"appCert"`
	// klientrakenduse HTTPS privaatvõti.
	AppKey string `json:"appKey"`

	// Usaldusankur TARA-Mock-i poole pöördumisel
	RootCAFile string `json:"rootCAFile"`

	// TARA-Mock-i otspunktid
	TaraMockAuthorizeEndpoint string `json:"taraMockAuthorizeEndpoint"`
	TaraMockTokenEndpoint     string `json:"taraMockTokenEndpoint"`
	TaraMockKeyEndpoint       string `json:"taraMockKeyEndpoint"`

	// OpenID Connect kohane tagasisuunamis-URL
	RedirectURI string `json:"redirectURI"`
}

// TARA-Mock kliendi sisseloetud seadistus.
var conf Config

// loadConf loeb JSON-failist f sisse seadistuse.
func loadConf(f string) Config {
	var config Config
	configFile, err := os.Open(f)
	if err != nil {
		log.Printf("TARA-Mock klient: Seadistuse lugemine ebaõnnestus. %s\n", err.Error())
	}
	defer configFile.Close()
	if err != nil {
		log.Printf("TARA-Mock klient: Seadistuse lugemine ebaõnnestus. %s\n", err.Error())
		os.Exit(1)
	}
	// Dekodeeri JSON-struktuur
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Println("TARA-Mock klient: Seadistuse dekodeerimine ebaõnnestus.")
		os.Exit(1)
	}
	return config
}
