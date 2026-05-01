# Zadanie 1 - Część Nieobowiązkowa

## 1. Analiza bezpieczeństwa obrazu podstawowego

Prace rozpoczęto od weryfikacji obrazu `zadanie1` zbudowanego w części obowiązkowej. Wynik skanowania został zapisany do pliku logu.

**Polecenie:**
```bash
docker scout cves zadanie1 > nieobowiazkowe/scout_report_glowny.log
```

**Wynik polecenia:**
```text
    ...Storing image for indexing
    ✓ Image stored for indexing
    ...Indexing
    ✓ Indexed 2 packages
    ✓ No vulnerable package detected
```

**Zawartość pliku `scout_report_glowny.log`:**
```text
## Overview

                   │       Analyzed Image        
───────────────────┼─────────────────────────────
 Target            │  zadanie1:latest            
   digest          │  6ca09d746ef3               
   platform        │ linux/arm64                 
   vulnerabilities │    0C     0H     0M     0L  
   size            │ 3.8 MB                      
   packages        │ 2                           


## Packages and Vulnerabilities

  No vulnerable packages detected
```

Zawartość pliku pojazuje ze w obrazie bazowym nie wystepuja podatnosci.

---

## 2. Zmodyfikowany Dockerfile (GitHub & SSH)

Zmodyfikowano konfigurację Dockerfile, która realizuje pobieranie kodu bezpośrednio z repozytorium GitHub przy użyciu funkcjonalności `mount=ssh`.

```dockerfile
# syntax=docker/dockerfile:1
# ETAP 1: Budowanie
FROM golang:1.26-alpine AS builder

# Instalacja narzędzi do pobierania kodu przez SSH i konfiguracja znanych hostów
RUN apk add --no-cache git openssh-client
RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

WORKDIR /app

# Pobieranie kodu z repozytorium GitHub przy użyciu funkcjonalności mount=ssh
RUN --mount=type=ssh git clone git@github.com:Adam-Sidor/Docker_lab.git .

WORKDIR /app/Zadanie1

# Pobieranie zależności i kompilacja aplikacji
RUN go mod download || true
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o weather-app .

# ETAP 2: Finalny obraz 
FROM scratch

# Dane zgodnie ze standardem OCI
LABEL org.opencontainers.image.authors="Adam Sidor"

# Kopiowanie certyfikatów SSL i skompilowanej aplikacji z etapu builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/Zadanie1/weather-app /weather-app

# Informacja o porcie
EXPOSE 8080

# Healthcheck wykorzystujący mechanizm flag wbudowany w aplikację
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/weather-app", "-health"]

# Uruchomienie aplikacji
ENTRYPOINT ["/weather-app"]
```

---

## 3. Budowanie wieloplatformowe (Pierwszy przebieg)

Skonfigurowano builder korzystający ze sterownika `docker-container` i wykonano pierwsze budowanie na dwie platformy (`amd64`, `arm64`) wraz z eksportem cache.

**Tworzenie buildera i ustawienie go jako aktualnie używanego:**
```bash
docker buildx create --name docker-container-builder --driver docker-container --use
```
**Budowanie na obie platformy i eksport cache:**
```bash
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t adamsidor/zadanie1:nieobowiazkowe \
  -f nieobowiazkowe/dockerfile \
  --ssh default \
  --push \
  --cache-to type=registry,ref=adamsidor/zadanie1:cache,mode=max \
  --cache-from type=registry,ref=adamsidor/zadanie1:cache \
  .
```

