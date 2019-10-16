#!/bin/bash

echo "--------------------------"
echo "Skript moodustab serdid, mida saab kasutada TARA-Mock ja"
echo "TARA-Mock klientrakenduse näite sertidena. Serdid moodustatakse"
echo "nimedega CN=Arendaja, kasutamiseks juhul, kui nii TARA-Mock"
echo "kui ka klientrakenduse näide paigaldatakse samasse masinasse"
echo "(localhost)."
echo "--------------------------"
echo "### NB! Kasutusel Git for Windows tee-eraldajad"


# Oluline allikas: https://www.freecodecamp.org/news/how-to-get-https-working-on-your-local-development-environment-in-5-minutes-7af615770eec/

# https://stackoverflow.com/questions/31506158/running-openssl-from-a-bash-script-on-windows-subject-does-not-start-with

# Windows-specific
# https://stackoverflow.com/questions/7360602/openssl-and-error-in-reading-openssl-conf-file
# set OPENSSL_CONF=c:/libs/openssl-0.9.8k/openssl.cnf
# echo %OPENSSL_CONF%
# echo

echo "### Moodustan võtme ja serdi root CA jaoks"
openssl req \
  -new \
  -x509 \
  -newkey rsa:2048 \
  -keyout rootCA.key \
  -out rootCA.pem \
  -nodes \
  -days 1024 \
  -subj "//C=EE\ST=\L=\O=Arendaja\CN=Arendaja"

# Kuva subject ja issuer
echo "OK"
echo "CA-le moodustatud sert:"
openssl x509 \
  -in rootCA.pem \
  -noout \
  -subject -issuer

# Moodusta serditaotlus localhost HTTPS serverile
openssl req \
  -new \
  -sha256 \
  -nodes \
  -out https-server.crs \
  -newkey rsa:2048 \
  -keyout https-server.key \
  -subj "//C=EE\ST=\L=\O=Arendaja\CN=Arendaja"

# Moodusta sert localhost HTTPS serverile
openssl x509 \
  -req \
  -in https-server.crs \
  -CA rootCA.pem \
  -CAkey rootCA.key \
  -CAcreateserial \
  -out https-server.crt \
  -days 500 \
  -sha256 \
  -extfile v3.ext

# Kuva subject ja issuer
echo
echo "Lokaalses masinas olevale HTTPS serverile moodustatud sert:"
openssl x509 \
  -in https-server.crt \
  -noout \
  -subject -issuer

# These work for application/json, but not for text/html in browser
# openssl genrsa -out https-server.key 2048
# openssl ecparam -genkey -name secp384r1 -out https-server.key
# openssl req -new -x509 -sha256 -key https-server.key -out https-server.crt -days 3650

