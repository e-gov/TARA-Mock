# TARA-Mock

TARA-Mock on rakendus, mis etendab TARA autentimist. 

[Ülevaade](#ülevaade)
[Kasutusstsenaarium](#kasutusstsenaarium)
[Lihtsustused](#lihtsustused)
[Paigaldamine](#paigaldamine)
[Klientrakenduse näidis](#klientrakenduse-näidis)

## Ülevaade

TARA-Mock on mõeldud kasutamiseks siis, kui [TARA testteenuse](https://e-gov.github.io/TARA-Doku/Testimine) võimalused jäävad klientrakenduse funktsionaalsuste testimisel napiks. TARA testteenuse abil saab autentida väga väikese hulga TARA poolt ette antud testkasutajatega.

TARA-Mock seevastu võimaldab kasutajal valida autentimise dialoogis identiteet  TARA-Mock seadistuses etteantud identiteetide hulgast või sisestada ise isikukood, ees- ja perekonnanimi, mille all soovitakse sisse logida. Sisuliselt saab sisse logida ükskõik millise identiteediga.

TARA-Mock on reaalse TARA-ga ühilduv, s.t TARA-Mock ja TARA on üksteisega vahetatavad lihtsa seadistusega.

TARA-Mock juures on ka klientrakenduse näidis.

TARA-Mock ei ole mõeldud kasutamiseks toodangus.

TARA-Mock ei ole mõeldud ka TARA-ga liidestamise testimiseks, sest TARA-Mock-is on ära jäetud mitmeid toodangus vajalikke kontrolle (vt allpool). TARA-ga liidestamise testimiseks on [TARA testteenus](https://e-gov.github.io/TARA-Doku/Testimine).

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
- klientrakenduse eelnev registreerimine ei ole nõutav
- turvaelemendid (`state` ja `nonce`), samuti `return_uri`, antakse ühes sammus edasi HTML vormi peidetud väljadena
- puudub päringuvõltsimise kaitse (CSRF). Märkus: Kuna TARA-Mock-s saab sisse logida suvalise kasutajana, siis ei ole kaitsel ka mõtet.
- juhusõned genereeritakse tavalise (`math/rand`), mitte krüptograafilise juhuarvugeneraatoriga (`crypto/rand`)
- minimaalne logimine; TARA-Mock väljastab mõningast logiteavet konsoolile
- klientrakenduse salasõna ei kontrollita
- parameetreid `scope` ja `response_type` ei kontrollita
- parameetrid `ui_locales` ei kontrollita ega toetata; TARA-Mock-i kasutajaliides on eesti keeles
- parameetrit `acr_values` ei kontrollita; identsustõend väljastatakse alati väite (_claim_) `acr` (tagatistase) väärtusega `high`
- identsustõendi väljastamisel `return_uri` ei kontrollita; identsustõend väljastatakse ainult volituskoodi alusel
- identsustõendi väljastamisel ei kontrollita, kas tõend on aegunud; tõendile järeletulemise aeg ei ole piiratud
- identsustõendite hoidlat ei puhastata aegunud tõenditest
- ei kontrollita, et identsustõend väljastatakse ainult üks kord
- isikukoodi vastavust Eesti isikukoodi standardile ei kontrollita; kui `date_of_birth` väärtust ei saa isikukoodist moodustada, siis tagastatakse väärtus `1961-07-12`
- ei kontrollita, kas isik on elus või üldse olemas
- ei kontrollita sisestatud nimede vastavust keelenormile
- autentimismeetodina näidatakse alati `mID`
- TARA-Mock-is ei ole teostatud UserInfo otspunkt (autenditud kasutaja andmete küsimine pääsutõendiga (_access token_)). TARA pakub UserInfo otspunkti, kuid selle kasutamine ei ole soovitatav. Kõik vajalikud andmed saab kätte juba identsustõendist


Mida siis kontrollitakse?

- `state` ja `nonce` peegeldatakse tagasi, nii nagu OIDC protokoll ette näeb
- `return_uri` peab olema kehtiv - muidu ei jõua kasutaja rakendusse tagasi
- vormikohaselt täidetakse kogu TARA kasutusvoog (v.a UserInfo otspunkt)
- identsustõend allkirjastatakse, allkirja kontrollimise võti on võtmeotspunktis


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

Muuda failis `service/main.go` olev vaikeseadistus oma konfiguratsioonile vastavaks:

```
const (
	taraMockHost       = "localhost"
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

TARA-Mock-ga kaasasolev klientrakenduse näidis pakub otspunkte:

- `/health` - elutukse
- `/` - avaleht; kasutaja saab sealt minna TARA-Mock-i autentima
- `/login` - kasutaja suunamine TARA-Mock-i autentima
- `/return` - autentimiselt tagasi suunatud kasutaja vastuvõtmine, identsustõendi pärimine TARA-Mock-st ja sisselogimise lõpuleviimine 

Klientrakenduse kasutamiseks:

1 Masinas peab olema paigaldatud Go, versioon 1.11 või hilisem.

2 Klooni repo [https://github.com/e-gov/TARA-Mock](https://github.com/e-gov/TARA-Mock) masinasse.

3 Kontrolli ja vajadusel muuda vaikeseadistus vastavaks oma konfiguratsioonile. Vaikeseadistus on tehtud eeldustel:

- TARA-Mock-i hostinimi on `localhost` ja port on `8080`
- klientrakenduse hostinimi on `localhost` ja port on `8081`

4 Valmista ja paigalda võtmed ja serdid, vt [Serdid](docs/Serdid.md)

5 Käivita klientrakendus:

```
cd client
go run .
```

6 Ava sirvikus klientrakenduse avaleht (vaikimisi `https://localhost:8081`)
 
