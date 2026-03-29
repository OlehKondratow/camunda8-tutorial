# Шпаргалка: GKE + Camunda Platform (пошагово)

Заметки по **gcloud / GKE / Helm / kubectl** для учебного репозитория **camunda8-tutorial** (рядом со `streamforge-infra` или отдельный clone).

Подставьте свои **`PROJECT_ID`**, **billing account**, **email** во все команды. В примерах ниже — **`my-camunda8-project`**; замените на свой проект.

Предполагается установленный [Google Cloud SDK](https://cloud.google.com/sdk/docs/install).

**Пример billing** (`gcloud billing accounts list` — возьмите свой `ACCOUNT_ID`):

| ACCOUNT_ID | Примечание |
|------------|------------|
| `YOUR_BILLING_ACCOUNT_ID` | Подставьте строку из вывода `gcloud billing accounts list` |

---

## Содержание

1. [Авторизация и проект](#1-авторизация-и-проект)  
2. [Billing](#2-billing)  
3. [Включение API](#3-включение-api)  
4. [Кластер GKE](#4-кластер-gke)  
5. [kubectl и доступ к кластеру](#5-kubectl-и-доступ-к-кластеру)  
6. [IAM для пользователя](#6-iam-для-пользователя)  
7. [Service Account (опционально)](#7-service-account-опционально)  
8. [Camunda Platform (Helm)](#8-camunda-platform-helm)  
9. [Доступ с локальной машины (port-forward)](#9-доступ-с-локальной-машины-port-forward)  
10. [Квоты, Free Tier и затраты](#10-квоты-free-tier-и-затраты)  
11. [PVC и полная очистка стенда](#11-pvc-и-полная-очистка-стенда)  
12. [Desktop Modeler](#12-desktop-modeler)  
13. [Job workers и связанные репозитории](#13-job-workers-и-связанные-репозитории)  
14. [Helm: список релизов и экспорт values](#14-helm-список-релизов-и-экспорт-values)  
15. [CI/CD и образы (Artifact Registry)](#15-cicd-и-образы-artifact-registry)  
16. [Полезные команды и типичные ситуации](#16-полезные-команды-и-типичные-ситуации)

---

## 1. Авторизация и проект

```bash
gcloud auth login
gcloud projects create my-camunda8-project --name="Camunda 8 GKE Project"
gcloud config set project my-camunda8-project
```

Если проект уже есть — **`gcloud projects create`** не выполняйте.

Проверка активного проекта:

```bash
gcloud config get-value project
```

---

## 2. Billing

```bash
gcloud billing accounts list

gcloud billing projects link my-camunda8-project \
  --billing-account=YOUR_BILLING_ACCOUNT_ID
```

Без привязки billing к проекту GKE и платные API работать не будут.

---

## 3. Включение API

```bash
gcloud services enable compute.googleapis.com \
  container.googleapis.com \
  cloudresourcemanager.googleapis.com
```

Для **Artifact Registry** и **Cloud Build** (образы воркеров) позже:

```bash
gcloud services enable artifactregistry.googleapis.com \
  cloudbuild.googleapis.com \
  --project=my-camunda8-project
```

---

## 4. Кластер GKE

### 4.1 Autopilot (регион)

```bash
gcloud container clusters create-auto camunda8-cluster \
  --region=europe-west3 \
  --project=my-camunda8-project
```

Учётные данные:

```bash
gcloud container clusters get-credentials camunda8-cluster \
  --region=europe-west3 \
  --project=my-camunda8-project
```

Если ответ **`409 Already exists`** — кластер с таким именем уже есть: выполните **`get-credentials`** для существующего имени или удалите старый кластер / задайте другое имя.

### 4.2 Standard (зона) — пример `camunda-stable`

Вариант с **3 нодами**, `e2-standard-2`, `pd-balanced` (при нехватке квот уменьшите `--num-nodes` или тип ВМ):

```bash
gcloud container clusters create camunda-stable \
  --zone=europe-west3-c \
  --num-nodes=3 \
  --machine-type=e2-standard-2 \
  --disk-type=pd-balanced \
  --disk-size=25 \
  --project=my-camunda8-project
```

**Альтернатива** с большим запасом по CPU/RAM (если квота позволяет):

```bash
gcloud container clusters create camunda-stable \
  --zone=europe-west3-c \
  --num-nodes=3 \
  --machine-type=e2-standard-4 \
  --disk-type=pd-balanced \
  --disk-size=40 \
  --project=my-camunda8-project
```

Подключение kubectl:

```bash
gcloud container clusters get-credentials camunda-stable \
  --zone=europe-west3-c \
  --project=my-camunda8-project

kubectl get nodes
```

**Важно:** в `gcloud config` может быть указана одна **зона** (`compute/zone`), а кластер создан в другой — для команд к кластеру всегда передавайте **`--zone`** / **`--region`** из описания кластера.

---

## 5. kubectl и доступ к кластеру

Если `gcloud` установлен через **apt**, то **`gcloud components install kubectl`** часто **недоступен**. Установите плагин и `kubectl` отдельно:

```bash
sudo apt update
sudo apt install -y google-cloud-cli-gke-gcloud-auth-plugin
sudo snap install kubectl --classic
# или бинарь: https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/
```

Проверка:

```bash
which kubectl gke-gcloud-auth-plugin
kubectl version --client
kubectl config current-context
```

---

## 6. IAM для пользователя

Замените **`YOUR_EMAIL@example.com`** на свой Google-аккаунт.

```bash
gcloud projects add-iam-policy-binding my-camunda8-project \
  --member="user:YOUR_EMAIL@example.com" \
  --role="roles/container.admin"

gcloud projects add-iam-policy-binding my-camunda8-project \
  --member="user:YOUR_EMAIL@example.com" \
  --role="roles/compute.viewer"

gcloud projects add-iam-policy-binding my-camunda8-project \
  --member="user:YOUR_EMAIL@example.com" \
  --role="roles/iam.serviceAccountUser"
```

Для запуска **Cloud Build** (`gcloud builds submit`) обычно нужна отдельная роль, например:

```bash
gcloud projects add-iam-policy-binding my-camunda8-project \
  --member="user:YOUR_EMAIL@example.com" \
  --role="roles/cloudbuild.builds.editor"
```

При **`PERMISSION_DENIED`** на сборках — проверьте эту роль или более широкую у владельца проекта.

---

## 7. Service Account (опционально)

```bash
gcloud iam service-accounts create camunda-runner \
  --display-name="Camunda 8 Service Account" \
  --project=my-camunda8-project

gcloud projects add-iam-policy-binding my-camunda8-project \
  --member="serviceAccount:camunda-runner@my-camunda8-project.iam.gserviceaccount.com" \
  --role="roles/container.developer"
```

Просмотр привязок:

```bash
gcloud projects get-iam-policy my-camunda8-project \
  --flatten="bindings[].members" \
  --format="table(bindings.role, bindings.members)"
```

---

## 8. Camunda Platform (Helm)

```bash
helm repo add camunda https://helm.camunda.io
helm repo update
```

### Вариант A: установка из файла values

Из каталога, где лежит **`camunda-values.yaml`** (например `streamforge-infra/camunda.yaml` или свой файл):

```bash
helm install camunda-platform camunda/camunda-platform \
  -f camunda-values.yaml \
  --namespace camunda \
  --create-namespace
```

### Вариант B: «лёгкий» профиль без Identity/Keycloak

Актуальные версии **чарта** и тег **образа** возьмите с [artifacthub camunda-platform](https://artifacthub.io/packages/helm/camunda/camunda-platform) / [docs Camunda Helm](https://docs.camunda.io/docs/self-managed/deployment/helm/). Пример для линии **8.6.x** и chart **`11.12.2`**:

```bash
helm install camunda-platform camunda/camunda-platform \
  --namespace camunda --create-namespace \
  --version 11.12.2 \
  --set global.image.tag=8.6.1 \
  --set identity.enabled=false \
  --set identityKeycloak.enabled=false \
  --set keycloak.enabled=false \
  --set optimize.enabled=false \
  --set connectors.enabled=false \
  --set zeebe.clusterSize=1 \
  --set zeebe.persistence.size=20Gi \
  --set elasticsearch.master.replicaCount=1 \
  --set elasticsearch.master.persistence.size=20Gi \
  --set operate.enabled=true \
  --set tasklist.enabled=true
```

После установки:

```bash
kubectl get pods -n camunda
kubectl get svc -n camunda
helm list -n camunda
```

**Примечание:** в свежих релизах чарт может тянуть образы **Camunda 8.7 / 8.8** и другие подчарты (Connectors, Optimize и т.д.) — смотрите `helm show chart camunda/camunda-platform` и заметки к выбранной версии.

Пароль первого пользователя Identity (если Identity **включён** и секрет создан):

```bash
kubectl get secret camunda-credentials -n camunda \
  -o jsonpath='{.data.identity-firstuser-password}' | base64 -d
echo
```

Предупреждения Helm про устаревшие поля **`global.documentStore.*`** сверяйте с [актуальной документацией Helm Camunda](https://docs.camunda.io/docs/self-managed/deployment/helm/).

### Обновление релиза и просмотр values

```bash
helm upgrade camunda-platform camunda/camunda-platform -n camunda \
  -f camunda-values.yaml
# или с теми же --set, что при install

helm get values camunda-platform -n camunda -o yaml
helm get values camunda-platform -n camunda --all -o yaml
```

Сохранить values в файлы этого репозитория: **[§14](#14-helm-список-релизов-и-экспорт-values)** и скрипт **`../scripts/fetch-helm-values.sh`**.

---

## 9. Доступ с локальной машины (port-forward)

```bash
kubectl port-forward svc/camunda-platform-operate 8081:80 -n camunda &
kubectl port-forward svc/camunda-platform-tasklist 8082:80 -n camunda &
kubectl port-forward svc/camunda-platform-zeebe-gateway 26500:26500 -n camunda &
```

| Сервис | URL / адрес |
|--------|---------------|
| Operate | http://127.0.0.1:8081 |
| Tasklist | http://127.0.0.1:8082 |
| Zeebe gRPC | `127.0.0.1:26500` → `ZEEBE_ADDRESS=127.0.0.1:26500` |

**demo/demo** бывает только при настроенной **Identity**; в «лёгком» профиле без Keycloak вход может отличаться — см. NOTES Helm.

Если включены **Connectors** (имя сервиса может отличаться по версии чарта):

```bash
# kubectl port-forward svc/camunda-platform-connectors 8086:8080 -n camunda &
```

Если порт занят: **`address already in use`** — старый `kubectl port-forward` ещё работает; найдите PID или смените локальный порт (`8083:80` и т.д.).

---

## 10. Квоты, Free Tier и затраты

Camunda на GKE обычно **не укладывается** в жёсткие лимиты бесплатного уровня по **CPU на все регионы** (исторически порядка **12 vCPU** суммарно — уточняйте в актуальной документации Google).

1. [Billing Console](https://console.cloud.google.com/billing) — при необходимости **Upgrade** на полный платёжный аккаунт; так проще запрашивать **увеличение квот**.  
2. [IAM → Quotas](https://console.cloud.google.com/iam-admin/quotas) — параметры вроде **CPUs (all regions)**; для кластера из нескольких `e2-standard-4` может понадобиться лимит порядка **24 vCPU** и достаточно **SSD**.  
3. Проверка из CLI:

```bash
gcloud compute project-info describe --format="yaml(quotas)" | grep -A 3 "CPUS_ALL_REGIONS"
gcloud compute project-info describe --format="yaml(quotas)" | grep -A 3 "SSD_TOTAL_GB"
```

4. **Budgets & alerts:** [Budgets](https://console.cloud.google.com/billing/budgets).  
5. Держите ресурсы в **одной зоне** (например `europe-west3-c`), если нет задачи на multi-zone — меньше сюрпризов по трафику и квотам.

---

## 11. PVC и полная очистка стенда

PVC **не всегда** удаляются вместе с Helm — диски могут продолжать тарифицироваться.

### Проверка томов

```bash
kubectl get pvc -n camunda
```

Ожидайте **BOUND** для томов Zeebe / Elasticsearch.

### Шаг A: удалить релиз Helm

```bash
helm uninstall camunda-platform -n camunda
```

### Шаг B: удалить PVC и namespace

```bash
kubectl delete pvc --all -n camunda
kubectl delete namespace camunda
```

### Шаг C: удалить кластер GKE (если сносите весь стенд)

```bash
gcloud container clusters delete camunda-stable \
  --zone=europe-west3-c \
  --project=my-camunda8-project
```

(Подставьте **имя кластера** и **зону** из `gcloud container clusters list`.)

### Шаг D: проверка «хвостов» в Compute Engine

```bash
gcloud compute disks list --project=my-camunda8-project
gcloud compute addresses list --project=my-camunda8-project
gcloud compute forwarding-rules list --project=my-camunda8-project
```

---

## 12. Desktop Modeler

1. Установка (Ubuntu): `sudo snap install camunda-modeler`  
2. **Camunda 8 Self-Managed** при port-forward на gateway: **Cluster URL** `http://localhost:26500`, **Auth: None** (если нет TLS/OAuth на gateway).  
3. **Deploy** диаграммы из Modeler; экземпляры процессов смотрите в **Operate** (если включён и доступен с хоста).

---

## 13. Job workers и связанные репозитории

| Где | Что |
|-----|-----|
| Этот репозиторий | **`zeebe-tutorial/`** (Go + Python), **`docker-compose.yaml`** локально |
| `streamforge-infra` | **`kubernetes/tutorial-form-log-worker/`**, **`ci/cloudbuild-tutorial-worker.yaml`**, Terraform **`live/dev/devops`** (Artifact Registry) |

Для Python предпочтительны **pyzeebe** или официальный Go-клиент Zeebe; старые сниппеты **zeebe-grpc** без сверки с текущим **gateway.proto** не копируйте.

---

## 14. Helm: список релизов и экспорт values

```bash
helm list -n camunda
helm status camunda-platform -n camunda
```

Сохранить values в каталог **`helm/`** этого репозитория:

```bash
cd /path/to/camunda8-tutorial
chmod +x scripts/fetch-helm-values.sh
./scripts/fetch-helm-values.sh camunda camunda-platform
```

Подробности — **`helm/README.md`**.

---

## 15. CI/CD и образы (Artifact Registry)

В **`streamforge-infra`**:

- Terraform **`terraform/live/dev/devops`** — Docker-репозиторий и IAM для Cloud Build / pull с GKE.  
- **`ci/cloudbuild-tutorial-worker.yaml`** — пример сборки воркера и push в `REGION-docker.pkg.dev/PROJECT/REPO/...`.

После push подставьте образ в **Deployment** воркера (`image: .../tutorial-form-log-worker:tag`).

---

## 16. Полезные команды и типичные ситуации

```bash
gcloud container clusters list --project=my-camunda8-project
kubectl get pods -n camunda
kubectl get svc -n camunda
helm list -n camunda
```

**Долгий старт:** сразу после `helm install` поды могут быть **Pending** / **0/1 Running** — смотрите события:

```bash
kubectl describe pod -n camunda <pod-name>
kubectl get events -n camunda --sort-by='.lastTimestamp'
```

Пример вывода при ожидании (имена/возраст могут отличаться):

```text
kubectl get pods -n camunda -w
NAME                                           READY   STATUS    RESTARTS   AGE
camunda-platform-elasticsearch-master-0        0/1     Pending   0          2m
camunda-platform-zeebe-0                       0/1     Running   0          2m
camunda-platform-zeebe-gateway-xxxx            0/1     Running   0          2m
```

**409 при создании кластера** — имя занято; используйте **`get-credentials`** к существующему или другое имя.

---

*Держите эту шпаргалку согласованной с практикой по проекту: после удаления кластера проверяйте **диски и IP**, чтобы не копить платёж.*
