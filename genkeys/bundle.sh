#!/bin/bash

echo "--- bundle.sh"
echo "Skript TARA-Mock elutukse otspunkti poole sirvikuga pöördumiseks"
echo "vajaliku sertipaki (PKCS#12 bundle). Pakk on vajalik, sest TARA-Mock"
echo "HTTPS server on seadistatud klienti autentima. Serdipakk tuleb laadida"
echo "sirvikusse."
echo "Export password küsimusele vajuta Enter ja salasõna ei looda."
echo " "

echo
echo "    Valmistan serdipaki"
openssl pkcs12 \
  -export \
  -in https.crt \
  -inkey https.key \
  -out bundle.p12 \

#  -certfile rootCA.pem \

# Kuva valmistatud pakk
echo "--- Valmistatud serdipakk:"
openssl pkcs12 -info \
  -in bundle.p12

# Abiteave
# Vt: https://medium.com/@sevcsik/authentication-using-https-
# client-certificates-3c9d270e8326