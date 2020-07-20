package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
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
	Amr   []string `json:"amr"` // Autentimismeetod
	State string   `json:"state"`
	Nonce string   `json:"nonce"`
	Acr   string   `json:"acr"` // Autentimistase
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

	// Võta päringu Query-osa (-> m)
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Debugf("ID Token request/Identsustõendi päring: %s", string(b))
	m, err := url.ParseQuery(string(b))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		for k, v := range m {
			log.Debugf("    %s: %s", k, v)
		}
	}

	// Võta identsustõendile vajalikud andmed (v) mälus hoitavast
	// identsustõendite andmete hoidlast.
	v, ok := idToendid[volituskood(m.Get("code"))]
	if !ok {
		http.Error(w, "Identsustõendile vajalikke andmeid ei leia", 404)
		return
	}

	// Koosta JWT väited
	claims := &Claims{
		Jti:      "001",
		Issuer:   "https://" + conf.TaraMockHost + conf.HTTPServerPort,
		Audience: v.clientID,
		// Identsustõendi kehtivusaeg - 1 minute
		ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		Subject:   v.sub,
		Amr:       []string{"mID"},
		State:     v.state,
		Nonce:     v.nonce,
		Acr:       "high", // Tagatistase 'kõrge' )
	}

	// Moodusta sünnikp isikukoodist.
	if dob, err := personCodeToDoB(v.sub); err != nil {
		// Kui sünnikp ei saa moodustada, siis kasuta fiks-d väärtust
		claims.ProfileAttributes.DateOfBirth = "1961-07-12"
	} else {
		claims.ProfileAttributes.DateOfBirth = dob
	}

	claims.ProfileAttributes.GivenName = v.givenName
	claims.ProfileAttributes.FamilyName = v.familyName

	// Koosta veebitõend
	// Token struktuuri vt:
	// https://github.com/dgrijalva/jwt-go/blob/master/token.go#L23
	// Oluline näide:
	// https://gist.github.com/cryptix/45c33ecf0ae54828e63b
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = conf.Kid
	// Create the JWT string
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		log.Errorf("Viga veebitõendi allkirjastamisel: %v", err)
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
	log.Debug("ID Token issued/Identsustõend väljastatud")
}

// jwt teegi näide: https://github.com/dgrijalva/jwt-go/issues/141
// Vt: https://stackoverflow.com/questions/24809235/initialize-a-nested-struct
