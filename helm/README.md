# Helm values с кластера (Camunda Platform)

Файлы **`camunda-platform-user-values.yaml`** и **`camunda-platform-all-values.yaml`** получают из активного kube-контекста скриптом:

```bash
chmod +x scripts/fetch-helm-values.sh
kubectl config current-context   # убедитесь, что это ваш GKE
./scripts/fetch-helm-values.sh camunda camunda-platform
```

| Файл | Содержимое |
|------|------------|
| `camunda-platform-user-values.yaml` | Только ваши переопределения (`helm get values` без `--all`) |
| `camunda-platform-all-values.yaml` | Итог после слияния defaults чарта и overrides (`--all`) |

Опционально полный манифест релиза (может быть очень большим):

```bash
FETCH_MANIFEST=1 ./scripts/fetch-helm-values.sh camunda camunda-platform
```

Если имя релиза или namespace другие — передайте их первым и вторым аргументом.
