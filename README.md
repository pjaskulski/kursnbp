# kursNBP
Tekstowy klient pobierający kursy walut i ceny złota z serwisu Narodowego Banku Polskiego.

    Parametry wywołania programu:
      -day string
        data tabeli kursów (RRRR-MM-DD lub: today, current lub zakres 
        dat RRRR-MM-DD:RRRR-MM-DD)
      -last string
        seria ostatnich tabel kursów
      -output string
        format wyjścia (json, table) (default "json")
      -table string
        typ tabeli kursów (A, B lub C)
      -currency string
        kod waluty, domyślnie ALL - zwraca tabelę kursów dla wszystkich walut, 
        kod np. CHF zwraca kurs dla wybranej waluty (kod waluty w standardzie 
        ISO 4217)
        kod GOLD zwraca cenę złota

      np. 
      kursnbp -table=A -current -output=table
      wyświetla bieżącą tabelę kursów A w formie tabeli 
      (zob. zrzut ekranu poniżej)

      kursnbp -table=A -day=2020-11-12:2020-11-13 -output=table
      wyświetla 2 tabele zgodnie z zadanym zakresem dat 
      (12-13 listopad 2020) w formie tabeli

      kursnbp -currency=CHF -last=5 -table=A -output=table
      wyświetla listę ostatnich 5 kursów dla waluty CHF, w formie tabeli
      (zob. zrzut ekranu poniżej)

      kursnbp -currency=GOLD -day=current
      wyświetla aktualną cenę złota (cena 1 g złota, w próbie 1000)
      (zob. zrzut ekranu poniżej)

Dokumentacja serwisu API Narodowego Banku Polskiego: [http://api.nbp.pl/](http://api.nbp.pl/)

TODO:
  - więcej testów
  - format wyjścia CSV


## English version:

Console application for downloading exchange rates and gold prices from the website of the National Bank of Poland.

    Program calling parameters:
      -day string
        date of the exchange rate table (YYYY-MM-DD or: today, current 
        or date range YYYY-MM-DD: YYYY-MM-DD)
      -last string
        series of recent exchange rate tables
      -output string
        output format (json, table) (default "table")
      -table string
        exchange rate table type (A, B or C)
      -currency string
        currency code, by default =ALL - returns the table of rates for 
        all currencies, specific code e.g. CHF returns the rate for the 
        selected currency (currency code in the ISO 4217 standard),
        code GOLD returns gold price.

      e.g.
      kursnbp -table = A -current -output = table
      displays the current A table of exchange rates in table format 
      (see screenshot below)

      kursnbp -table = A -day = 2020-11-12: 2020-11-13
      displays 2 tables in accordance with a given date range (November 12-13, 2020) 
      in the form of a text table

      kursnbp -currency=EUR -last=5 -table=A
      displays a list of the last 5 CHF rates in a table format
      (see screenshot below)

      kursnbp -currency=GOLD -day=current
      displays the current gold price (1g of gold, of 1000 millesimal fineness)
      (see screenshot below)

Documentation of the API service of the National Bank of Poland
[http://api.nbp.pl/en.html](http://api.nbp.pl/en.html)


TODO:

  - more tests
  - output format CSV


![Screen](/doc/kursnbp.png)

![Screen](/doc/kursnbp2.png)

![Screen](/doc/kursnbp_gold.png)

![Screen](/doc/kursnbp_tabc.png)