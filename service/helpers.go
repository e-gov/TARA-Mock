package main

import (
	"crypto/rand"
	"fmt"
)

// PseudoUUID genereerib juhusõne. Vt: https://stackoverflow.com/questions/15130321/is-there-a-method-to-generate-a-uuid-with-go-language.
func PseudoUUID() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