**Pełny log budowania:**
```text
[+] Building 982.7s (35/35) FINISHED                                                                                                                                                                                               docker-container:docker-container-builder
 => [internal] load build definition from dockerfile                                                                                                                                                                                                                    0.0s
 => => transferring dockerfile: 857B                                                                                                                                                                                                                                    0.0s
 => resolve image config for docker-image://docker.io/docker/dockerfile:1                                                                                                                                                                                               1.0s
 => [auth] docker/dockerfile:pull token for registry-1.docker.io                                                                                                                                                                                                        0.0s
 => CACHED docker-image://docker.io/docker/dockerfile:1@sha256:2780b5c3bab67f1f76c781860de469442999ed1a0d7992a5efdf2cffc0e3d769                                                                                                                                         0.0s
 => => resolve docker.io/docker/dockerfile:1@sha256:2780b5c3bab67f1f76c781860de469442999ed1a0d7992a5efdf2cffc0e3d769                                                                                                                                                    0.0s
 => [linux/arm64 internal] load metadata for docker.io/library/golang:1.26-alpine                                                                                                                                                                                       0.4s
 => [linux/amd64 internal] load metadata for docker.io/library/golang:1.26-alpine                                                                                                                                                                                       0.6s
 => [auth] library/golang:pull token for registry-1.docker.io                                                                                                                                                                                                           0.0s
 => [internal] load .dockerignore                                                                                                                                                                                                                                       0.0s
 => => transferring context: 2B                                                                                                                                                                                                                                         0.0s
 => ERROR importing cache manifest from adamsidor/zadanie1:cache                                                                                                                                                                                                        0.6s
 => [linux/amd64 builder 1/8] FROM docker.io/library/golang:1.26-alpine@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1                                                                                                                         0.0s
 => => resolve docker.io/library/golang:1.26-alpine@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1                                                                                                                                             0.0s
 => [linux/arm64 builder 1/8] FROM docker.io/library/golang:1.26-alpine@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1                                                                                                                         0.0s
 => => resolve docker.io/library/golang:1.26-alpine@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1                                                                                                                                             0.0s
 => [auth] adamsidor/zadanie1:pull token for registry-1.docker.io                                                                                                                                                                                                       0.0s
 => [linux/arm64 builder 2/8] RUN apk add --no-cache git openssh-client                                                                                                                                                                                          0.5s
 => [linux/arm64 builder 3/8] RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts                                                                                                                                                        0.8s
 => [linux/arm64 builder 4/8] WORKDIR /app                                                                                                                                                                                                                       0.1s
 => [linux/arm64 builder 5/8] RUN --mount=type=ssh git clone git@github.com:Adam-Sidor/Docker_lab.git .                                                                                                                                                                 7.8s
 => [linux/amd64 builder 2/8] RUN apk add --no-cache git openssh-client                                                                                                                                                                                          0.5s
 => [linux/amd64 builder 3/8] RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts                                                                                                                                                        0.8s
 => [linux/amd64 builder 4/8] WORKDIR /app                                                                                                                                                                                                                       0.1s
 => [linux/amd64 builder 5/8] RUN --mount=type=ssh git clone git@github.com:Adam-Sidor/Docker_lab.git .                                                                                                                                                                 7.1s
 => [linux/amd64 builder 6/8] WORKDIR /app/Zadanie1                                                                                                                                                                                                                     0.0s
 => [linux/amd64 builder 7/8] RUN go mod download || true                                                                                                                                                                                                               0.1s 
 => [linux/amd64 builder 8/8] RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o weather-app .                                                                                                                                                                  25.1s
 => [linux/arm64 builder 6/8] WORKDIR /app/Zadanie1                                                                                                                                                                                                                     0.0s 
 => [linux/arm64 builder 7/8] RUN go mod download || true                                                                                                                                                                                                               0.1s
 => [linux/arm64 builder 8/8] RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o weather-app .                                                                                                                                                                   5.5s
 => [linux/arm64 stage-1 1/2] COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/                                                                                                                                                             0.1s
 => [linux/arm64 stage-1 2/2] COPY --from=builder /app/Zadanie1/weather-app /weather-app                                                                                                                                                                         0.1s
 => [linux/amd64 stage-1 1/2] COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/                                                                                                                                                             0.1s
 => [linux/amd64 stage-1 2/2] COPY --from=builder /app/Zadanie1/weather-app /weather-app                                                                                                                                                                         0.1s
 => exporting to image                                                                                                                                                                                                                                                  6.2s
 => => exporting layers                                                                                                                                                                                                                                                 0.0s
 => => exporting manifest sha256:e244713cccf6d6614e3185906ed6dfce042de65db5c3176cf97add6d92445df7                                                                                                                                                                       0.0s
 => => exporting config sha256:91f063c5080e71af9b6da7b2e24f63a9cdc71e5396bf65ecda5e0d28a65b43c7                                                                                                                                                                         0.0s
 => => exporting attestation manifest sha256:fbc643dbf18b69e817c9169e27d3f002096b327a80ac3504b0cdfa321d8f7fdd                                                                                                                                                           0.0s
 => => exporting manifest sha256:92404d8c4cd88fd478a845b3d1e0fe789e5819234623bd6fffc57b58b21fd851                                                                                                                                                                       0.0s
 => => exporting config sha256:e7e921527d86221c9823abcea30b2feb4edd29fdf5aa9defca18980e154098f2                                                                                                                                                                         0.0s
 => => exporting attestation manifest sha256:0ed217faec29b5afac16a4bcfbfd97859ffb6d9dc19e28fb937bc2239de78a2f                                                                                                                                                           0.0s
 => => exporting manifest list sha256:453f9ca01bd2ac102d084aa9ef1ee1d1c3937fc35714b3c730b648953ba20d99                                                                                                                                                                  0.0s
 => => pushing layers                                                                                                                                                                                                                                                   1.8s
 => => pushing manifest for docker.io/adamsidor/zadanie1:nieobowiazkowe@sha256:453f9ca01bd2ac102d084aa9ef1ee1d1c3937fc35714b3c730b648953ba20d99                                                                                                                         4.3s
 => exporting cache to registry                                                                                                                                                                                                                                       947.9s
 => => preparing build cache for export                                                                                                                                                                                                                                 4.0s
 => => sending cache export                                                                                                                                                                                                                                           943.9s
 => => writing layer sha256:2b1d8027cecb0f2dba0632bbc75e65c19bc18709686d22fb9902f9852d28f14e                                                                                                                                                                            0.2s
 => => writing layer sha256:17c96d64519c8439835eb056e8d714425d174e8a2370d8b447da3c66973131d4                                                                                                                                                                            0.1s
 => => writing layer sha256:05e89efc4b5ec39fa30d639b12ad2c7fd994a11ffa608e77427a882d73768d76                                                                                                                                                                          836.4s
 => => writing layer sha256:2a1381d785be07d85833f2dd3f4d0c86f4778fcf6ff112b487f2125c331a3752                                                                                                                                                                            0.2s
 => => writing layer sha256:2fb3bf578efb2a24965c6d12872aa0b4c8e360712fcbacf024d4a053867669a3                                                                                                                                                                          133.9s
 => => writing layer sha256:4aea57d93368be65a5f6fae264039d084106c503a0d1bcce101057e9d84358a2                                                                                                                                                                          465.6s
 => => writing layer sha256:4b83012ae5a32376a38507b6c37f34be83b8cde6d8d12d8ad5b895c35b0d0014                                                                                                                                                                            1.2s
 => => writing layer sha256:4c3aca1499f06ff02939c8f744463049feecb3965f1319b86b5e91790a9631e6                                                                                                                                                                            5.3s
 => => writing layer sha256:4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1                                                                                                                                                                            0.0s
 => => writing layer sha256:4fc331b5d0f09d74ddff683ea0835758b8819f54d5f07877d1201a060862d472                                                                                                                                                                          139.2s
 => => writing layer sha256:5358d0d3dca3b638b9e10e9af750136a843af77c17f35ca3a033d88f4735ff97                                                                                                                                                                            0.2s
 => => writing layer sha256:539f236c7c641fbc0316a5a0e5b602e9c790c1c8811514b3021961a891c571ca                                                                                                                                                                            1.3s
 => => writing layer sha256:5475db821899e81df61acc30229a8d58907cd1e3b673720ce55dfbf23fdbae21                                                                                                                                                                          435.9s
 => => writing layer sha256:64226052dc000123096737e625b3b5f1b37f4dafd8425229838f74933110287a                                                                                                                                                                          144.8s
 => => writing layer sha256:6a0ac1617861a677b045b7ff88545213ec31c0ff08763195a70a4a5adda577bb                                                                                                                                                                           70.6s
 => => writing layer sha256:7681648c89730867ec9970c4fb2635554c004c4bbec7ad384648fd05c912c2aa                                                                                                                                                                            1.3s
 => => writing layer sha256:89244935705494e38d5fa051ae5221080d8446b496516d3888d3cb64c7a7d51f                                                                                                                                                                            6.0s
 => => writing layer sha256:973f1b4898e9c2c22c54ff9d1bc2e2d45d59d4c2e53c60281dde84ef43df4c79                                                                                                                                                                            1.1s
 => => writing layer sha256:b105095b1a57fc2c5c46dbddf6137fcfaa0d6448bab218f93d6a19aa410a0faa                                                                                                                                                                            1.5s
 => => writing layer sha256:b55da06e3b41084804b2e5dbba71d35d26478b19ba8055d07a393cae22e9935f                                                                                                                                                                          567.9s
 => => writing layer sha256:c7b6b08c67539e0e5025fc9de65c1bae7a1e382a7545ea9439b5247db25fdf96                                                                                                                                                                            0.3s
 => => writing layer sha256:d17f077ada118cc762df373ff803592abf2dfa3ddafaa7381e364dd27a88fca7                                                                                                                                                                           73.5s
 => => writing layer sha256:e813a1647e7ea7745c1c15436c1f86a82eadb5436933163091db860d16ed923d                                                                                                                                                                            1.6s
 => => writing layer sha256:f3fcd6f34d996b5d5a262efb56d207261989034b9d2e6ce5927d96639dd92f60                                                                                                                                                                            1.2s
 => => writing config sha256:1541d159290593efc21e669f24f0ea73e4a5b030d56943d74aec75a2d5e5fe23                                                                                                                                                                           1.5s
 => => writing cache image manifest sha256:25fd93fe583d9427067cd023189b4301f725b01a4469107fe88cc0b718b6d566                                                                                                                                                             2.2s
 => [auth] adamsidor/zadanie1:pull,push token for registry-1.docker.io                                                                                                                                                                                                  0.0s
 => [auth] adamsidor/zadanie1:pull,push token for registry-1.docker.io                                                                                                                                                                                                  0.0s
 => [auth] adamsidor/zadanie1:pull,push token for registry-1.docker.io                                                                                                                                                                                                  0.0s
------
 > importing cache manifest from adamsidor/zadanie1:cache:
------

View build details: docker-desktop://dashboard/build/docker-container-builder/docker-container-builder0/l2u9osw3vrbpkn9c07djjb922

```

