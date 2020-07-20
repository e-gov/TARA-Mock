package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"os"
)

// Identity on autenditav identiteet (isik).
type Identity struct {
	Isikukood     string `json:"isikukood"`
	Eesnimi       string `json:"eesnimi"`
	Perekonnanimi string `json:"perekonnanimi"`
}

// identities hoiab kasutajale valimiseks pakutavaid identiteete (isikuid).
var identities []Identity

// loadIdentities loeb JSON-failist f sisse ettevalmistatud identiteedid.
func loadIdentities(f string) []Identity {
	var d []Identity // Dekodeeritud identiteedid

	// Ava fail
	fh, err := os.Open(f) // File handle
	if err != nil {
		log.Errorf("TARA-Mock: Ettevalmistatud identiteetide lugemine ebaõnnestus. %v", err)
		os.Exit(1)
	}

	defer fh.Close()
	// Dekodeeri JSON-struktuur
	jsonParser := json.NewDecoder(fh)
	err = jsonParser.Decode(&d)
	if err != nil {
		log.Error("TARA-Mock: Ettevalmistatud identiteetide  dekodeerimine ebaõnnestus.")
		os.Exit(1)
	}
	log.Info("Loaded identities/Loetud identiteedid:")
	for _, id := range d {
		log.Infof("	%s, %s, %s", id.Isikukood, id.Eesnimi, id.Perekonnanimi)
	}
	return d
}
