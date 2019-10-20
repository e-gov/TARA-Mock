# TARA-Mock

TARA-Mock on rakendus, mis etendab TARA autentimist. 

## Ülevaade

TARA-Mock on mõeldud kasutamiseks siis, kui [TARA testteenuse](https://e-gov.github.io/TARA-Doku/Testimine) võimalused jäävad klientrakenduse funktsionaalsuste testimisel napiks. TARA testteenuse abil saab autentida väga väikese hulga TARA poolt ette antud testkasutajatega.

TARA-Mock seevastu võimaldab kasutajal valida autentimise dialoogis identiteet  TARA-Mock seadistuses etteantud identiteetide hulgast või sisestada ise isikukood, ees- ja perekonnanimi, mille all soovitakse sisse logida. Sisuliselt saab sisse logida ükskõik millise identiteediga.

TARA-Mock on reaalse TARA-ga ühilduv, s.t TARA-Mock ja TARA on üksteisega vahetatavad lihtsa seadistusega.

TARA-Mock juures on ka klientrakenduse näidis.

TARA-Mock ei ole mõeldud kasutamiseks toodangus. TARA-Mock ei ole mõeldud ka TARA-ga liidestamise testimiseks - sest TARA-Mock-is on ära jäetud mitmeid toodangus vajalikke kontrolle (vt allpool).

TARA-Mock on kirjutatud Go-s.

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

## Paigaldamine

TARA-Mock töötab pordil 8080. Paigalda TARA-Mock sobivasse masinasse. TARA-Mock on kasutatav ka oma masinas (`localhost`). Otspunktid:

- `/health` - elutukse
- `/` - avaleht teabega TARA-Makett kohta
- `/oidc/authorize` - autentimisele suunamine
- `/token` - identsustõendi väljastamine
- `/oidc/jwks` - identsustõendi avalik võti

Nt TARA-Mock kasutamisel oma masinas: `https://localhost:8080/health`.

1 Masinas peab olema paigaldatud Go, versioon 1.11 või hilisem.

2 Klooni repo [https://github.com/e-gov/TARA-Mock](https://github.com/e-gov/TARA-Mock) masinasse.

3 Kui soovid, muuda etteantud identiteete

TARA-Mock-is määrab kasutaja ise identiteedi (isikukoodi, ees- ja perekonnanime), millega ta autenditakse. Selleks ta kas valib etteantud identiteetide hulgast või sisestab ise identiteeti.

Tarkvaraga on kaasas 3 etteantud identiteeti. Etteantud identiteetide muutmiseks redigeeri faili `service/data.go`:

```
...
		Identity{"Isikukood1", "Eesnimi1", "Perekonnanimi1"},
		Identity{"Isikukood2", "Eesnimi2", "Perekonnanimi2"},
		Identity{"Isikukood3", "Eesnimi3", "Perekonnanimi3"},
...
```
Muudatusi saab teha ka hiljem. Siis tuleb TARA-Mock-i uuesti käivitada.

4 Kontrolli ja vajadusel muuda TARA-Mock-is seadistatud hostinimesid ja porte. Vaikeseadistus on tehtud eeldustel:

- TARA-Mock töötab arendaja masinas (`localhost`), pordil `8080`
- rakendus, millest TARA-Mock-i poole pöördutakse (klientrakendus), töötab arendaja masinas (`localhost`), pordil `8081`
- klientrakendusse tagasisuunamise URL on `https://localhost:8081/return`.

Muuda failis `service/main.go` olev vaikeseadistus oma konfiguratsioonile vastavaks:

```
const (
	taraMockHost       = "localhost"
	returnURL          = "https://localhost:8081/return"
	httpServerPort     = ":8080"
...
```

5 Valmista ja paigalda võtmed ja serdid, vt [Serdid](docs/Serdid.md)

  - sh lisa sirvikusse TARA-Mock-i CA sert

6 Käivita TARA-Mock:

```
cd service
go run .
```

TARA-Mock on klientrakenduse teenindamiseks valmis.

## Klientrakenduse näidis

Klientrakenduse näidis töötab lokaalses masinas, pordil 8081. Otspunktid:

- `/health` - elutukse
- `/` - avaleht; kasutaja saab sealt minna TARA-Mock-i autentima
- `/login` - kasutaja suunamine TARA-Mock-i autentima
- `/return` - autentimiselt tagasi suunatud kasutaja vastuvõtmine, identsustõendi pärimine TARA-Mock-st ja sisselogimise lõpuleviimine 

TARA-Mock-ga kaasasoleva klientrakenduse kasutamiseks:

- Masinas peab olema paigaldatud Go, versioon 1.11 või hilisem.
- Klooni repo [https://github.com/e-gov/TARA-Mock](https://github.com/e-gov/TARA-Mock) masinasse.
- Valmista ja paigalda võtmed ja serdid, vt [Serdid](docs/Serdid.md)
- Vajadusel sea klientrakenduses õige TARA-Mock-i hostinimi (vaikimisi `localhost`), TARA-Mock-i ja klientrakenduse pordinumbrid (vaikimisi `8080` ja `8081`)
- Käivita klientrakendus:

```
cd client
go run .
```

- Ava sirvikus klientrakenduse avaleht (vaikimisi `https://localhost:8081`)

## TO DO

- Võib teostada UserInfo otspunkti, kui selleks on nõudlus. TARA pakub UserInfo otspunkti, kuid selle kasutamine ei ole soovitatav. Kõik vajalikud andmed saab kätte juba identsustõendist. 