---

## 4. Weryfikacja obrazów w rejestrze

Po zakończeniu budowania sprawdzono manifesty obrazu głównego oraz obrazu cache.

**Sprawdzenie platform:**
```bash
docker buildx imagetools inspect adamsidor/zadanie1:nieobowiazkowe
```
**Wynik polecenia:**
```bash
Name:      docker.io/adamsidor/zadanie1:nieobowiazkowe
MediaType: application/vnd.oci.image.index.v1+json
Digest:    sha256:453f9ca01bd2ac102d084aa9ef1ee1d1c3937fc35714b3c730b648953ba20d99
           
Manifests: 
  Name:        docker.io/adamsidor/zadanie1:nieobowiazkowe@sha256:e244713cccf6d6614e3185906ed6dfce042de65db5c3176cf97add6d92445df7
  MediaType:   application/vnd.oci.image.manifest.v1+json
  Platform:    linux/amd64
               
  Name:        docker.io/adamsidor/zadanie1:nieobowiazkowe@sha256:92404d8c4cd88fd478a845b3d1e0fe789e5819234623bd6fffc57b58b21fd851
  MediaType:   application/vnd.oci.image.manifest.v1+json
  Platform:    linux/arm64
               
  Name:        docker.io/adamsidor/zadanie1:nieobowiazkowe@sha256:fbc643dbf18b69e817c9169e27d3f002096b327a80ac3504b0cdfa321d8f7fdd
  MediaType:   application/vnd.oci.image.manifest.v1+json
  Platform:    unknown/unknown
  Annotations: 
    vnd.docker.reference.digest: sha256:e244713cccf6d6614e3185906ed6dfce042de65db5c3176cf97add6d92445df7
    vnd.docker.reference.type:   attestation-manifest
               
  Name:        docker.io/adamsidor/zadanie1:nieobowiazkowe@sha256:0ed217faec29b5afac16a4bcfbfd97859ffb6d9dc19e28fb937bc2239de78a2f
  MediaType:   application/vnd.oci.image.manifest.v1+json
  Platform:    unknown/unknown
  Annotations: 
    vnd.docker.reference.digest: sha256:92404d8c4cd88fd478a845b3d1e0fe789e5819234623bd6fffc57b58b21fd851
    vnd.docker.reference.type:   attestation-manifest
```

