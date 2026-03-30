# Warianty powiązania BPMN + DMN + Form (Camunda 8)

Przegląd wariantów powiązania BPMN, DMN i formularzy oraz kryteria wyboru wariantu.

## Wariant A — Pełny stek (formularz → decyzja DMN → workery)

**Idea:** użytkownik wprowadza dane w **Tasklist** (formularz), silnik liczy trasę wg **DMN**, potem **różne service task** (`c8jw-*`).  
**Plusy:** widać formularz UI, zarządzalną politykę (tabela) i kod workerów.  
**Minusy:** potrzebne **Operate/Tasklist** i wdrożenie **trzech** zasobów w sensownej kolejności.

**Implementacja w repozytorium:** zob. zestaw plików poniżej (prefiks `c8jw_` / `credit`).

## Wariant B — Tylko BPMN + DMN (bez formularza)

**Idea:** zmienne podajesz przy starcie instancji (Operate / API), potem **Business Rule Task** i workery.  
**Plusy:** prostsze demo bez Tasklist.  
**Minusy:** brak „ładnego” wprowadzania dla osób nietechnicznych.

## Wariant C — Formularz + workery, bez DMN

**Idea:** warunek na **Exclusive Gateway** bezpośrednio w BPMN (`amount < 5000`).  
**Plusy:** mniej artefaktów, blisko `bpmn/diagram_1.bpmn`.  
**Minusy:** brak osobnej **zarządzalnej** tabeli reguł (biznes mniej wyraźnie widzi zmianę bez przebudowy diagramu).

## Wariant D — DMN jako „słownik”, wynik tylko w zmiennych

**Idea:** po DMN proces idzie jednym przepływem (np. tylko logowanie `c8jw-python`).  
**Plusy:** minimum rozgałęzień, akcent na wywołaniu DMN.  
**Minusy:** słabiej widać rozgałęzienie według decyzji.

---

## Wdrożenie zestawu A w Modeler / Connectors

1. **`forms/c8jw-demo-application.form`** — Form (id = `c8jw_demo_application`).  
2. **`dmn/credit-route.dmn`** — decyzja `credit_route_decision`.  
3. **`bpmn/examples/credit-orchestration.bpmn`** — proces `c8jw_credit_orchestration`.

Zaleca się jednym **deploy** wysłać wszystkie trzy pliki (lub najpierw formularz + DMN, potem BPMN). W **`credit-orchestration.bpmn`** do wywołania DMN użyto **`bindingType="latest"`** — można opublikować DMN przed BPMN. Wariant **`deployment`** jest wygodny, gdy DMN i proces zawsze idą w jednej paczce.

Po starcie procesu: wypełnij formularz (**imię**, **kwota**). DMN ustawi `creditRoute` = `"auto"` lub `"manual"`; dalej zadziała `c8jw-golang` lub `c8jw-python` (potrzebne oba workery albo usuń gałąź w BPMN).

---

## Formularz z wyborem workera → dynamiczny job type → formularz wyniku

**Idea:** w Tasklist użytkownik wybiera **job type** z listy; service task tworzy się z **typem ze zmiennej** (`<zeebe:taskDefinition type="=jobType" />`). Po **complete** workera zmienne procesu trafiają do **drugiego** user task z formularzem, pola **read-only** z tym samym `key` co wyjście workera (w granicach możliwości Tasklist / wersji formularzy).

**Pliki**

1. `forms/c8jw-worker-picker.form` — `id`: **c8jw_worker_picker** (`jobType`, opcjonalnie `name`, `amount`).  
2. `forms/c8jw-worker-result.form` — `id`: **c8jw_worker_result**.  
3. `bpmn/examples/worker-picker-orchestration.bpmn` — proces **c8jw_worker_form_demo**.

Wdrożenie **dwóch** formularzy i BPMN. Musi działać worker wybranego typu (dla demo wygodnie podnieść wszystkie cztery w GKE lub kilka job types w jednym lokalnym workerze).

---

## Złożone DRG DMN (drzewo decyzji z zależnościami)

**Cel:** pokazać **Decision Requirements Graph** — dwie decyzje tabelaryczne (`risk_level`, `loyalty_tier`) i **literał FEEL** łączący wyniki w jedną trasę (`approval_routing_tree`).

**Pliki**

1. `dmn/complex-decision-tree.dmn` — DRG: `risk_level` (kwota + scoring), `loyalty_tier` (staż + opóźnienia), korzeń **`approval_routing_tree`** → zmienna wynikowa **`serviceRoute`** (wartości m.in. `INSTANT_CREDIT`, `MANUAL_UNDERWRITING`, `COMPLIANCE_REVIEW`). W literale użyto **`risk_level.band`** i **`loyalty_tier.tier`** (nazwy kolumn wyjściowych tabel).
2. `bpmn/examples/dmn-complex-tree-demo.bpmn` — proces **`c8jw_dmn_complex_tree_demo`**: jeden **Business Rule Task** wywołuje **`approval_routing_tree`** (`bindingType="latest"`, `resultVariable="serviceRoute"`).

**Zmienne startowe (liczby):** `loanAmount`, `creditScore` (np. 0–1000), `monthsActive`, `priorDefaults` (np. 0). Wdróż **DMN i BPMN** (jak w pozostałych przykładach).
