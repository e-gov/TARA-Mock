package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims on truktuur, mis kodeeritakse veebitõendina (JWT).
type Claims struct {
	Jti               string `json:"jti"`
	Issuer            string `json:"iss"`
	Audience          string `json:"aud"`
	ExpiresAt         string `json:"exp"`
	IssuedAt          string `json:"iat"`
	NotBefore         string `json:"nbf"`
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

// SendIdentityToken väljastab klientrakendusele identsustõendi (OIDC identsustõendi otspunkt /oidc/token).
func SendIdentityToken(w http.ResponseWriter, r *http.Request) {
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

	// Paki päringu keha lahti (b -> msg).
	var msg IDTokenReqBody
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Võta identsustõend (v) mälus hoitavast identsustõendite hoidlast.
	v, ok := idToendid[msg.Code]
	if !ok {
		http.Error(w, "Identsustõend ei eksisteeri", 404)
		return
	}

	/* Primitiivne meetod: Koosta identsustõend
	mt := IdentityToken{v.sub, v.givenName, v.familyName}
	output, err := json.Marshal(mt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	*/

	// Koosta JWT
	// Koosta JWT väited
	claims := &Claims{
		Jti:      "001",
		Issuer:   "TARA-Mock",
		Audience: "Klientrakendus",
		// Identsustõendi kehtivusaeg - 1 minute
		ExpiresAt: string(time.Now().Add(1 * time.Minute).Unix()),
		IssuedAt:  string(time.Now().Unix()),
		NotBefore: string(time.Now().Unix()),
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

	// Allkirjavõti
	var jwtKey = []byte("my_secret_key")

	// Koosta veebitõend
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Väljasta
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(tokenString))

	// Primitiivne meetod:
	// w.Write(output)
}
