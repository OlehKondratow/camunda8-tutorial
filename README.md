# Camunda 8 — repozytorium szkoleniowe

Samodzielny zestaw do ćwiczeń: **Zeebe job workers** (Go i Python), przykładowe BPMN, formularz, Docker Compose, ściąga **GKE**.

**Repozytorium:** [github.com/OlehKondratow/camunda8-tutorial](https://github.com/OlehKondratow/camunda8-tutorial)

Może równolegle służyć jako **materiał referencyjny do portfolio** (Camunda 8, Zeebe, Kubernetes); treść poniżej ma wyłącznie charakter **techniczny**.

## Zakres techniczny

- **Camunda 8 / Zeebe** (w odróżnieniu od Camunda 7): orchestracja przez **job workers** (gRPC), nie przez osadzony silnik w aplikacji.
- **BPMN** w Modelerze; wykonanie i monitoring po stronie platformy (**Operate**, **Tasklist**).
- **Service Task → job type** (`c8jw-*`): aktywacja joba, **complete** z **variables**; spójne nazewnictwo obrazu, StatefulSet w Kubernetes i typu zadania w BPMN.
- Przykład **formularz + DMN + BPMN** (end-to-end): `zeebe-tutorial/bpmn-and-dmn-bundles.md`, proces **`c8jw_credit_orchestration`**.
- Workery w **Go**, **Python** i **Java** — jeden protokół Zeebe, różne środowiska uruchomieniowe.
- **Obrazy** w **Artifact Registry**, build w **Cloud Build**, wdrożenie workerów w **GKE** (`kubernetes/`, `ci/`).
- **Docker Compose** — lokalny Zeebe i workery bez klastra.

Typowy scenariusz weryfikacji: uruchomienie Zeebe i workera → wdrożenie procesu → start instancji w **Operate** → wykonanie joba → zmienne procesu i kolejny krok na diagramie.

Diagram z **dwoma job types** pod rząd: `zeebe-tutorial/bpmn/examples/portfolio-pipeline.bpmn` (process id **`c8jw-portfolio-pipeline`**). Przy tworzeniu instancji podaj zmienne, np. `{"name":"Alice","amount":1500}` — w przeciwnym razie krok **c8jw-golang** zwróci błąd (implementacja: `zeebe-tutorial/internal/tutorial/decision_task.go`).

## Struktura

| Ścieżka | Rola |
|---------|------|
| `zeebe-tutorial/` | Workery `c8jw-*`; formularze+BPMN: **`bpmn-and-dmn-bundles.md`**; kredyt: **credit-orchestration** + **credit-route.dmn** + **c8jw-demo-application.form**; **złożone DRG:** **complex-decision-tree.dmn** + **dmn-complex-tree-demo.bpmn** (`c8jw_dmn_complex_tree_demo`); wybór workera: **worker-picker-orchestration.bpmn**, **c8jw-worker-*.form**; także `example-task-process`, **portfolio-pipeline** |
| `images/operate/` | Zrzuty ekranu **Camunda Operate** |
| `images/modeler/` | Zrzuty **Camunda Modeler** |
| `docker-compose.yaml` | Lokalny **Zeebe 8.5** + oba workery |
| `docs/gke-camunda-cheatsheet.md` | GKE, Helm, port-forward, PVC, sprzątanie |
| `docs/service-tasks-and-gateways-reunico.md` | Service tasks / bramy (Reunico) |
| `examples/application.form` | Przykładowy formularz (Camunda) |
| ~~`credit-scoring-camunda/`~~ | **Przeniesiony** poza to repozytorium — osobny katalog roboczy (np. `/data/projects/credit-scoring-camunda` obok tego klona). Pełna dokumentacja w `README.md` tamże. |
| `examples/diagram_1.bpmn` | Dodatkowy diagram |
| `helm/` | Values Camunda: przykład **`camunda-platform-user-values.yaml`**, eksport z klastra — `helm/README.md`, `scripts/fetch-helm-values.sh` |
| `kubernetes/` | Kilka szkoleniowych **job workers** dla GKE (zob. `kubernetes/README.md`) |
| `ci/cloudbuild-workers.yaml` | Cloud Build: wszystkie obrazy workerów → Artifact Registry |
| `scripts/deploy-workers-to-gke.sh` | Wdrożenie workerów w namespace `c8-tutorial-workers` |

## Szybki start (lokalnie)

Wymagania: Docker lub Podman Compose; opcjonalnie Go 1.21+, Python 3.12.

```bash
git clone https://github.com/OlehKondratow/camunda8-tutorial.git
cd camunda8-tutorial
docker compose up --build
```

Brama: **localhost:26500**. Workery nasłuchują **`c8jw-golang`** i **`c8jw-python`**.

- **Go:** `cd zeebe-tutorial && export ZEEBE_ADDRESS=127.0.0.1:26500 && go run ./cmd/example-worker`  
  Zmienna **`JOB_TYPES`** (po przecinku), np. `JOB_TYPES=c8jw-golang`.
- **Python:** `zeebe-tutorial/python/example-zeebe-worker/README.md`

BPMN: **`zeebe-tutorial/bpmn/examples/example-task.bpmn`** (`example-task-process`) oraz łańcuch **Python → Golang** — **`portfolio-pipeline.bpmn`** (`c8jw-portfolio-pipeline`). Proces z kursu — **`examples/diagram_1.bpmn`** / Modeler (gałąź automatycznej obsługi, Type **`c8jw-golang`**) ewentualnie wdrożyć osobno.

W **Modeler** w Service Task pole **Type** musi być zgodne ze schemą **`c8jw-*`** w `kubernetes/README.md`. **Process ID** nie jest powiązany z wdrożeniem workera.

## Chmura

- Ściąga: **`docs/gke-camunda-cheatsheet.md`**
- Build obrazów: `gcloud builds submit --config=ci/cloudbuild-workers.yaml --project=my-camunda8-project .`
- Wdrożenie do klastra: `./scripts/deploy-workers-to-gke.sh` (potrzebne repozytorium **`c8-tutorial-docker`** w Artifact Registry — zob. `kubernetes/README.md`)

## Licencja

MIT — zob. `LICENSE`.
