# Kontekst dla Cursor (krótko)

Pełny opis repozytorium — w **[README.md](../README.md)**. Tu tylko to, co wygodne przy samej przestrzeni roboczej.

| Co | Gdzie |
|----|-------|
| Mapa katalogów, szybki start | **README.md** |
| Workery GKE (obrazy, job types, rejestr) | **kubernetes/README.md** |
| GKE, Artifact Registry, Helm | **docs/gke-camunda-cheatsheet.md** |
| Eksport Helm values z klastra | **helm/README.md**, **scripts/fetch-helm-values.sh** |
| Przykład user values Camunda (lekki profil) | **helm/camunda-platform-user-values.yaml** |

U autora ścieżka do repo: `/data/projects/camunda8-tutorial`.

Chmura (build workerów i wdrożenie):

```bash
gcloud builds submit --config=ci/cloudbuild-workers.yaml --project=my-camunda8-project .
./scripts/deploy-workers-to-gke.sh
```
