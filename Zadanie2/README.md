# Zadanie 2 - GitHub Actions

## 1. Zasoby projektu

- **Kod źródłowy aplikacji**: [Zadanie 1](https://github.com/Adam-Sidor/Docker_lab/tree/master/Zadanie1)
- **Przebieg Pipeline (Actions)**: [GitHub Actions Logs](https://github.com/Adam-Sidor/Docker_lab/actions)
- **Opublikowane Obrazy (Packages)**: [GitHub Container Registry](https://github.com/Adam-Sidor/Docker_lab/packages)
- **Dane Cache**: [DockerHub Registry Cache](https://hub.docker.com/repository/docker/adamsidor/zadanie2/tags)

## 2. Pełny kod Pipeline (GitHub Actions)

Poniżej znajduje się konfiguracja pliku `.github/workflows/docker-pipeline.yml` z opisem poszczególnych etapów:

```yaml
name: Zadanie 2

on:
  workflow_dispatch: # Umożliwia ręczne uruchomienie pipeline'u
  push:
    tags: [ 'zadanie2-v*' ] # Uruchamia proces tylko po wypchnięciu tagu wersji, np. zadanie2-v1.0

jobs:
  ci_step:
    name: Budowa, skanowanie i wysyłka obrazu Docker
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write # Wymagane uprawnienie do wysyłki obrazów do GHCR

    steps:
      - name: Pobranie kodu repozytorium
        uses: actions/checkout@v4

      - name: Definicja metadanych Docker
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: |
            ghcr.io/${{ github.repository_owner }}/docker_lab
          flavor: |
            latest=true # Automatyczne dodawanie tagu :latest
          tags: |
            type=sha,priority=100,prefix=sha-,format=short # Tagowanie skrótem SHA
            type=match,pattern=zadanie2-v(.*),group=1 # Wyciąganie wersji z tagu zadanie2-v*

      - name: Konfiguracja QEMU
        uses: docker/setup-qemu-action@v3
        with:
          platforms: arm64,amd64 # Przygotowanie do budowy multiplatformowej

      - name: Konfiguracja Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Logowanie do GitHub Container Registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Logowanie do DockerHub (dla cache)
        uses: docker/login-action@v3
        with:
          username: ${{ vars.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      - name: Budowa obrazu lokalnego do skanowania
        uses: docker/build-push-action@v5
        with:
          context: ./Zadanie1 # Kontekst budowania to folder Zadania 1
          file: ./Zadanie1/dockerfile
          load: true # Załadowanie obrazu do lokalnego dockera w celu skanowania
          tags: weather-app:scan
          cache-from: type=registry,ref=${{ vars.DOCKERHUB_USERNAME }}/zadanie2:cache
          cache-to: type=registry,ref=${{ vars.DOCKERHUB_USERNAME }}/zadanie2:cache,mode=max

      - name: Skanowanie obrazu skanerem Trivy
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: 'weather-app:scan'
          format: 'table'
          exit-code: '1' # Przerwanie procesu, jeśli zostaną znalezione błędy
          ignore-unfixed: true
          vuln-type: 'os,library'
          severity: 'CRITICAL,HIGH' # Skanowanie pod kątem zagrożeń Krytycznych i Wysokich

      - name: Budowa i wypchnięcie obrazu wieloarchitekturowego
        if: success() # Krok wykonuje się tylko, jeśli test CVE (Trivy) zakończył się sukcesem
        uses: docker/build-push-action@v5
        with:
          context: ./Zadanie1
          file: ./Zadanie1/dockerfile
          platforms: linux/amd64,linux/arm64
          push: true # Ostateczne wypchnięcie do GHCR
          tags: ${{ steps.meta.outputs.tags }}
          cache-from: type=registry,ref=${{ vars.DOCKERHUB_USERNAME }}/zadanie2:cache
          cache-to: type=registry,ref=${{ vars.DOCKERHUB_USERNAME }}/zadanie2:cache,mode=max
```

## 3. Strategia tagowania - Uzasadnienie Techniczne

### Wyzwalacz Pipeline'u (Push Tags)
W konfiguracji workflow zdecydowano się na wyzwalanie akcji na podstawie wypchnięcia określonych tagów (`push: tags: [ 'zadanie2-v*' ]`).

**Uzasadnienie:**
Ponieważ repozytorium jest współdzielone dla wszystkich laboratoriów z tego przedmiotu, konieczne było odizolowanie uruchamiania pipeline'ów dla poszczególnych zadań. Konfiguracja tagów `zadanie2-v*` gwarantuje, że proces budowania dla Zadania 2 uruchomi się tylko i wyłącznie wtedy, gdy zmiany dotyczą tego konkretnego zadania. Zapobiega to niepotrzebnemu uruchamianiu akcji podczas commitowania i pushowania zmian do innych laboratoriów.

### Tagowanie obrazów aplikacyjnych
W rozwiązaniu przyjęto nowoczesną strategię wielopoziomowego tagowania, co zapewnia elastyczność i bezpieczeństwo:
- **Tag SHA (`sha-<skrót>`)**: Zapewnia **identyfikowalność (traceability)**. Pozwala na precyzyjne powiązanie obrazu z konkretnym commitem w Git. Jest to kluczowe dla audytów bezpieczeństwa i debugowania, eliminując niepewność co do zawartości binarnej obrazu.
- **Tagowanie wersjami (`SemVer`)**: Gwarantuje **niezmienność (immutability)** wydanych już wersji. Użytkownik ma pewność, że raz pobrana wersja nie ulegnie zmianie, co mogłoby nastąpić przy użyciu samego tagu `latest`.
- **Tag `latest`**: Służy jako dynamiczny alias ułatwiający pobieranie najnowszej stabilnej wersji w środowiskach deweloperskich.

### Tagowanie danych cache
Zgodnie z wymaganiami zadania, dane cache są przechowywane w dedykowanym repozytorium na DockerHub z użyciem tagu `:cache`.

**Uzasadnienie wyboru:**
Użycie dedykowanego tagu zapewnia izolację danych tymczasowych (warstw budowania) od obrazów produkcyjnych. Ułatwia to zarządzanie rejestrem, zapewnia przejrzystość struktury artefaktów i pozwala na optymalne wykorzystanie trybu `mode=max` bez "zaśmiecania" głównego repozytorium obrazów.

## 4. Wybór skanera CVE: Trivy vs Docker Scout

Do testów bezpieczeństwa wybrany został skaner **Trivy**.

**Uzasadnienie:**
Wybór Trivy podyktowany został przede wszystkim łatwością integracji typu "Plug & Play" w środowisku GitHub Actions. W przeciwieństwie do Docker Scout, Trivy nie wymaga skomplikowanej autoryzacji w chmurze zewnętrznej do przeprowadzenia lokalnego skanowania, co czyni proces bardziej hermetycznym. Trivy generuje niezwykle przejrzyste raporty tabelaryczne bezpośrednio w logach, co ułatwia szybką analizę bezpieczeństwa obrazu przed jego publikacją.

## 5. Pełne Logi z przebiegu prac (CLI)

### Etap 1: Konfiguracja zmiennych i sekretów
W kroku tym skonfigurowano niezbędne dane uwierzytelniające dla GitHub Actions, aby proces budowania mógł zalogować się do rejestrów zewnętrznych (DockerHub) bez ujawniania hasła w kodzie.

```text
adam@Adams-MacBook-Pro Docker % gh variable set DOCKERHUB_USERNAME                    

? Paste your variable adamsidor

✓ Created variable DOCKERHUB_USERNAME for Adam-Sidor/Docker_lab


adam@Adams-MacBook-Pro Docker % gh secret set DOCKERHUB_TOKEN 
? Paste your secret: ************************************

✓ Set Actions secret DOCKERHUB_TOKEN for Adam-Sidor/Docker_lab
```

### Etap 2: Wersjonowanie i wypchnięcie kodu
Zatwierdzono plik pipeline'u oraz stworzono specjalny tag `zadanie2-v1.0`, który zgodnie z konfiguracją `on: push: tags` uruchomił proces budowania.

```text
adam@Adams-MacBook-Pro Docker % git add .github 


adam@Adams-MacBook-Pro Docker % git commit -m "Zadanie 2: gotowa konfiguracja z tagowaniem zadanie2-v*"
[master d765867] Zadanie 2: gotowa konfiguracja z tagowaniem zadanie2-v*
 1 file changed, 83 insertions(+)
 create mode 100644 .github/workflows/docker-pipeline.yml

adam@Adams-MacBook-Pro Docker % git push      
Enumerating objects: 6, done.
Counting objects: 100% (6/6), done.
Delta compression using up to 10 threads
Compressing objects: 100% (3/3), done.
Writing objects: 100% (5/5), 1.28 KiB | 1.28 MiB/s, done.
Total 5 (delta 1), reused 0 (delta 0), pack-reused 0 (from 0)
remote: Resolving deltas: 100% (1/1), completed with 1 local object.
To https://github.com/Adam-Sidor/Docker_lab.git
   9da540f..d765867  master -> master

adam@Adams-MacBook-Pro Docker % git tag zadanie2-v1.0

adam@Adams-MacBook-Pro Docker % git push origin zadanie2-v1.0
Total 0 (delta 0), reused 0 (delta 0), pack-reused 0 (from 0)
To https://github.com/Adam-Sidor/Docker_lab.git
 * [new tag]         zadanie2-v1.0 -> zadanie2-v1.0
```

### Etap 3: Monitorowanie i weryfikacja sukcesu (Action Dashboard)
Użyto narzędzia `gh run watch`, aby śledzić postęp prac na serwerach GitHub w czasie rzeczywistym. Widoczny jest sukces każdego z etapów, w tym skanowania Trivy oraz budowy wieloplatformowej.

```text
 adam@Adams-MacBook-Pro Docker % gh run watch

? Select a workflow run * Zadanie 2: gotowa konfiguracja z tagowaniem zadanie2-v*, Zadanie 2 [zadanie2-v1.0] 1s ago
✓ zadanie2-v1.0 Zadanie 2 · 25604593368
Triggered via push about 4 minutes ago

JOBS
✓ Budowa, skanowanie i wysyłka obrazu Docker in 4m16s (ID 75164208635)
  ✓ Set up job
  ✓ Pobranie kodu repozytorium
  ✓ Definicja metadanych Docker
  ✓ Konfiguracja QEMU
  ✓ Konfiguracja Docker Buildx
  ✓ Logowanie do GitHub Container Registry
  ✓ Logowanie do DockerHub (dla cache)
  ✓ Budowa obrazu lokalnego do skanowania
  ✓ Skanowanie obrazu skanerem Trivy
  ✓ Budowa i wypchnięcie obrazu wieloarchitekturowego
  ✓ Post Budowa i wypchnięcie obrazu wieloarchitekturowego
  ✓ Post Skanowanie obrazu skanerem Trivy
  ✓ Post Budowa obrazu lokalnego do skanowania
  ✓ Post Logowanie do DockerHub (dla cache)
  ✓ Post Logowanie do GitHub Container Registry
  ✓ Post Konfiguracja Docker Buildx
  ✓ Post Konfiguracja QEMU
  ✓ Post Pobranie kodu repozytorium
  ✓ Complete job

✓ Run Zadanie 2 (25604593368) completed with 'success'
```
