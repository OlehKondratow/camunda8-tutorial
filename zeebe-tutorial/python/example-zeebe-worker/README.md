# example-zeebe-worker (Python)

Minimalny **Job Worker** dla Camunda 8 (Zeebe), odpowiednik `cmd/example-worker` w Go (katalog `zeebe-tutorial`).

## Instalacja

```bash
cd camunda8-tutorial/zeebe-tutorial/python/example-zeebe-worker
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

## Zmienne środowiskowe

Te same co dla workerów Go:

- `ZEEBE_ADDRESS` — np. `127.0.0.1:26500` (lokalny Zeebe) lub adres klastra Camunda SaaS.
- Dla **Camunda SaaS**: `ZEEBE_CLIENT_ID`, `ZEEBE_CLIENT_SECRET`, `ZEEBE_AUTHORIZATION_SERVER_URL` (gdy potrzebne).
- Lokalnie bez OAuth: nie ustawiaj client id/secret — w kodzie dla uproszczenia jest plaintext do lokalnego brokera; dla SaaS włącz TLS wg dokumentacji [pyzeebe](https://github.com/camunda-community-hub/pyzeebe).

## Uruchomienie

```bash
export ZEEBE_ADDRESS=127.0.0.1:26500
python -m app.worker
```

Job types (schemat **`c8jw-*`**):

- **`c8jw-golang`** — gałąź „Obsłużyć automatycznie” w `application.bpmn` (Reunico), zmienne `name`, `amount`.
- **`c8jw-python`** — `../../bpmn/examples/example-task.bpmn`.

## Uwaga

Dla **production** i SaaS lepiej dopracować konfigurację TLS/OAuth wg oficjalnych przykładów klienta Zeebe Python. Ten katalog to **szkoleniowy minimum**.
