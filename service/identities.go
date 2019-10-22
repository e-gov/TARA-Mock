package main

import (
	"encoding/json"
	"fmt"
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
		fmt.Printf("TARA-Mock: Ettevalmistatud identiteetide lugemine ebaõnnestus. %s\n", err.Error())
		os.Exit(1)
	}

	defer fh.Close()
	// Dekodeeri JSON-struktuur
	jsonParser := json.NewDecoder(fh)
	err = jsonParser.Decode(&d)
	if err != nil {
		fmt.Println("TARA-Mock: Ettevalmistatud identiteetide  dekodeerimine ebaõnnestus.")
		os.Exit(1)
	}
	fmt.Println("Loetud identiteedid:")
	for _, id := range d {
		fmt.Printf("  %s, %s, %s\n", id.Isikukood, id.Eesnimi, id.Perekonnanimi)
	}
	return d
}
