package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// getIdentityToken pärib TARA-Mock-lt identsustõendi.
func getIdentityToken(vk string) ([]byte, bool) {

	// Lae kliendi võti ja sert
	cert, err := tls.LoadX509KeyPair(
		"keys/https-server.crt",
		"keys/https-server.key")
	if err != nil {
		log.Fatal(err)
	}

	// Lae CA sert
	caCert, err := ioutil.ReadFile(
		"keys/rootCA.pem",
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
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Koosta POST päringu keha
	requestBody, err := json.Marshal(map[string]string{
		"grant_type":   "authorization_code",
		"code":         vk,
		"redirect_uri": redirectURI,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Saada POST päring
	resp, err := client.Post(
		"https://localhost:8080/oidc/token",
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Loeb päringu keha, kujule []byte
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body, true
}
