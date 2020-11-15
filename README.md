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

![Screen](/doc/kursnbp.png)
