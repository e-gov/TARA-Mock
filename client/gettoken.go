package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	// Identsustõendi verifitseerimiseks
	// Asendatud (apr 2022):
	// Dok-n: https://godoc.org/github.com/dgrijalva/jwt-go
	// "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v4"
	// JWK (veebivõtme) toiminguteks
	// Dok-n: https://godoc.org/github.com/lestrrat-go/jwx
	"github.com/lestrrat-go/jwx/jwk"
)

var idTokenPublicKey *rsa.PublicKey

// MyCustomClaims kirjeldab identsustõendi väited.
type MyCustomClaims struct {
	ProfileAttributes struct {
		DateOfBirth string `json:"date_of_birth"`
		GivenName   string `json:"given_name"`
		FamilyName  string `json:"family_name"`
	} `json:"profile_attributes"`
	Amr   []string `json:"amr"` // Autentimismeetod
	State string   `json:"state"`
	Nonce string   `json:"nonce"`
	Acr   string   `json:"acr"` // Autentimistase
	// Vt: https://godoc.org/github.com/dgrijalva/jwt-go#StandardClaims
	jwt.StandardClaims
}

// Valid kontrollib identsustõendi õigsust.
func (MyCustomClaims) Valid() error {
	return nil
}

// getKey tagastab identsustõendi allkirja avaliku võtme.
// Vt: https://stackoverflow.com/questions/41077953/go-language-and-verify-jwt
func getKey(token *jwt.Token) (interface{}, error) {
	// idTokenPublicKey tüüp on jwk.Key. Kuidas sobib teegiga jwt??
	// Vt näide: https://godoc.org/github.com/lestrrat-go/jwx/jwk
	return idTokenPublicKey, nil
}

