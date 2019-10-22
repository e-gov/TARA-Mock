#!/bin/bash

echo "--------------------------"
echo "Skript moodustab TARA-Mock-i poolt väljastatava"
echo "identsustõendi allkirjastamise privaat- ja avaliku võtme:"
echo "  idtoken.key"
echo "  idtoken.pub"
echo "Vt: https://github.com/e-gov/TARA-Mock/blob/master/docs/Turvalisus.md"

echo
echo "### 3 Genereerin identsustõendi allkirjastamise privaat- ja avaliku võtme"
openssl genrsa \
  -out idtoken.key \
  2048
openssl rsa \
  -in idtoken.key \
  -pubout > idtoken.pub


