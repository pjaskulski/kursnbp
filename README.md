# kursNBP
Tekstowy klient pobierający kursy walut z serwisu Narodowego Banku Polskiego.

    Parametry wywołania programu:
      -day string
        data tabeli kursów (RRRR-MM-DD lub: today, current lub zakres dat RRRR-MM-DD:RRRR-MM-DD)
      -last string
        seria ostatnich tabel kursów
      -output string
        format wyjścia (json, table) (default "json")
      -table string
        typ tabeli kursów (A, B lub C)
      -currency string
        kod waluty, domyślnie ALL - zwraca tabelę kursów dla wszystkich walut, kod np. CHF zwraca kurs dla wybranej waluty (kod waluty w standardzie ISO 4217)

      np. 
      kursnbp -table=A -current -output=table
      wyświetla bieżącą tabelę kursów A w formie tabeli (zob. zrzut ekranu poniżej)

      kursnbp -table=A -day=2020-11-12:2020-11-13 -output=table
      wyświetla 2 tabele zgodnie z zadanym zakresem dat (12-13 listopad 2020) w formie tabeli

      kursnbp -currency=CHF -last=5 -table=A -output=table
      wyświetla listę ostatnich 5 kursów dla waluty CHF, w formie tabeli

Dokumentacja serwisu API Narodowego Banku Polskiego: [http://api.nbp.pl/](http://api.nbp.pl/).

TODO:
  - testy
  - pobieranie kursów złota


## English version:

Text client for downloading exchange rates from the website of the National Bank of Poland.

    Program calling parameters:
      -day string
        date of the exchange rate table (YYYY-MM-DD or: today, current or date range YYYY-MM-DD: YYYY-MM-DD)
      -last string
        series of recent exchange rate tables
      -output string
        output format (json, table) (default "json")
      -table string
        exchange rate table type (A, B or C)
      -currency string
        currency code, by default =ALL - returns the table of rates for all currencies, specific code e.g. CHF returns the rate for the selected currency (currency code in the ISO 4217 standard)

      e.g.
      kursnbp -table = A -current -output = table
      displays the current A table of exchange rates in table format (see screenshot below)

      kursnbp -table = A -day = 2020-11-12: 2020-11-13 -output = table
      displays 2 tables in accordance with a given date range (November 12-13, 2020) in the form of a text table

      kursnbp -currency=CHF -last=5 -table=A -output=table
      displays a list of the last 5 CHF rates in a table format


Documentation of the API service of the National Bank of Poland
[http://api.nbp.pl/en.html](http://api.nbp.pl/en.html)


TODO:

  - tests
  - downloading gold rates


![Screen](/doc/kursnbp.png)

![Screen](/doc/kursnbp2.png)