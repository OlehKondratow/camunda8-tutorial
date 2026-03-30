# Helm values z klastra (Camunda Platform)

W repozytorium może leżeć zacommitowany przykład **`camunda-platform-user-values.yaml`** (ręczna kopia lub szablon startowy). Aktualny zrzut z klastra pobiera skrypt poniżej — **nadpisuje** `*-user-values.yaml` i `*-all-values.yaml`.

Pliki **`camunda-platform-user-values.yaml`** i **`camunda-platform-all-values.yaml`** z aktywnego kontekstu kube:

```bash
chmod +x scripts/fetch-helm-values.sh
kubectl config current-context   # upewnij się, że to Twój GKE
./scripts/fetch-helm-values.sh camunda camunda-platform
```

| Plik | Zawartość |
|------|-----------|
| `camunda-platform-user-values.yaml` | Tylko Twoje nadpisania (`helm get values` bez `--all`) |
| `camunda-platform-all-values.yaml` | Wynik po scaleniu domyślnych chartu i overrides (`--all`) |

Opcjonalnie pełny manifest release (może być bardzo duży):

```bash
FETCH_MANIFEST=1 ./scripts/fetch-helm-values.sh camunda camunda-platform
```

Jeśli inna nazwa release lub namespace — podaj je jako pierwszy i drugi argument.
