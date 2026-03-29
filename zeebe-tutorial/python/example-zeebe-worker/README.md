# example-zeebe-worker (Python)

Минимальный **Job Worker** для Camunda 8 (Zeebe), аналог `cmd/example-worker` на Go (каталог `zeebe-tutorial`).

## Установка

```bash
cd camunda8-tutorial/zeebe-tutorial/python/example-zeebe-worker
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

## Переменные окружения

Те же, что и для Go-воркеров:

- `ZEEBE_ADDRESS` — например `127.0.0.1:26500` (локальный Zeebe) или адрес кластера Camunda SaaS.
- Для **Camunda SaaS**: `ZEEBE_CLIENT_ID`, `ZEEBE_CLIENT_SECRET`, `ZEEBE_AUTHORIZATION_SERVER_URL` (при необходимости).
- Локально без OAuth: не задавайте client id/secret — в коде ниже для простоты используется plaintext к локальному брокеру; для SaaS включите TLS по документации [pyzeebe](https://github.com/camunda-community-hub/pyzeebe).

## Запуск

```bash
export ZEEBE_ADDRESS=127.0.0.1:26500
python -m app.worker
```

Job types:

- **`decision`** — сервисная задача «Обработать автоматически» в `application.bpmn` (Reunico), переменные `name`, `amount`.
- **`example-task`** — `../../bpmn/examples/example-task.bpmn` (относительно этого каталога: `zeebe-tutorial/bpmn/examples/example-task.bpmn`).

## Примечание

Для **production** и SaaS лучше доработать конфиг TLS/OAuth по официальным примерам Zeebe Python client. Этот каталог — **учебный минимум**.
