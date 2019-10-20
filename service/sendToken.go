package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims on truktuur, mis kodeeritakse veebitõendina (JWT).
type Claims struct {
	Jti               string `json:"jti"`
	Issuer            string `json:"iss"`
	Audience          string `json:"aud"`
	ExpiresAt         int64  `json:"exp"`
	IssuedAt          int64  `json:"iat"`
	NotBefore         int64  `json:"nbf"`
	Subject           string `json:"sub"`
	ProfileAttributes struct {
		DateOfBirth string `json:"date_of_birth"`
		GivenName   string `json:"given_name"`
		FamilyName  string `json:"family_name"`
	} `json:"profile_attributes"`
	Amr   string `json:"amr"` // Autentimismeetod
	State string `json:"state"`
	Nonce string `json:"nonce"`
	Acr   string `json:"acr"` // Autentimistase
}

// Valid kontrollib identsustõendi õigsust.
func (Claims) Valid() error {
	return nil
}

// sendIdentityToken väljastab klientrakendusele identsustõendi
// (otspunkt /oidc/token).
func sendIdentityToken(w http.ResponseWriter, r *http.Request) {
	// Võta päringust volituskood, seejärel võta volituskoodile
	// vastavad identsustõendi andmed, koosta identsustõend ja saada
	// päringu vastuses.

	// Loe päringu keha (b).
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// IDTokenReqBody on vastuvõetava identsustõendi päringu keha
	// struktuur.
	type idTokenReqBody struct {
		GrantType  string      `json:"grant_type"`
		Code       volituskood `json:"code"`
		RequestURI string      `json:"request_uri"`
	}

	// Paki päringu keha lahti (b -> msg).
	var msg idTokenReqBody
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Võta identsustõendile vajalikud andmed (v) mälus hoitavast
	// identsustõendite andmete hoidlast.
	v, ok := idToendid[msg.Code]
	if !ok {
		http.Error(w, "Identsustõendile vajalikke andmeid ei leia", 404)
		return
	}

	// Koosta JWT
	// Koosta JWT väited
	claims := &Claims{
		Jti:      "001",
		Issuer:   "TARA-Mock",
		Audience: "Klientrakendus",
		// Identsustõendi kehtivusaeg - 1 minute
		ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		Subject:   v.sub,
		Amr:       "mID",
		State:     v.state,
		Nonce:     v.nonce,
		Acr:       "high",
	}
	claims.ProfileAttributes.DateOfBirth = "1961-07-12"
	claims.ProfileAttributes.GivenName = v.givenName
	claims.ProfileAttributes.FamilyName = v.familyName
	// Vt: https://stackoverflow.com/questions/24809235/initialize-a-nested-struct

	// Koosta veebitõend
	// Token struktuuri vt:
	// https://github.com/dgrijalva/jwt-go/blob/master/token.go#L23
	// Oluline näide:
	// https://gist.github.com/cryptix/45c33ecf0ae54828e63b
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "taramock"
	// Create the JWT string
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		log.Printf("Viga veebitõendi allkirjastamisel: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	// Koosta vastuse keha
	type IDTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		// "RECOMMENDED.  The lifetime in seconds of the access token."
		// -- OAuth 2.0 spec
		ExpiresIn int    `json:"expires_in"`
		IDToken   string `json:"id_token"`
	}
	var IDTR IDTokenResponse
	IDTR.AccessToken = "eiolekasutusel"
	IDTR.TokenType = "bearer"
	IDTR.ExpiresIn = 3600
	IDTR.IDToken = tokenString

	saadetis, err := json.Marshal(IDTR)
	if err != nil {
		log.Printf("Viga veebitõendi väljastamisel: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	// Väljasta
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	//	w.Write([]byte(tokenString))
	w.Write(saadetis)
}

// jwt teegi näide: https://github.com/dgrijalva/jwt-go/issues/141
