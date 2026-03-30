# Job workery dla GKE (szkoleniowe)

Namespace: **`c8-tutorial-workers`**.

## Schemat **c8jw** (Camunda 8 job workers)

Jedna nazwa na kontur: **obraz** = **StatefulSet** = label **`app`** = **job type** w BPMN (**Type** przy Service Task). Dla każdego workera tworzony jest headless **Service** `*-headless` (wymóg StatefulSet dla nazwy sieciowej poda).

| Katalog źródeł | Nazwa `c8jw-*` | Język |
|----------------|----------------|-------|
| `form-log-worker/` | `c8jw-python` | Python |
| `go-example-worker/` | `c8jw-golang` | Golang |
| `echo-demo-worker/` | `c8jw-python-echo` | Python |
| `java-example-worker/` | `c8jw-java` | Java |

Nadpisanie typu dla Java: zmienna środowiskowa **`JOB_TYPE`** w StatefulSet.

Repliki: **po 1** na każdy typ workera (`replicas: 1`).

Ponadto: `app.kubernetes.io/component` = `python` | `golang` | `java`.

**Process ID** w Modeler jest dowolny — nie jest powiązany z workerem.

Rejestr: `…/c8-tutorial-docker/c8jw-*`. Szczegóły IAM — **`docs/gke-camunda-cheatsheet.md`** (pkt. 15).

## Build i wdrożenie

```bash
gcloud builds submit --config=ci/cloudbuild-workers.yaml --project=my-camunda8-project .
./scripts/deploy-workers-to-gke.sh
```

## Logi

```bash
kubectl logs -n c8-tutorial-workers -l app=c8jw-python --tail=50 -f
```

## Usunięcie starych zasobów / pełne czyszczenie workerów

Stare **Deployment** (nazwy `tutorial-*`, `c8t-worker-*`), jeśli jeszcze istnieją:

```bash
kubectl delete deployment \
  tutorial-form-log-worker tutorial-go-worker tutorial-echo-demo-worker tutorial-java-worker \
  c8t-worker-example-task c8t-worker-decision c8t-worker-echo-demo c8t-worker-java-demo \
  c8t-worker-python c8t-worker-golang c8t-worker-python-echo c8t-worker-java \
  -n c8-tutorial-workers --ignore-not-found
```

Bieżące **StatefulSet** i headless **Service**:

```bash
kubectl delete statefulset \
  c8jw-python c8jw-golang c8jw-python-echo c8jw-java \
  -n c8-tutorial-workers --ignore-not-found
kubectl delete service \
  c8jw-python-headless c8jw-golang-headless c8jw-python-echo-headless c8jw-java-headless \
  -n c8-tutorial-workers --ignore-not-found
```
