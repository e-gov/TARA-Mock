# TARA-Mock

TARA-Mock on rakendus, mis etendab TARA autentimist. 

TARA-Mock on mõeldud kasutamiseks siis, kui [TARA testteenuse](https://e-gov.github.io/TARA-Doku/Testimine) võimalused jäävad klientrakenduse funktsionaalsuste testimisel napiks. TARA testteenuse abil saab autentida väga väikese hulga TARA poolt ette antud testkasutajatega. TARA-Mock seevastu võimaldab kasutajal valida autentimise dialoogis identiteet kas TARA-Mock seadistuses etteantud identiteetide hulgast või sisestada ise isikukood, ees- ja perekonnanimi, mille all soovitakse sisse logida s.t logida sisse suvalise identiteediga.

TARA-Mock pakub samu otspunkte kui TARA, s.t TARA-Mock ja TARA on üksteisega vahetatavad lihtsa seadistusega.

TARA-Mock juures on ka klientrakenduse näidis.

## Kasutusstsenaarium

1) Suunamine autentimisele

<img src="docs/Rakendus_01.PNG" width="700">

2) Autentimisdialoog

<img src="docs/TARA-Mock_01.PNG" width="700">

3) Tagasi autentimiselt

<img src="docs/Rakendus_02.PNG" width="700">

## Lihtsustused

TARA-Mock on tehtud rida lihtsustusi ja jäetud ära kontrolle:

- aktsepteeritakse kõiki klientrakendusi (`client_id` väärtust ei kontrollita)
- turvaelemendid (`state` ja `nonce`) antakse ühes sammus edasi HTML vormi peidetud väljadena
- puudub päringuvõltsimise kaitse (CSRF)
- juhusõned genereeritakse tavalise (`math/rand`), mitte krüptograafilise juhuarvugeneraatoriga (`crypto/rand`)
- ainult Eesti isikukoodiga isikute autentimine
- piiratud logimine: TARA-Mock väljastab mõningast logiteavet konsoolile
- klientrakenduse salasõna ei kontrollita
- identsustõendi väljastamisel `redirect_uri` ei kontrollita; identsustõend väljastatakse ainult volituskoodi alusel
- identsustõendi väljastamisel ei kontrollita, kas tõend on aegunud

## Kaitsed

TARA-Mock töötab HTTPS-ga, rakendatakse TLS kliendi autentimist.

## TARA-Mock

TARA-Mock töötab pordil 8080.

**Seadistamine**:

Genereeri ja sea TARA-Mock HTTPS serveri serdid. Sertide genereerimise näiteskript on failis `keys/genkeys.sh`.

Etteantud identiteedid on koodi sissekirjutatud TARA-Mock koodi. Muuda identiteedid oma vajadustele vastavaks.

**Paigaldamine**:

Paigalda TARA-Mock sobivasse masinasse. TARA-Mock on kasutatav ka oma masinas (`localhost`).

**Otspunktid**:
- `/health` - elutukse
- `/` - avaleht teabega TARA-Makett kohta
- `/oidc/authorize` - autentimisele suunamine
- `/token` - identsustõendi väljastamine
- `/oidc/jwks` - identsustõendi avalik võti

Nt TARA-Mock kasutamisel oma masinas:

`https://localhost:8080/oidc/token` - identsustõendi väljastamine

**Käivitamine**:

```
cd service
go run .
```

## Klientrakenduse näidis

**Otspunktid**:
- `/health` - elutukse
- `/` - avaleht; kasutaja saab sealt minna TARA-Mock-i autentima
- `/login` - kasutaja suunamine TARA-Mock-i autentima
- `/return` - autentimiselt tagasi suunatud kasutaja vastuvõtmine, identsustõendi pärimine TARA-Mock-st ja sisselogimise lõpuleviimine 

**Seadistamine**:

Klientrakenduse näidis töötab lokaalses masinas, pordil 8081. Genereeri ja sea TARA-Mock HTTPS serveri serdid. Sertide genereerimise näiteskript on failis `keys/genkeys.sh`.

**Käivitamine**:

```
cd service
go run .
```
