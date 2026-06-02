# Sprawozdanie z Laboratorium 12 - Stack LEMP w Docker Compose

## 1. Konfiguracja pliku `docker-compose.yml`

Poniżej znajduje się pełna konfiguracja pliku `docker-compose.yml` realizująca uruchomienie stacka LEMP (Nginx, PHP-FPM, MySQL) wraz z phpMyAdmin:

```yaml
services:
  nginx:
    image: nginx:1.25-alpine
    container_name: nginx_server
    ports:
      - "4001:80"
    volumes:
      - ./nginx/default.conf:/etc/nginx/conf.d/default.conf:ro
      - ./html:/var/www/html:ro
    networks:
      - frontend
      - backend
    depends_on:
      - php

  php:
    image: php:8.2-fpm
    container_name: php_app
    volumes:
      - ./html:/var/www/html
    networks:
      - backend

  mysql:
    image: mysql:8.0
    container_name: mysql_db
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: test_db
      MYSQL_USER: test_user
      MYSQL_PASSWORD: test_password
    volumes:
      - mysql_data:/var/lib/mysql
    networks:
      - backend

  phpmyadmin:
    image: phpmyadmin:5.2
    container_name: phpmyadmin_ui
    ports:
      - "6001:80"
    environment:
      PMA_HOST: mysql
      MYSQL_ROOT_PASSWORD: root_password
    networks:
      - backend
    depends_on:
      - mysql

volumes:
  mysql_data:

networks:
  frontend:
  backend:
```

---

## 2. Uzasadnienie konfiguracji sieciowej (phpMyAdmin)

W architekturze zdefiniowano dwie izolowane sieci: `frontend` oraz `backend`. 

* **`phpmyadmin`** został celowo przyłączony **tylko i wyłącznie do sieci `backend`**.
* **Uzasadnienie:**
  * Kontener `phpmyadmin` musi komunikować się bezpośrednio z bazą danych `mysql_db` w celu jej administracji, dlatego wymagane jest, aby znajdował się w tej samej sieci wewnętrznej (`backend`).
  * Interfejs graficzny phpMyAdmin jest wystawiany na zewnątrz za pomocą bezpośredniego mapowania portu hosta (`6001:80`). Oznacza to, że użytkownik łączy się z nim bezpośrednio z systemu hosta, a ruch nie przechodzi przez serwer `nginx` (który jako jedyny spaja sieć `frontend` z `backend`). Przyłączenie phpMyAdmin do sieci `frontend` byłoby nadmiarowe i niezgodne z zasadą minimalnych uprawnień sieciowych.

---

## 3. Logi z uruchomienia i weryfikacji kontenerów (CLI)

### Krok 1: Uruchomienie całego stacka za pomocą Docker Compose

Poniżej log z pierwszego pobrania obrazów oraz uruchomienia wszystkich kontenerów:

```text
adam@Adams-MacBook-Pro Lab12 % docker compose up -d
[+] up 61/61
 ✔ Image php:8.2-fpm       Pulled                                                                                                                                                                                                                                      364.6s
 ✔ Image nginx:1.25-alpine Pulled                                                                                                                                                                                                                                      102.0s
 ✔ Image mysql:8.0         Pulled                                                                                                                                                                                                                                      361.7s
 ✔ Image phpmyadmin:5.2    Pulled                                                                                                                                                                                                                                      193.6s
 ✔ Network lab12_frontend  Created                                                                                                                                                                                                                                       0.0s
 ✔ Network lab12_backend   Created                                                                                                                                                                                                                                       0.0s
 ✔ Volume lab12_mysql_data Created                                                                                                                                                                                                                                       0.0s
 ✔ Container mysql_db      Started                                                                                                                                                                                                                                       0.4s
 ✔ Container php_app       Started                                                                                                                                                                                                                                       0.4s
 ✔ Container nginx_server  Started                                                                                                                                                                                                                                       0.3s
 ✔ Container phpmyadmin_ui Started                                                                                                                                                                                                                                       0.3s
```

### Krok 2: Sprawdzenie stanu działających usług