Wynik potwierdził obecność manifestów dla `linux/amd64` oraz `linux/arm64`.

**Sprawdzenie cache:**
```bash
docker buildx imagetools inspect adamsidor/zadanie1:cache
```
**Wynik polecenia:**
```bash
Name:      docker.io/adamsidor/zadanie1:cache
MediaType: application/vnd.oci.image.manifest.v1+json
Digest:    sha256:25fd93fe583d9427067cd023189b4301f725b01a4469107fe88cc0b718b6d566
```

---

## 5. Optymalizacja procesu (Użycie Cache Registry)

Wykonano ponowne budowanie obrazu, aby zweryfikować poprawność działania mechanizmu cache.

**Szczegółowy log z widocznym użyciem CACHED:**
```text
[+] Building 8.7s (33/33) FINISHED                                                                                                                                                                                                 docker-container:docker-container-builder
 => [internal] load build definition from dockerfile                                                                                                                                                                                                                    0.0s
 => => transferring dockerfile: 857B                                                                                                                                                                                                                                    0.0s
 => resolve image config for docker-image://docker.io/docker/dockerfile:1                                                                                                                                                                                               1.1s
 => [auth] docker/dockerfile:pull token for registry-1.docker.io                                                                                                                                                                                                        0.0s
 => CACHED docker-image://docker.io/docker/dockerfile:1@sha256:2780b5c3bab67f1f76c781860de469442999ed1a0d7992a5efdf2cffc0e3d769                                                                                                                                         0.0s
 => => resolve docker.io/docker/dockerfile:1@sha256:2780b5c3bab67f1f76c781860de469442999ed1a0d7992a5efdf2cffc0e3d769                                                                                                                                                    0.0s
 => [linux/arm64 internal] load metadata for docker.io/library/golang:1.26-alpine                                                                                                                                                                                       0.5s
 => [linux/amd64 internal] load metadata for docker.io/library/golang:1.26-alpine                                                                                                                                                                                       0.5s
 => [auth] library/golang:pull token for registry-1.docker.io                                                                                                                                                                                                           0.0s
 => [internal] load .dockerignore                                                                                                                                                                                                                                       0.0s
 => => transferring context: 2B                                                                                                                                                                                                                                         0.0s
 => importing cache manifest from adamsidor/zadanie1:cache                                                                                                                                                                                                              1.5s
 => => inferred cache manifest type: application/vnd.oci.image.manifest.v1+json                                                                                                                                                                                         0.0s
 => [linux/amd64 builder 1/8] FROM docker.io/library/golang:1.26-alpine@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1                                                                                                                         0.0s
 => => resolve docker.io/library/golang:1.26-alpine@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1                                                                                                                                             0.0s
 => [linux/arm64 builder 1/8] FROM docker.io/library/golang:1.26-alpine@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1                                                                                                                         0.0s
 => => resolve docker.io/library/golang:1.26-alpine@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1                                                                                                                                             0.0s
 => [auth] adamsidor/zadanie1:pull token for registry-1.docker.io                                                                                                                                                                                                       0.0s
 => CACHED [linux/arm64 builder 2/8] RUN apk add --no-cache git openssh-client                                                                                                                                                                                          0.0s
 => CACHED [linux/arm64 builder 3/8] RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts                                                                                                                                                        0.0s
 => CACHED [linux/arm64 builder 4/8] WORKDIR /app                                                                                                                                                                                                                       0.0s
 => CACHED [linux/arm64 builder 5/8] RUN --mount=type=ssh git clone git@github.com:Adam-Sidor/Docker_lab.git .                                                                                                                                                          0.0s
 => CACHED [linux/arm64 builder 6/8] WORKDIR /app/Zadanie1                                                                                                                                                                                                              0.0s
 => CACHED [linux/arm64 builder 7/8] RUN go mod download || true                                                                                                                                                                                                        0.0s
 => CACHED [linux/arm64 builder 8/8] RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o weather-app .                                                                                                                                                            0.0s
 => CACHED [linux/arm64 stage-1 1/2] COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/                                                                                                                                                             0.0s
 => CACHED [linux/arm64 stage-1 2/2] COPY --from=builder /app/Zadanie1/weather-app /weather-app                                                                                                                                                                         0.0s
 => CACHED [linux/amd64 builder 2/8] RUN apk add --no-cache git openssh-client                                                                                                                                                                                          0.0s
 => CACHED [linux/amd64 builder 3/8] RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts                                                                                                                                                        0.0s
 => CACHED [linux/amd64 builder 4/8] WORKDIR /app                                                                                                                                                                                                                       0.0s
 => CACHED [linux/amd64 builder 5/8] RUN --mount=type=ssh git clone git@github.com:Adam-Sidor/Docker_lab.git .                                                                                                                                                          0.0s
 => CACHED [linux/amd64 builder 6/8] WORKDIR /app/Zadanie1                                                                                                                                                                                                              0.0s
 => CACHED [linux/amd64 builder 7/8] RUN go mod download || true                                                                                                                                                                                                        0.0s
 => CACHED [linux/amd64 builder 8/8] RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o weather-app .                                                                                                                                                            0.0s
 => CACHED [linux/amd64 stage-1 1/2] COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/                                                                                                                                                             0.0s
 => CACHED [linux/amd64 stage-1 2/2] COPY --from=builder /app/Zadanie1/weather-app /weather-app                                                                                                                                                                         0.0s
 => exporting to image                                                                                                                                                                                                                                                  5.4s
 => => exporting layers                                                                                                                                                                                                                                                 0.0s
 => => exporting manifest sha256:e244713cccf6d6614e3185906ed6dfce042de65db5c3176cf97add6d92445df7                                                                                                                                                                       0.0s
 => => exporting config sha256:91f063c5080e71af9b6da7b2e24f63a9cdc71e5396bf65ecda5e0d28a65b43c7                                                                                                                                                                         0.0s
 => => exporting attestation manifest sha256:8c19957f74955881eb877e9aa00ad4d39543056883b605b46799c8a7563cab29                                                                                                                                                           0.0s
 => => exporting manifest sha256:92404d8c4cd88fd478a845b3d1e0fe789e5819234623bd6fffc57b58b21fd851                                                                                                                                                                       0.0s
 => => exporting config sha256:e7e921527d86221c9823abcea30b2feb4edd29fdf5aa9defca18980e154098f2                                                                                                                                                                         0.0s
 => => exporting attestation manifest sha256:71f55d25eb2779420e892a6b0c3a69152e54906c3362b03e74f2e85d62c17fd5                                                                                                                                                           0.0s
 => => exporting manifest list sha256:16edb64e4f5045d554fac4c2711cb94a550aa69f9060861dc0f23bded78e1d0d                                                                                                                                                                  0.0s
 => => pushing layers                                                                                                                                                                                                                                                   1.7s
 => => pushing manifest for docker.io/adamsidor/zadanie1:nieobowiazkowe@sha256:16edb64e4f5045d554fac4c2711cb94a550aa69f9060861dc0f23bded78e1d0d                                                                                                                         3.6s
 => exporting cache to registry                                                                                                                                                                                                                                         4.7s
 => => preparing build cache for export                                                                                                                                                                                                                                 0.1s
 => => sending cache export                                                                                                                                                                                                                                             4.6s
 => => writing layer sha256:2b1d8027cecb0f2dba0632bbc75e65c19bc18709686d22fb9902f9852d28f14e                                                                                                                                                                            0.5s
 => => writing layer sha256:17c96d64519c8439835eb056e8d714425d174e8a2370d8b447da3c66973131d4                                                                                                                                                                            0.5s
 => => writing layer sha256:2a1381d785be07d85833f2dd3f4d0c86f4778fcf6ff112b487f2125c331a3752                                                                                                                                                                            0.6s
 => => writing layer sha256:05e89efc4b5ec39fa30d639b12ad2c7fd994a11ffa608e77427a882d73768d76                                                                                                                                                                            0.6s
 => => writing layer sha256:2fb3bf578efb2a24965c6d12872aa0b4c8e360712fcbacf024d4a053867669a3                                                                                                                                                                            0.1s
 => => writing layer sha256:4aea57d93368be65a5f6fae264039d084106c503a0d1bcce101057e9d84358a2                                                                                                                                                                            0.1s
 => => writing layer sha256:4b83012ae5a32376a38507b6c37f34be83b8cde6d8d12d8ad5b895c35b0d0014                                                                                                                                                                            0.1s
 => => writing layer sha256:4c3aca1499f06ff02939c8f744463049feecb3965f1319b86b5e91790a9631e6                                                                                                                                                                            0.1s
 => => writing layer sha256:4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1                                                                                                                                                                            0.1s
 => => writing layer sha256:4fc331b5d0f09d74ddff683ea0835758b8819f54d5f07877d1201a060862d472                                                                                                                                                                            0.2s
 => => writing layer sha256:5358d0d3dca3b638b9e10e9af750136a843af77c17f35ca3a033d88f4735ff97                                                                                                                                                                            0.1s
 => => writing layer sha256:539f236c7c641fbc0316a5a0e5b602e9c790c1c8811514b3021961a891c571ca                                                                                                                                                                            0.1s
 => => writing layer sha256:5475db821899e81df61acc30229a8d58907cd1e3b673720ce55dfbf23fdbae21                                                                                                                                                                            0.1s
 => => writing layer sha256:64226052dc000123096737e625b3b5f1b37f4dafd8425229838f74933110287a                                                                                                                                                                            0.2s
 => => writing layer sha256:6a0ac1617861a677b045b7ff88545213ec31c0ff08763195a70a4a5adda577bb                                                                                                                                                                            0.1s
 => => writing layer sha256:7681648c89730867ec9970c4fb2635554c004c4bbec7ad384648fd05c912c2aa                                                                                                                                                                            0.1s
 => => writing layer sha256:89244935705494e38d5fa051ae5221080d8446b496516d3888d3cb64c7a7d51f                                                                                                                                                                            0.1s
 => => writing layer sha256:973f1b4898e9c2c22c54ff9d1bc2e2d45d59d4c2e53c60281dde84ef43df4c79                                                                                                                                                                            0.2s
 => => writing layer sha256:b105095b1a57fc2c5c46dbddf6137fcfaa0d6448bab218f93d6a19aa410a0faa                                                                                                                                                                            0.2s
 => => writing layer sha256:b55da06e3b41084804b2e5dbba71d35d26478b19ba8055d07a393cae22e9935f                                                                                                                                                                            0.2s
 => => writing layer sha256:c7b6b08c67539e0e5025fc9de65c1bae7a1e382a7545ea9439b5247db25fdf96                                                                                                                                                                            0.2s
 => => writing layer sha256:d17f077ada118cc762df373ff803592abf2dfa3ddafaa7381e364dd27a88fca7                                                                                                                                                                            0.2s
 => => writing layer sha256:e813a1647e7ea7745c1c15436c1f86a82eadb5436933163091db860d16ed923d                                                                                                                                                                            0.2s
 => => writing layer sha256:f3fcd6f34d996b5d5a262efb56d207261989034b9d2e6ce5927d96639dd92f60                                                                                                                                                                            0.2s
 => => writing config sha256:9ad24ef671879e8be1e13a14fefe7140ca062987c0f0b1442add53d69f0c7e76                                                                                                                                                                           1.1s
 => => writing cache image manifest sha256:e91ecf22bb294cc57b03826a5a976a7d1b700d296f0e71570ba7fa50ad3c603f                                                                                                                                                             2.1s
 => [auth] adamsidor/zadanie1:pull,push token for registry-1.docker.io                                                                                                                                                                                                  0.0s

View build details: docker-desktop://dashboard/build/docker-container-builder/docker-container-builder0/jb80o05y57cnhetluisgbza80                                                                                                                                                                                                                                   4.7s
```
*(Czas budowania spadł z 982.7s do 8.7s).*
Dzięki wykorzystaniu Cache Registry czas budowania skrócił się z ponad 16 minut do niecałych 9 sekund.

