# kursNBP
Tekstowy klient pobierający kursy walut z serwisu Narodowego Banku Polskiego

    Parametry wywołania programu:
      -day string
        data tabeli kursów (RRRR-MM-DD lub: today, current lub zakres dat RRRR-MM-DD:RRRR-MM-DD)
      -last string
        seria ostatnich tabel kursów
      -output string
        format wyjścia (json, table) (default "json")
      -table string
        typ tabeli kursów (A, B lub C)

      np. 
      kursnbp -table=A -current -output=table
      wyświetla bieżącą tabelę kursów A w formie tabeli (zob. zrzut ekranu poniżej)

      kursnbp -table=A -day=2020-11-12:2020-11-13 -output=table
      wyświetla 2 tabele zgodnie z zadanym zakresem dat (12-13 listopad 2020) w formie tabeli

Dokumentacja serwisu API Narodowego Banku Polskiego: [http://api.nbp.pl/](http://api.nbp.pl/).

TODO:
  - testy
  - pobieranie kursów wskazanej dla waluty
  - pobieranie kursów złota

![Screen](/doc/kursnbp.png)
