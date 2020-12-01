# kursNBP


[![GitHub](https://img.shields.io/github/license/pjaskulski/kursnbp)](https://opensource.org/licenses/MIT) 
![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/pjaskulski/kursnbp?include_prereleases) 
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/pjaskulski/kursnbp) 
[![go report](https://goreportcard.com/badge/github.com/pjaskulski/kursnbp)](https://goreportcard.com/report/github.com/pjaskulski/kursnbp) 

**[Description (en)](#description)**<br>
**[Opis (pl)](#opis)**<br>
**[Screenshots](#screenshots)**<br>
**[TODO](#todo)**<br>
**[Contact](#contact)**<br>

## Description:

kursNBP - a command line tool for downloading exchange rates and gold prices from the website of the National Bank of Poland

The project uses the nbpapi library: [https://github.com/pjaskulski/nbpapi](https://github.com/pjaskulski/nbpapi)

Downloads (v0.3.2):<br> 
[linux](https://github.com/pjaskulski/kursnbp/releases/download/v0.3.2/kursnbp_v0.3.2_linux_amd64.tar.gz) 
[windows](https://github.com/pjaskulski/kursnbp/releases/download/v0.3.2/kursnbp_v0.3.2_windows_amd64.zip) 
[macos](https://github.com/pjaskulski/kursnbp/releases/download/v0.3.2/kursnbp_v0.3.2_macos_amd64.tar.gz) 
[FreeBSD](https://github.com/pjaskulski/kursnbp/releases/download/v0.3.2/kursnbp_v0.3.2_freebsd_amd64.tar.gz)

    Usage:
      kursnbp table | currency | gold [--flag]

    Commands:
      table       returns a table of exchange rates (or a series of tables)
      currency    returns the rate of the specified currency or a series of 
                  rates
      gold        returns the price of gold (or series), the price of 1g of
                  gold (of 1000 millesimal fineness) calculated at NBP

    Global Flags:
         --version    Displays the version of the program
      -h --help       Displays help describing program commands and flags
      -b --clipboard  Copy output to clipboard instead of printing 
                      (Linux/FreeBSD: requires 'xclip' or 'xsel' command 
                      to be installed)
  
    Flags for commands:
      table:
        -t --table    Table type: A, B or C, default: A
        -d --date     Date in the format: 'YYYY-MM-DD' (ISO 8601 standard), 
                      or a range of dates in the format: 
                      'YYYY-MM-DD:YYYY-MM-DD' or 'today' (rate for today) 
                      or 'current' - current table / rate (last published),
                      default: current
        -l --last     As an alternative to --date, the last <n> tables/rates
                      can be retrieved e.g. -l = 5, default: 0
        -o --output   Output format: table (formatted text table),
                      json (json format), csv (data with comma separated 
                      fields, field names in the first line), xml (xml 
                      format), default: table
        -i --lang     Output language (also error messages), allowed
                      values: 'en', 'pl', default: 'en', currency names 
                      returned by the NBP service are always in Polish
    
      currency:
        -t --table    Table type: A, B or C, default: A
        -d --date     Date in the format: 'YYYY-MM-DD' (ISO 8601 standard), 
                      or a range of dates in format: 'YYYY-MM-DD:YYYY-MM-DD' 
                      or 'today' (rate for today) or 'current' - current 
                      table/rate (last published), default: current
        -l --last     As an alternative to --date, the last <n> tables/rates
                      can be retrieved e.g. -l = 5, default: 0
        -c --code     ISO 4217 currency code, depending on the type of the 
                      table available currencies may vary
        -o --output   Output format: table (formatted text table),
                      json (json format), csv (data with comma separated 
                      fields, field names in the first line), xml (xml 
                      format), default: table
        -i --lang     Output language (also error messages), allowed
                      values: 'en', 'pl', default: 'en', currency names 
                      returned by the NBP service are always in Polish

      gold:
        -d --date    Date in the format: 'YYYY-MM-DD' (ISO 8601 standard), 
                     or a range of dates in format: 'YYYY-MM-DD: YYYY-MM-DD' 
                     or 'today' (today's price) or 'current' - current 
                     price (last published), default: current
        -l --last    As an alternative to --date, the last <n> gold prices 
                     can be retrieved e.g. -last=5, default: 0
        -o --output  Output format: table (formatted text table),
                     json (json format), csv (data with comma separated 
                     fields, field names in the first line), xml (xml 
                     format), default: table
        -i --lang    Output language (also error messages), allowed
                     values: 'en', 'pl', default: 'en', currency names 
                     returned by the NBP service are always in Polish

Examples:
    
    kursnbp --help
    Displays help describing program commands and global flags

    kursnbp --help currency
    Displays help for 'currency' command
    
    kursnbp table
    Displays the current table of exchange rate of type A
    
    kursnbp table --last=2 --table=C
    Series of latest 2 tables of exchange rates of type C

    kursnbp table --date=2020-11-19 --table=A
    Exchange rate table of type A published on 19 November 2020

    kursnbp table --date=today --output=csv
    Exchange rate table of type A published today,  in csv format

    kursnbp table --date=today --output=csv --clipboard
    Exchange rate table of type A published today,  in csv format, 
    copied to the clipboard

    kursnbp table --date=today --output=xml
    Exchange rate table of type A published today, in xml format

    kursnbp currency --code=CHF
    Current exchange rate CHF (Swiss franc) from the exchange rate 
    table of type A

    kursnbp currency --code=EUR --last=10
    Series of latest 10 exchange rates of currency EUR (euro) 
    from the exchange rate table of type A

    kursnbp gold
    Current price of gold

    kursnbp gold --date=2020-11-12:2020-11-19
    Series of gold prices published from 12 November 2020 
    to 19 November 2020

Documentation of the API service of the National Bank of Poland
[http://api.nbp.pl/en.html](http://api.nbp.pl/en.html)

License: MIT <br>
Modules used: <br>
https://github.com/integrii/flaggy (Unlicense License)
https://github.com/atotto/clipboard (BSD-3-Clause License)


## Opis:

kursNBP - konsolowy program do pobierania kursów walut i notowań cen złota z serwisu Narodowego Banku Polskiego

Wykorzystano bibliotekę nbpapi: [https://github.com/pjaskulski/nbpapi](https://github.com/pjaskulski/nbpapi)

Do pobrania (wersja 0.3.2):<br> 
[linux](https://github.com/pjaskulski/kursnbp/releases/download/v0.3.2/kursnbp_v0.3.2_linux_amd64.tar.gz) 
[windows](https://github.com/pjaskulski/kursnbp/releases/download/v0.3.2/kursnbp_v0.3.2_windows_amd64.zip) 
[macos](https://github.com/pjaskulski/kursnbp/releases/download/v0.3.2/kursnbp_v0.3.2_macos_amd64.tar.gz) 
[FreeBSD](https://github.com/pjaskulski/kursnbp/releases/download/v0.3.2/kursnbp_v0.3.2_freebsd_amd64.tar.gz)

    Użycie:
      kursnbp table|currency|gold [--flag]

    Polecenia: 
      table      Zwraca tabelę kursów wymiany walut (lub serię tabel)
      currency   Zwraca kurs wskazanej waluty lub serię kursów
      gold       Zwraca cenę złota lub serię notowań cen złota (cena 1 g 
                 złota, w próbie 1000)

    Flagi globalne: 
         --version    Wyświetla wersję programu
      -h --help       Wyświetla pomoc z opisem poleceń i flag programu 
      -b --clipboard  Kopiuje dane wyjściowe do schowka, zamiast drukowania
                      na ekranie (Linux/FreeBSD: dodatkowe programy 'xsel'
                      lub 'xclip' muszą być zainstalowane w systemie)
  
    Flagi:
      table: 
        -t --table   Typ tabeli kursów: A, B lub C, domyślnie: A
        -d --date    Data tabeli kursów w formacie: 'RRRR-MM-DD' (standard 
                     ISO 8601), lub zakres dat 'RRRR-MM-DD:RRRR-MM-DD' 
                     lub 'today' (kurs na dziś) lub 'current' - bieżąca 
                     tabela/kurs (ostatnio opublikowany) domyślnie: current
        -l --last    Alternatywnie do --date można pobrać ostatnich <n> 
                     tabel/kursów np. -l=5, domyślnie: 0
        -o --output  Format danych wyjściowych: table (sformatowana tabela 
                     tekstowa), json (format json), csv (dane z polami 
                     rozdzielanymi przecinkiem, nazwy pól w pierwszym 
                     wierszu), xml (format xml), domyślnie: table
        -i --lang    Język danych wyjściowych (także komunikatów 
                     o błędach), nazwy walut zwracane przez serwis NBP
                     zawsze w języku polskim, dozwolone wartości: 'en', 'pl', 
                     domyślnie 'en'
    
      currency:
        -t --table   Typ tabeli kursów: A, B lub C, domyślnie: A
        -d --date    Data tabeli kursów w formacie: 'RRRR-MM-DD' (standard 
                     ISO 8601), lub zakres dat 'RRRR-MM-DD:RRRR-MM-DD' 
                     lub 'today' (kurs na dziś) lub 'current' - bieżąca 
                     tabela/kurs (ostatnio opublikowany) domyślnie: current
        -l --last    Alternatywnie do --date można pobrać ostatnich <n> 
                     tabel/kursów np. -l=5, domyślnie: 0
        -c --code    Kod waluty w standardzie ISO 4217, zależnie od typu 
                     tabeli lista dostępnych walut może się różnić
        -o --output  Format danych wyjściowych: table (sformatowana tabela 
                     tekstowa), json (format json), csv (dane z polami 
                     rozdzielanymi przecinkiem, nazwy pól w pierwszym 
                     wierszu), xml (format xml), domyślnie: table
        -i --lang    Język danych wyjściowych (także komunikatów 
                     o błędach), nazwy walut zwracane przez serwis NBP
                     zawsze w języku polskim, dozwolone wartości: 'en', 'pl', 
                     domyślnie 'en'

      gold:
        -d --date    Data notowania ceny złota w formacie: 'RRRR-MM-DD', 
                     (standard ISO 8601) lub zakres dat: 
                     'RRRR-MM-DD:RRRR-MM-DD' lub 'today' (cena na dziś) lub 
                     'current' - bieżąca cena (ostatnio opublikowana) 
                     domyślnie: current
        -l --last    Alternatywnie do --date można pobrać ostatnich <n> cen 
                     złota np. -l=5, domyślnie: 0
        -o --output  Format danych wyjściowych: table (sformatowana tabela 
                     tekstowa), json (format json), csv (dane z polami 
                     rozdzielanymi przecinkiem, nazwy pól w pierwszym 
                     wierszu), xml (format xml), domyślnie: table
        -i --lang    Język danych wyjściowych (także komunikatów 
                     o błędach), nazwy walut zwracane przez serwis NBP
                     zawsze w języku polskim, dozwolone wartości: 'en', 'pl', 
                     domyślnie 'en'

Przykłady:
    
    kursnbp --help
    Wyświetla ogólny help do programu

    kursnbp --help gold
    Wyświetla szczegółowy help do polecenia gold

    kursnbp table
    Wyświetla bieżącą tabelę kursów typu A
    
    kursnbp table --last=2 --table=C
    Wyświetla 2 ostatnie tabele kursów typu C

    kursnbp table --date=2020-11-19 --table=A
    Wyświetla tabelę kursów walut z podanego dnia

    kursnbp table --date=today --output=csv
    Wyświetla dzisiejszą tabelę kursów w formacie csv

    kursnbp table --date=today --output=csv --clipboard
    Pobiera dzisiejszą tabelę kursów w formacie csv 
    i kopiuje do schowka
    
    kursnbp table --date=today --output=xml
    Wyświetla dzisiejszą tabelę kursów w formacie xml

    kursnbp currency --code=CHF
    Wyświetla bieżący kurs waluty CHF (frank szwajcarski)

    kursnbp currency --code=EUR --last=10
    Wyświetla 10 ostatnich kursów waluty EUR (euro)

    kursnbp gold
    Wyświetla bieżącą cenę złota

    kursnbp gold --date=2020-11-12:2020-11-19
    Wyświetla listę notowań cen złota w podanym przedziale dat

Dokumentacja serwisu API Narodowego Banku Polskiego: [http://api.nbp.pl/](http://api.nbp.pl/)

Licencja: MIT<br>
Wykorzystano moduły:<br> 
https://github.com/integrii/flaggy (Unlicense License)
https://github.com/atotto/clipboard (BSD-3-Clause License)

## Screenshots:

![Screen](/docs/kursnbp_table.png)

![Screen](/docs/free_bsd_kursnbp.png)

![Screen](/docs/kursnbp_currency.png)

## TODO

-  [ ] more tests
-  [ ] better tables in terminal


## Contact
- Please use [Github issue tracker](https://github.com/pjaskulski/kursnbp/issues) for filing bugs or feature requests.