---

## 6. Końcowa analiza bezpieczeństwa (Docker Scout)

Na koniec przeprowadzono skanowanie nowego obrazu. Skanowanie przeprowadzono dla obu architektur, potwierdzając brak podatności.

**Polecenia:**
- Dla ARM64:
  ```bash
  docker scout cves adamsidor/zadanie1:nieobowiazkowe > nieobowiazkowe/scout_report_nieobowiazkowe.log

      ...Storing image for indexing
      ✓ Image stored for indexing
      ...Indexing
      ✓ Indexed 2 packages
      ✓ No vulnerable package detected
  ```
- Dla AMD64:
  ```bash
  docker scout cves adamsidor/zadanie1:nieobowiazkowe --platform linux/amd64 > nieobowiazkowe/scout_report_amd64.log

      ...Pulling
      ✓ Pulled
      ...Storing image for indexing
      ✓ Image stored for indexing
      ...Indexing
      ✓ Indexed 2 packages
      ✓ Provenance obtained from attestation
      ✓ No vulnerable package detected
  ```

**Raporty:**
- `scout_report_nieobowiazkowe.log` (ARM64)
  ```text
  ## Overview

                    │           Analyzed Image            
  ──────────────────┼─────────────────────────────────────
  Target            │  adamsidor/zadanie1:nieobowiazkowe  
    digest          │  e7e921527d86                       
    platform        │ linux/arm64                         
    vulnerabilities │    0C     0H     0M     0L          
    size            │ 3.8 MB                              
    packages        │ 2                                   


  ## Packages and Vulnerabilities

    No vulnerable packages detected                                            
  ```
- `scout_report_amd64.log` (AMD64)
  ```text
  ## Overview

                    │                Analyzed Image                
  ──────────────────┼──────────────────────────────────────────────
  Target            │  adamsidor/zadanie1:nieobowiazkowe           
    digest          │  e244713cccf6                                
    platform        │ linux/amd64                                  
    provenance      │ https://github.com/Adam-Sidor/Docker_lab.git 
                    │  https://github.com/Adam-Sidor/Docker_lab/blob/ea10e19c406f60438de92e3523691fcf48f145a9    
    vulnerabilities │    0C     0H     0M     0L                   
    size            │ 3.8 MB                                       
    packages        │ 2                                            


  ## Packages and Vulnerabilities

    No vulnerable packages detected                                      
  ```

---
**Docker Hub:** [https://hub.docker.com/r/adamsidor/zadanie1/tags](https://hub.docker.com/r/adamsidor/zadanie1/tags)
