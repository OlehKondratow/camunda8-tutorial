# Camunda 8: zadania serwisowe i bramy (notatki)

**Źródło:** [Reunico — Materiał o Service Tasks i Gateway w Camunda 8](https://reunico.com/blog/camunda-8-service-tasks-gateways/)  
**Data na stronie:** 19.12.2024  
**Autor:** Mstislav Martyniuk  
**Czas materiału:** ~15 min (szkoleniowe)

> Ćwiczenie w oryginale zakłada doświadczenie z **Spring Boot** (Spring Zeebe SDK). Ilustracje i wideo — tylko na stronie źródłowej.

---

## Zakres tematów

- Użycie **bramy** w procesie — w przykładzie **Exclusive Gateway** (XOR / „rozgałęzienie LUB/LUB”).
- Realizacja **Service Task** przez wzorzec **Job Worker**.

Część cyklu „Uczymy się Camunda 8 z Reunico”; zob. też playlistę na YouTube i sąsiednie artykuły.

---

## Bramy (Gateways)

Typy sensownie:

- decyzja wg **warunku**: Exclusive / Inclusive Gateway;
- decyzja wg **zdarzenia**: Event Gateway;
- **równoległe rozgałęzienie** i **scalenie**: Parallel Gateway.

### Exclusive Gateway

Przy wejściu tokena sprawdzane są **warunki na wychodzących sequence flow**. W Camunda 8 warunki zapisuje się w **FEEL**. Token idzie **pierwszą** ścieżką, dla której warunek jest prawdziwy.

**Przykład z artykułu:** zmienna `amount` (kwota kredytu):

- jeśli `amount < 5000` — gałąź do **zadania serwisowego** (autoobsługa);
- w przeciwnym razie — gałąź do **zadania użytkownika**.

**Default flow:** oznaczany „kluczem”; działa, gdy **wszystkie pozostałe** warunki są fałszywe. W przykładzie: ścieżka „≥5000” może być default, ścieżka „<5000” — z warunkiem `amount < 5000`.

---

## Zadania serwisowe i Job Worker

Camunda 8 orkiestruje nie tylko kroki użytkownika, ale i wywołania usług/API.

**Service Task** z implementacją **Job Worker**:

1. Wykonywanie procesu **zatrzymuje się** na zadaniu.
2. Powstaje **job** o danym **Job type**.
3. Klient **Zeebe** (gRPC) pobiera job i wykonuje logikę.
4. Po wykonaniu — **complete** lub **fail** (z retry).

Implementacja workera: własny klient Zeebe lub SDK (w artykule — **Spring Zeebe**); alternatywnie inne oficjalne klienty.

---

## Praktyka (krótko wg scenariusza z artykułu)

1. Dodać **Exclusive Gateway** przed zadaniem „obsłuż wniosek”, nazwa bramy jako pytanie (np. „Kwota mniejsza niż 5000?”).
2. Nadać wychodzącym flow nazwy i warunki FEEL; jeden flow ustawić jako **default**.
3. Dla gałęzi autoobsługi — **Service Task**, **Job type** (w kursie często `decision`; w tym repozytorium szkoleniowe workery używają **`c8jw-golang`** itd., prefiks **`c8jw-*`** — zob. `kubernetes/README.md`).
4. Wdrożyć proces do klastra Camunda.
5. Zaimplementować **Job Worker** (np. Spring Boot + Spring Zeebe).
6. Uruchamiać proces z różnymi `amount` (np. 3000, 10000, 1000) i sprawdzać w **Operate** / **Tasklist**, która gałąź została wybrana.

Oczekiwane zachowanie (wg tekstu artykułu):

- przy **3000** i **1000** — wykonuje się autoobsługa;
- przy **10000** — powstaje user task w Tasklist;
- po ręcznym ukończeniu user task instancje śledzi się w Operate.

---

## Przydatne linki ze strony źródłowej

- Cykl / sąsiednie materiały na [reunico.com](https://reunico.com/blog/).
- Dokumentacja Camunda o klientach Zeebe i Job Workers: [docs.camunda.io](https://docs.camunda.io).
