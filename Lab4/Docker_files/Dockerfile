# zmiana na ubuntu
FROM ubuntu:latest

# ustawienie autora
LABEL maintainer="Adam Sidor <s101661@edu.pollub.pl>"

# aktualizacja i instalacja apache2
RUN apt-get update && \
    apt-get install -y apache2 && \
    apt-get clean

# kopiowanie strony do domyślnego katalogu serwera
COPY index.html /var/www/html/index.html

# otworzenie portu 80
EXPOSE 80

# uruchomienie serwera
CMD ["apachectl", "-D", "FOREGROUND"]