#!/bin/bash

RED='\033[0;31m'
NC='\033[0m'

echo -e "${RED}--- genkeys.sh"
echo -e "    Skript moodustab TARA-Mock tööks vajalikud võtmed ja serdid."
echo -e "    Vt: https://github.com/e-gov/TARA-Mock/blob/master/docs/Turvalisus.md"

echo -e " "
echo -e "--- NB! subj väärtustes tuleb Git for Windows kasutamisel seada"
echo -e "    tee-eraldajad: //CN\ ..."
echo -e " "
echo -e "    NB! Skript ei tööta WSL-s (Windows Subsystem for Linux)"
echo -e " "

echo -e " "
echo -e "--- 1 Valmistan CA võtme ja serdi${NC}"
openssl req \
  -new \
  -x509 \
  -newkey rsa:2048 \
  -keyout rootCA.key \
  -out rootCA.pem \
  -nodes \
  -days 1024 \
  -subj "/C=EE/ST=/L=/O=TEST-CA/CN=TEST-CA"

# Kuva subject ja issuer
echo -e "${RED}--- Valmistatud CA sert:${NC}"
openssl x509 \
  -in rootCA.pem \
  -noout \
  -subject -issuer

echo -e "${RED} "
echo -e "--- 2 Valmistan TARA-Mock HTTPS privaatvõtme ja serdi${NC}"
# Serditaotlus
openssl req \
  -new \
  -sha256 \
  -nodes \
  -out https.crs \
  -newkey rsa:2048 \
  -keyout https.key \
  -subj "/C=EE/ST=/L=/O=Arendaja/CN=Arendaja"
# Sert
openssl x509 \
  -req \
  -in https.crs \
  -CA rootCA.pem \
  -CAkey rootCA.key \
  -CAcreateserial \
  -out https.crt \
  -days 500 \
  -sha256 \
  -extfile v3.ext

# Kuva subject ja issuer
echo -e "${RED} "
echo -e "--- Valmistatud sert:${NC}"
openssl x509 \
  -in https.crt \
  -noout \
  -subject -issuer

echo -e "${RED} "
echo -e "--- 3 Genereerin identsustõendi allkirjastamise privaat- ja avaliku võtme${NC}"
openssl genrsa \
  -out idtoken.key \
  2048
openssl rsa \
  -in idtoken.key \
  -pubout > idtoken.pub

echo -e "${RED} "
echo -e "--- 4 Eemaldan vanad võtmed ja serdid${NC}"
rm -f ../service/vault/*.*
rm -f ../client/vault/*.*

echo -e "${RED} "
echo -e "--- 5 Paigaldan võtmed ja serdid TARA-Mock-i${NC}"
cp rootCA.pem ../service/vault
cp https.key ../service/vault
cp https.crt ../service/vault
cp idtoken.key ../service/vault
cp idtoken.pub ../service/vault

echo -e "${RED} "
echo -e "--- 6 Paigaldan võtmed ja serdid klientrakendusse${NC}"
cp rootCA.pem ../client/vault
cp https.key ../client/vault
cp https.crt ../client/vault

echo -e "${RED}--- Võtmed ja serdid genereeritud ja paigaldatud"
echo -e "--- Ära unusta sirvikusse usaldusankrut paigaldada${NC}"

# -------------------
# Abiteave
# These work for application/json, but not for text/html in browser
# openssl genrsa -out https.key 2048
# openssl ecparam -genkey -name secp384r1 -out https.key
# openssl req -new -x509 -sha256 -key https.key -out https.crt -days 3650

# Oluline allikas: https://www.freecodecamp.org/news/how-to-get-https-working-on-your-local-development-environment-in-5-minutes-7af615770eec/

# Windows-specific
# https://stackoverflow.com/questions/31506158/running-openssl-from-a-bash-script-on-windows-subject-does-not-start-with
# https://stackoverflow.com/questions/7360602/openssl-and-error-in-reading-openssl-conf-file
# set OPENSSL_CONF=c:/libs/openssl-0.9.8k/openssl.cnf
# echo %OPENSSL_CONF%
# echo