// getIdentityToken : 1) pärib TARA-Mock-lt identsustõendi allkirja
// avaliku võtme (otspunktist /oidc/jwks); 2) pärib TARA-Mock-lt
// identsustõendi; 3) parsib identsustõendi;
// 4) tagastab identsustõendilt loetud isikuandmed, stringina.
// TO DO: Löö f-deks. Probleem: kuidas http.Client üle kanda?
func getIdentityToken(vk string) (string, bool) {

	// Loe kliendi HTTPS võti ja sert
	cert, err := tls.LoadX509KeyPair(
		conf.AppCert,
		conf.AppKey)
	if err != nil {
		log.Fatal(err)
	}

	// Loe CA sert
	caCert, err := ioutil.ReadFile(
		conf.RootCAFile,
	)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Sea HTTPS klient valmis
	// Vt: https://golang.org/pkg/net/http/
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	// tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// client := &http.Client{Timeout: 10 * time.Second}

	// ----------------
	// Päri allkirja avalik võti
	fmt.Printf("\ngetIdentityToken:\n    Pärin identsustõendi allkirjavõtme\n")
	resp1, err := client.Get(conf.TaraMockKeyEndpoint)
	if err != nil {
		log.Fatalln("Viga allkirja avaliku võtme pärimisel: ", err)
	}
	defer resp1.Body.Close()

	type Key struct {
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
	}

	type KeySet struct {
		// Keys []jwk.Key `json:"keys"`
		Keys []Key `json:"keys"`
	}

	// Vt: https://stackoverflow.com/questions/17156371/how-to-get-json-response-from-http-get
	decoder := json.NewDecoder(resp1.Body)
	var data KeySet
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatalln("getIdentityToken: Viga võtmepäringu vastuse dekodeerimisel: ", err)
	}
	// fmt.Println("Saadud võti kid: ", data.Keys[0].Kid)
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("getIdentityToken: Viga JSON-kujule teisendamisel: ", err)
	}
	// fmt.Println("Saadud võtmed: ", string(jsonBytes))

	// Teisenda võtmehulk (stringina) teegi lestrrat-go/jwx/jwk
	// pakutud Go-kujule. kSet tüüp on jwk.Set.
	kSet, err := jwk.ParseString(string(jsonBytes))
	if err != nil {
		log.Fatalln("getIdentityToken: Viga teisendamisel lestrrat-go/jwx/jwk kujule: ", err)
	}
	fmt.Println("    Saadud võti, kid = ", kSet.Keys[0].KeyID())
	// Materialize() peaks jwk.Key-st tegema *rsa.PublicKey.
	m, err := kSet.Keys[0].Materialize()
	if err != nil {
		log.Printf("getIdentityToken: Avaliku RSA võtme moodustamine ebaõnnestus: %s", err)
		return "Viga: Avaliku RSA võtme moodustamine ebaõnnestus", true
	}
	idTokenPublicKey = m.(*rsa.PublicKey)

	// ----------------
	// Päri identsustõend
	// Koosta POST päringu query-osa ja keha
	qp := "grant_type=authorization_code" +
		"&code=" + vk +
		"&redirect_uri=" + conf.RedirectURI
	fmt.Printf("--- Pärin identsustõendi: %v\n", qp)
	var requestBody = []byte(qp)

	// Saada POST päring
	resp, err := client.Post(
		conf.TaraMockTokenEndpoint,
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Loe vastuse keha, kujule []byte
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("    Saadud vastus: ", string(body))

	type IDTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		// "RECOMMENDED.  The lifetime in seconds of the access token."
		// -- OAuth 2.0 spec
		ExpiresIn int    `json:"expires_in"`
		IDToken   string `json:"id_token"`
	}
	var IDTR IDTokenResponse

	// Parsi JSON
	if err := json.Unmarshal(body, &IDTR); err != nil {
		log.Printf("Vastuse JSON parsimine ebaõnnestus: %s", err)
	}
	// fmt.Println("token_type: ", IDTR.TokenType)

	// Parsi ja kontrolli JWT
	// Vt: https://stackoverflow.com/questions/41077953/go-language-and-verify-jwt
	fmt.Println("--- Kontrollin saadud identsustõendit")
	var p1 jwt.Parser
	_, err2 := p1.Parse(IDTR.IDToken, getKey)
	if err2 != nil {
		log.Printf("Identsustõendi kontroll ebaõnnestus %s", err)
		return "Identsustõendi kontroll ebaõnnestus", true
	}
	/*
		fmt.Println("    Tõendil on väited:")
		claims1 := token1.Claims.(jwt.MapClaims)
		for key, value := range claims1 {
			fmt.Printf("    %s\t%v\n", key, value)
		}
	*/

	var myClaims MyCustomClaims

	// Parsi tõend ilma allkirjakontrollita. Saaks ühitada eelmise
	// parsimisega.
	var p jwt.Parser
	t, _, err := p.ParseUnverified(IDTR.IDToken, &myClaims)
	if err != nil {
		log.Printf("JWT parsimine ebaõnnestus: %s", err)
		return "Identsustõendi töötlemine ebaõnnestus", true
	}

	fmt.Println("--- Tõendilt loetud:")
	fmt.Println("    võtme id (kid): ", t.Header["kid"].(string))
	fmt.Println("    algoritm (alg): ", t.Header["alg"].(string))
	fmt.Println("    tüüp (typ): ", t.Header["typ"].(string))
	claims := t.Claims.(*MyCustomClaims)
	fmt.Printf("    state %v\n", claims.State)
	fmt.Printf("    nonce %v\n", claims.Nonce)
	fmt.Printf("    autentimismeetod %v\n", claims.Amr)
	fmt.Printf("    tagatistase %v\n", claims.Acr)
	fmt.Println("    -- standardväited --")
	fmt.Printf("    id (jti): %v\n", claims.StandardClaims.Id)
	fmt.Printf("    väljaandja (iss): %v\n", claims.StandardClaims.Issuer)
	fmt.Printf("    (aud): %v\n", claims.StandardClaims.Audience)
	fmt.Printf("    kehtib kuni (exp): %v\n", time.Unix((claims.StandardClaims.ExpiresAt), 0))
	fmt.Printf("    väljaandmiskp (iat): %v\n", time.Unix((claims.StandardClaims.IssuedAt), 0))
	fmt.Printf("    mitte enne (nbf): %v\n", time.Unix((claims.StandardClaims.NotBefore), 0))
	fmt.Printf("    subjekt (sub): %v\n", claims.StandardClaims.Subject)
	fmt.Println("    -- isiku profiil --")
	fmt.Printf("    eesnimi %v\n", claims.ProfileAttributes.GivenName)
	fmt.Printf("    perekonnanimi %v\n", claims.ProfileAttributes.FamilyName)
	fmt.Printf("    sünniaeg %v\n", claims.ProfileAttributes.DateOfBirth)

	return claims.StandardClaims.Subject + ", " +
		claims.ProfileAttributes.GivenName + ", " +
		claims.ProfileAttributes.FamilyName + ", " +
		claims.ProfileAttributes.DateOfBirth, true
}

// Token Claim tüübikinnitamise (type assertion) näide:
// https://github.com/gomango/clientside-skeleton/blob/master/main.go#L73
