# Camunda 8 — учебный репозиторий

Самостоятельный набор для практики: **Zeebe job workers** (Go и Python), пример BPMN, форма, Docker Compose, шпаргалка по **GKE**.

Ранее лежал в `streamforge-infra/camunda-tutorial/`; личные заметки по-прежнему в `streamforge-infra/temp/` (не в git).

## Состав

| Путь | Назначение |
|------|------------|
| `zeebe-tutorial/` | Go-модуль `zeebe-tutorial` — job types **`decision`**, **`example-task`**; Python (pyzeebe); BPMN `example-task-process` |
| `docker-compose.yaml` | Локальный **Zeebe 8.5** + оба воркера |
| `docs/gke-camunda-cheatsheet.md` | GKE, Helm, port-forward, PVC, очистка |
| `docs/service-tasks-and-gateways-reunico.md` | Service tasks / gateways (Reunico) |
| `examples/application.form` | Пример формы (Camunda) |
| `examples/diagram_1.bpmn` | Дополнительная диаграмма |
| `helm/` | Values Camunda Platform с кластера (см. `helm/README.md`, скрипт `scripts/fetch-helm-values.sh`) |

## Быстрый старт (локально)

Требования: Docker или Podman Compose; опционально Go 1.21+, Python 3.12.

```bash
git clone <url> camunda8-tutorial
cd camunda8-tutorial
docker compose up --build
```

Шлюз: **localhost:26500**. Воркеры слушают **`decision`** и **`example-task`**.

- **Go:** `cd zeebe-tutorial && export ZEEBE_ADDRESS=127.0.0.1:26500 && go run ./cmd/example-worker`
- **Python:** `zeebe-tutorial/python/example-zeebe-worker/README.md`

BPMN **`zeebe-tutorial/bpmn/examples/example-task.bpmn`** (process id `example-task-process`). Процесс **`application.bpmn`** из курсов (ветка «Обработать автоматически», job type **`decision`**) задеплойте из **Camunda Modeler** отдельно.

## Облако и CI

- Шпаргалка: **`docs/gke-camunda-cheatsheet.md`**
- Образы воркеров, Artifact Registry, Cloud Build: репозиторий **`streamforge-infra`** (`kubernetes/tutorial-form-log-worker`, `ci/cloudbuild-tutorial-worker.yaml`, Terraform `live/dev/devops`).

## Лицензия

MIT — см. `LICENSE`.
