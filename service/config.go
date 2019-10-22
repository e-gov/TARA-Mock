package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config on TARA-Mock seadistuse tüüp.
type Config struct {
	TaraMockHost   string `json:"taraMockHost"`
	HTTPServerPort string `json:"httpServerPort"`
	// TARA-Mock-i HTTPS sert
	TaraMockCert string `json:"taraMockCert"`
	// TARA-Mock-i HTTPS privaatvõti
	TaraMockKey string `json:"taraMockKey"`
	// TARA-Mock-i identsustõendi allkirjastamise avalik võti
	IDTokenPrivKeyPath string `json:"idTokenPrivKeyPath"`
	// TARA-Mock-i identsustõendi allkirjastamise privaatvõti
	IDTokenPubKeyPath string `json:"idTokenPubKeyPath"`
	// Identsustõendi allkirjavõtme identifikaator
	Kid string `json:"kid"`
}

// TARA-Mock sisseloetud seadistus.
var conf Config

// loadConf loeb JSON-failist f sisse seadistuse.
func loadConf(f string) Config {
	var config Config
	configFile, err := os.Open(f)
	defer configFile.Close()
	if err != nil {
		fmt.Printf("TARA-Mock: Seadistuse lugemine ebaõnnestus. %s\n", err.Error())
		os.Exit(1)
	}
	// Dekodeeri JSON-struktuur
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		fmt.Println("TARA-Mock: Seadistuse dekodeerimine ebaõnnestus.")
		os.Exit(1)
	}
	return config
}