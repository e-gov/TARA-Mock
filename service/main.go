package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	// Identsustõendi koostamiseks ja allkirjastamiseks
	// Dok-n: https://godoc.org/github.com/dgrijalva/jwt-go
	"github.com/dgrijalva/jwt-go"
)

const (
	taraMockHost   = "localhost"
	httpServerPort = ":8080"
	// TARA-Mock-i HTTPS sert
	taraMockCert = "vault/https.crt"
	// TARA-Mock-i HTTPS privaatvõti
	taraMockKey = "vault/https.key"
	// TARA-Mock-i identsustõendi allkirjastamise avalik võti
	idTokenPrivKeyPath = "vault/idtoken.key"
	// TARA-Mock-i identsustõendi allkirjastamise privaatvõti
	idTokenPubKeyPath = "vault/idtoken.pub"
	// Identsustõendi allkirjavõtme identifikaator
	kid = "taramock"
)

// Identsustõendi allkirjastamise RSA võtmepaar
var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

type volituskood string

// Andmed identsustõendi moodustamiseks ja väljastamiseks.
// Identsustõend koostatakse vahetult enne väljastamist.
type dataForTokenType struct {
	sub        string // subject, isikutõendi väli "sub"
	familyName string // family_name
	givenName  string // given_name
	state      string // autentimispäringus saadetud turvaväärtus
	nonce      string // autentimispäringus saadetud turvaväärtus
}

// Identsustõendite hoidla
var idToendid = make(map[volituskood]dataForTokenType)

var mutex = &sync.Mutex{}

func main() {

	// Marsruudid
	// Go-s "/" käsitleb ka need teed, millele oma käsitlejat ei leidu.
	http.HandleFunc("/", landingPage)
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/oidc/.well-known/openid-configuration", sendConf)
	http.HandleFunc("/oidc/authorize", authenticateUser)
	http.HandleFunc("/back", sendUserBack)
	http.HandleFunc("/oidc/token", sendIdentityToken)
	http.HandleFunc("/oidc/jwks", sendKey)

	// Loe sisse identsustõendi allkirjastamise võtmepaar.
	readRSAKeys()

	// fileServer serveerib kasutajaliidese muutumatuid faile.
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Käivita HTTPS server
	fmt.Printf("** TARA-Mock käivitatud pordil %v **\n", httpServerPort)
	err := http.ListenAndServeTLS(
		httpServerPort,
		taraMockCert,
		taraMockKey,
		nil)
	if err != nil {
		log.Fatal(err)
	}
}

// readRSAKeys loeb sisse identsustõendi allkirjastamise võtmepaari
// ja valmistab ette allkirjastamise avaliku võtme otspunkti.
// Kasutab teeki dgrijalva/jwt-go.
func readRSAKeys() {
	// Vt: https://github.com/dgrijalva/jwt-go/blob/master/http_example_test.go
	signBytes, err := ioutil.ReadFile(idTokenPrivKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}

	verifyBytes, err := ioutil.ReadFile(idTokenPubKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
	}
}

// Kaalu allkirjastamispaketi kasutamist:
// "github.com/lestrrat-go/jwx/jws"
// Dok-n: https://godoc.org/github.com/lestrrat-go/jwx/jws
// Alternatiiv on ka: square/go-jose (v3), vt:
// https://godoc.org/github.com/square/go-jose