```text
adam@Adams-MacBook-Pro Lab12 % docker compose ps
NAME            IMAGE               COMMAND                  SERVICE      CREATED         STATUS         PORTS
mysql_db        mysql:8.0           "docker-entrypoint.s…"   mysql        3 minutes ago   Up 3 minutes   3306/tcp, 33060/tcp
nginx_server    nginx:1.25-alpine   "/docker-entrypoint.…"   nginx        3 minutes ago   Up 3 minutes   0.0.0.0:4001->80/tcp, [::]:4001->80/tcp
php_app         php:8.2-fpm         "docker-php-entrypoi…"   php          3 minutes ago   Up 3 minutes   9000/tcp
phpmyadmin_ui   phpmyadmin:5.2      "/docker-entrypoint.…"   phpmyadmin   3 minutes ago   Up 3 minutes   0.0.0.0:6001->80/tcp, [::]:6001->80/tcp
```

### Krok 3: Inspekcja przydziału sieciowego kontenerów

Weryfikacja podłączenia kontenerów do sieci `frontend` oraz `backend`:

```text
adam@Adams-MacBook-Pro Lab12 % docker inspect lab12_frontend | jq '.[].Containers'
{
  "9c0c151a4a309144d7a898fc9cef97eff882a600cf3992cf6f562a9bbbb79bd5": {
    "Name": "nginx_server",
    "EndpointID": "e8209c74f12082af7cc51bc16927b9e40342ddef663a0b004021b5139d966bc0",
    "MacAddress": "d2:be:c6:5b:9d:31",
    "IPv4Address": "172.19.0.2/16",
    "IPv6Address": ""
  }
}

adam@Adams-MacBook-Pro Lab12 % docker inspect lab12_backend | jq '.[].Containers'
{
  "223ec2c3e570d03ef061307b58cfcdf00e88d256df40f4a8825fdacc31b8fa21": {
    "Name": "php_app",
    "EndpointID": "a02cceba48b4c55c454a9e2c4d94d3fb6386ea43a866cdab17a1a6c16c30a90f",
    "MacAddress": "7e:58:0e:9b:08:08",
    "IPv4Address": "172.20.0.2/16",
    "IPv6Address": ""
  },
  "5b6e1c90b61510bf217468540f2beadf8dc345841e632b4e49419220bedbee09": {
    "Name": "mysql_db",
    "EndpointID": "7f9f5bce77473d30c60b49195948dfab4cb937cf5035ae9c25f36bed7155b6ec",
    "MacAddress": "fa:db:2b:9e:8c:38",
    "IPv4Address": "172.20.0.3/16",
    "IPv6Address": ""
  },
  "9c0c151a4a309144d7a898fc9cef97eff882a600cf3992cf6f562a9bbbb79bd5": {
    "Name": "nginx_server",
    "EndpointID": "a9f09aff220c140e7ebd7428de90fe8421269fec5cbeb9e1d0a6fa5485d4387c",
    "MacAddress": "92:c0:c7:b8:16:6c",
    "IPv4Address": "172.20.0.5/16",
    "IPv6Address": ""
  },
  "d640aa85b76cf9e2bbd0ab000445640077bb33b435f94eff481f091d3c32ad2c": {
    "Name": "phpmyadmin_ui",
    "EndpointID": "d0e588ea8950b8fc7ee487c4758b23a0685039031714db1a0242c751b50e5047",
    "MacAddress": "f6:df:1c:36:af:2b",
    "IPv4Address": "172.20.0.4/16",
    "IPv6Address": ""
  }
}
```

---

## 4. Dowód poprawnego działania stacka LEMP i bazy danych

### Dowód 1: Poprawne wyświetlanie strony PHP i statusu połączenia z MySQL

Poniższy zrzut ekranu przedstawia poprawnie uruchomioną stronę `index.php` pod adresem `http://localhost:4001`, co potwierdza poprawne działanie Nginx, interpretera PHP oraz nawiązanie połączenia z bazą danych MySQL:

![Strona startowa php](img/index_php.png)

---

### Dowód 2: Możliwość zainicjowania testowej bazy danych

Poniższy zrzut ekranu przedstawia panel administracyjny phpMyAdmin pod adresem `http://localhost:6001` z widoczną nowo utworzoną, testową bazą danych, co potwierdza pełną operacyjność środowiska bazodanowego:

![Panel phpMyAdmin](img/phpmyadmin.png)

