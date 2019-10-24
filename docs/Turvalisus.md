# Turvalisus

[Võtmeplaan](#võtmeplaan) · 
[Võtmete ja sertide valmistamine](#võtmete-ja-sertide-valmistamine) · 
[Turvalisus sirvikus](#turvalisus-sirvikus)

## Võtmeplaan

TARA-Mock vajab oma tööks järgmisi võtmeid ja serte (kaustas `service/vault`):

- `https.key` - TARA-Mock HTTPS serveri privaatvõti
- `https.crt` - TARA-Mock HTTPS serveri sert
- `idtoken.key` - identsustõendi allkirjastamise privaatvõti
- `idtoken.pub` - identsustõendi allkirjastamise avalik võti
- `rootCA.pem` - klientrakenduse serdi väljaandja sert

Sirvikus tuleb luua usaldus TARA-Mock-i vastu. Selleks tuleb sirvikusse,  millega TARA-Mock-i kasutatakse, paigaldada TARA-Mock HTTPS serveri serdi väljaandja (edaspidi - TARA-Mock-i CA) sert:

- `rootCA.pem` - TARA-Mock-i CA sert

Usaldusankru sirvikusse paigaldamise juhiseid vt jaotises [Sirvik](#sirvik)

Klientrakenduses, mis TARA-Mock-i poole pöördub, tuleb paigaldada TARA-Mock-i CA sert. Samuti peab klientrakendusel endal olema privaatvõti ja sert. Oma privaatvõtit ja serti kasutab klientrakendus nii oma HTTPS serveris (kasutajaliidese pakkumisel läbi sirviku) kui ka HTTPS kliendina, pöördumisel TARA-Mock-i poole (identsustõendi pärimisel). Need failid tuleb panna kausta `client/vault`:

- `rootCA.pem` - TARA-Mock-i CA sert
- `https.key` - klientrakenduse privaatvõti
- `https.crt` - klientrakenduse sert

## Võtmete ja sertide valmistamine

Võtmed ja serdid saab valmistada standardsete vahenditega (OpenSSL), vastavalt paigalduskonfiguratsioonile.

Soovi korral võib kasutada TARA-Mock koodirepo koosseisus (kaustas `genkeys`) olevat skripti `genkeys.sh`. Skript valmistab võtmed ja serdid konfiguratsioonile, kus:

- nii TARA-Mock kui ka  klientrakendus paigaldatakse samasse masinasse (localhost)
- TARA-Mock-l ja klientrakendusel on sama CA
- CA organisatsiooninimi (`O=`) ja nimi (`CN=`) on `Arendaja`
- TARA-Mock-i ja klientrakenduse organisatsiooninimi (`O=`) ja nimi (`CN=`) on samuti `Arendaja`.

Muu konfiguratsiooni puhul saab skripti muuta, kirjutades sisse õiged hostinimed, `O=` ja `CN=` vaartused ja tehes võtmete ning sertide kaustadesse paigutamise muude vahenditega. NB! Skriptis on failiteed Windows-i nõuete kohaselt. Linux-i kasutamisel muuta vastavalt.

Skriptis täidetakse järgmised sammud:

1 Sertifitseerimisasutuse (CA) privaatvõtme ja serdi valmistamine. Valmistatakse failid: `rootCA.key` (privaatvõti) ja `rootCA.pem` (sert). Samuti moodustatakse fail `rootCA.srl`, milles CA peab arvestust sertide seerianumbrite üle).

```
openssl req \
  -new \
  -x509 \
  -newkey rsa:2048 \
  -keyout rootCA.key \
  -out rootCA.pem \
  -nodes \
  -days 1024 \
  -subj "//C=EE\ST=\L=\O=Arendaja\CN=Arendaja"
```

Kaitske CA privaatvõtit.

2 TARA-Mock HTTPS serveri privaatvõtme (`https.key`) ja serdi (`https.crt`) valmistamine. Kuna TARA-Mock HTTPS server teenindab ka sirvikuid, siis on vaja serdile kantav väli `subjectAltName` määratleda eraldi failis `v3.ext`. Kui paigalduskohaks ei ole `localhost`, siis tuleb failis `v3.ext` seada õige hostinimi. 

```
# Moodusta serditaotlus localhost HTTPS serverile
openssl req \
  -new \
  -sha256 \
  -nodes \
  -out https.crs \
  -newkey rsa:2048 \
  -keyout https.key \
  -subj "//C=EE\ST=\L=\O=Arendaja\CN=Arendaja"

# Moodusta sert localhost HTTPS serverile
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
```

3 Identsustõendi allkirjastamise võtmepaari valmistamine. Valmistatakse failid `idtoken.key` (privaatvõti) ja `idtoken.pub` (avalik võti):

```
openssl genrsa \
  -out idtoken.key \
  2048
openssl rsa \
  -in idtoken.key \
  -pubout > idtoken.pub
```
4 Võtmete ja sertide paigaldamine TARA-Mock-i ja klientrakendusse (vastavalt kaustadesse `service/vault` ja `client/vault`).

## Turvalisus sirvikus

Sirvikusse, millega TARA-Mock-i kasutatakse, tuleb paigaldada TARA-Mock-i CA sert. Sellega pannakse sirvik TARA-Mock-i usaldama. Usalduseta ei ava sirvik TARA-Mock-i kasutajaliidest.

- kopeeri TARA-Mock-i CA sert sirviku arvutisse (või tee muul viisil kättesaadavaks)
  - Chrome: `chrome://settings/privacy`, `Manage Certificates`, `Trusted Root Certification Authorities`
  - Firefox: `Tools`, `Options`, `Privacy & Security`, `Certificates`, `View Certificates`, `Authorities`, `Import`

Kui kasutate ka TARA-Mock-ga kaasas olevat klientrakendust, siis peab ka selle CA serdi paigaldama sirvikusse. Kui TARA-Mock-i ja klientrakenduse serdid anda välja ühe CA poolt, piisab ühest paigaldamisest. 

TARA-Mock-i kasutamise järel (kui rakendus sai testitud), eemaldage CA sert sirvikust.