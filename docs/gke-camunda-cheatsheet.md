# Ściąga: GKE + Camunda Platform (krok po kroku)

Notatki o **gcloud / GKE / Helm / kubectl** dla szkoleniowego repozytorium **camunda8-tutorial**.

Podstaw swoje **`PROJECT_ID`**, **billing account**, **email** we wszystkich poleceniach. W przykładach poniżej — **`my-camunda8-project`**; zamień na swój projekt.

Zakładamy zainstalowany [Google Cloud SDK](https://cloud.google.com/sdk/docs/install).

**Przykład billing** (`gcloud billing accounts list` — wybierz swój `ACCOUNT_ID`):

| ACCOUNT_ID | Uwaga |
|------------|-------|
| `YOUR_BILLING_ACCOUNT_ID` | Wstaw wartość z wyjścia `gcloud billing accounts list` |

---

## Spis treści

1. [Autoryzacja i projekt](#1-autoryzacja-i-projekt)  
2. [Billing](#2-billing)  
3. [Włączenie API](#3-włączenie-api)  
4. [Klaster GKE](#4-klaster-gke)  
5. [kubectl i dostęp do klastra](#5-kubectl-i-dostęp-do-klastra)  
6. [IAM dla użytkownika](#6-iam-dla-użytkownika)  
7. [Service Account (opcjonalnie)](#7-service-account-opcjonalnie)  
8. [Camunda Platform (Helm)](#8-camunda-platform-helm)  
9. [Dostęp z maszyny lokalnej (port-forward)](#9-dostęp-z-maszyny-lokalnej-port-forward)  
10. [Limity, Free Tier i koszty](#10-limity-free-tier-i-koszty)  
11. [PVC i pełne sprzątanie środowiska](#11-pvc-i-pełne-sprzątanie-środowiska)  
12. [Desktop Modeler](#12-desktop-modeler)  
13. [Job workery i powiązane repozytoria](#13-job-workery-i-powiązane-repozytoria)  
14. [Helm: lista release'ów i eksport values](#14-helm-lista-releaseów-i-eksport-values)  
15. [CI/CD i obrazy (Artifact Registry)](#15-cicd-i-obrazy-artifact-registry)  
16. [Przydatne polecenia i typowe sytuacje](#16-przydatne-polecenia-i-typowe-sytuacje)

---

## 1. Autoryzacja i projekt

```bash
gcloud auth login
gcloud projects create my-camunda8-project --name="Camunda 8 GKE Project"
gcloud config set project my-camunda8-project
```

Jeśli projekt już istnieje — **nie** wykonuj **`gcloud projects create`**.

Sprawdzenie aktywnego projektu:

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

Bez powiązania billing z projektem GKE i płatne API nie zadziałają.

---

## 3. Włączenie API

```bash
gcloud services enable compute.googleapis.com \
  container.googleapis.com \
  cloudresourcemanager.googleapis.com
```

Dla **Artifact Registry** i **Cloud Build** (obrazy workerów) później:

```bash
gcloud services enable artifactregistry.googleapis.com \
  cloudbuild.googleapis.com \
  --project=my-camunda8-project
```

---

## 4. Klaster GKE

### 4.1 Autopilot (region)

```bash
gcloud container clusters create-auto camunda8-cluster \
  --region=europe-west3 \
  --project=my-camunda8-project
```

Dane uwierzytelniające:

```bash
gcloud container clusters get-credentials camunda8-cluster \
  --region=europe-west3 \
  --project=my-camunda8-project
```

Jeśli odpowiedź **`409 Already exists`** — klaster o tej nazwie już istnieje: wykonaj **`get-credentials`** dla istniejącej nazwy lub usuń stary klaster / użyj innej nazwy.

### 4.2 Standard (strefa) — przykład `camunda-stable`

Wariant z **3 węzłami**, `e2-standard-2`, `pd-balanced` (przy braku limitów zmniejsz `--num-nodes` lub typ VM):

```bash
gcloud container clusters create camunda-stable \
  --zone=europe-west3-c \
  --num-nodes=3 \
  --machine-type=e2-standard-2 \
  --disk-type=pd-balanced \
  --disk-size=25 \
  --project=my-camunda8-project
```

**Alternatywa** z większym zapasem CPU/RAM (jeśli limit pozwala):

```bash
gcloud container clusters create camunda-stable \
  --zone=europe-west3-c \
  --num-nodes=3 \
  --machine-type=e2-standard-4 \
  --disk-type=pd-balanced \
  --disk-size=40 \
  --project=my-camunda8-project
```

Podłączenie kubectl:

```bash
gcloud container clusters get-credentials camunda-stable \
  --zone=europe-west3-c \
  --project=my-camunda8-project

kubectl get nodes
```

**Ważne:** w `gcloud config` może być ustawiona jedna **strefa** (`compute/zone`), a klaster utworzony w innej — dla poleceń do klastra zawsze podawaj **`--zone`** / **`--region`** z opisu klastra.

---

## 5. kubectl i dostęp do klastra

Jeśli `gcloud` zainstalowano przez **apt**, często **`gcloud components install kubectl`** jest **niedostępny**. Zainstaluj plugin i `kubectl` osobno:

```bash
sudo apt update
sudo apt install -y google-cloud-cli-gke-gcloud-auth-plugin
sudo snap install kubectl --classic
# lub binaria: https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/
```

Weryfikacja:

```bash
which kubectl gke-gcloud-auth-plugin
kubectl version --client
kubectl config current-context
```

---

## 6. IAM dla użytkownika

Zamień **`YOUR_EMAIL@example.com`** na swój kont Google.

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

Do uruchamiania **Cloud Build** (`gcloud builds submit`) zwykle potrzebna jest osobna rola, np.:

```bash
gcloud projects add-iam-policy-binding my-camunda8-project \
  --member="user:YOUR_EMAIL@example.com" \
  --role="roles/cloudbuild.builds.editor"
```

Przy **`PERMISSION_DENIED`** na buildach — sprawdź tę rolę lub szerszą u właściciela projektu.

---

## 7. Service Account (opcjonalnie)

```bash
gcloud iam service-accounts create camunda-runner \
  --display-name="Camunda 8 Service Account" \
  --project=my-camunda8-project

gcloud projects add-iam-policy-binding my-camunda8-project \
  --member="serviceAccount:camunda-runner@my-camunda8-project.iam.gserviceaccount.com" \
  --role="roles/container.developer"
```

Podgląd przypisań:

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

### Wariant A: instalacja z pliku values

Z katalogu, gdzie leży **`camunda-values.yaml`** (własny plik lub eksport z `helm/`):

```bash
helm install camunda-platform camunda/camunda-platform \
  -f camunda-values.yaml \
  --namespace camunda \
  --create-namespace
```

### Wariant B: „lekki” profil bez Identity/Keycloak

Aktualne wersje **chartu** i tag **obrazu** — weź z [artifacthub camunda-platform](https://artifacthub.io/packages/helm/camunda/camunda-platform) / [docs Camunda Helm](https://docs.camunda.io/docs/self-managed/deployment/helm/). Przykład dla linii **8.6.x** i chartu **`11.12.2`**:

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

Po instalacji:

```bash
kubectl get pods -n camunda
kubectl get svc -n camunda
helm list -n camunda
```

**Uwaga:** w nowych wydaniach chart może pobierać obrazy **Camunda 8.7 / 8.8** i inne subcharty (Connectors, Optimize itd.) — sprawdź `helm show chart camunda/camunda-platform` i notatki do wybranej wersji.

Hasło pierwszego użytkownika Identity (jeśli Identity jest **włączone** i utworzono secret):

```bash
kubectl get secret camunda-credentials -n camunda \
  -o jsonpath='{.data.identity-firstuser-password}' | base64 -d
echo
```

Ostrzeżenia Helm o przestarzałych polach **`global.documentStore.*`** zestaw z [aktualną dokumentacją Helm Camunda](https://docs.camunda.io/docs/self-managed/deployment/helm/).

### Aktualizacja release i podgląd values

```bash
helm upgrade camunda-platform camunda/camunda-platform -n camunda \
  -f camunda-values.yaml
# lub z tymi samymi --set, co przy install

helm get values camunda-platform -n camunda -o yaml
helm get values camunda-platform -n camunda --all -o yaml
```

Zapis values do plików tego repozytorium: **[§14](#14-helm-lista-releaseów-i-eksport-values)** oraz skrypt **`../scripts/fetch-helm-values.sh`**.

---

## 9. Dostęp z maszyny lokalnej (port-forward)

```bash
kubectl port-forward svc/camunda-platform-operate 8081:80 -n camunda &
kubectl port-forward svc/camunda-platform-tasklist 8082:80 -n camunda &
kubectl port-forward svc/camunda-platform-zeebe-gateway 26500:26500 -n camunda &
```

| Serwis | URL / adres |
|--------|-------------|
| Operate | http://127.0.0.1:8081 |
| Tasklist | http://127.0.0.1:8082 |
| Zeebe gRPC | `127.0.0.1:26500` → `ZEEBE_ADDRESS=127.0.0.1:26500` |

**demo/demo** bywa tylko przy skonfigurowanej **Identity**; w lekkim profilu bez Keycloak logowanie może być inne — zob. NOTES Helm.

Jeśli włączone są **Connectors** (nazwa serwisu może różnić się według wersji chartu):

```bash
# kubectl port-forward svc/camunda-platform-connectors 8086:8080 -n camunda &
```

Jeśli port zajęty: **`address already in use`** — stary `kubectl port-forward` nadal działa; znajdź PID lub zmień lokalny port (`8083:80` itd.).

---

## 10. Limity, Free Tier i koszty

Camunda na GKE zwykle **nie mieści się** w twardych limitach darmowego poziomu **CPU we wszystkich regionach** (historycznie ok. **12 vCPU** łącznie — sprawdź w aktualnej dokumentacji Google).

1. [Billing Console](https://console.cloud.google.com/billing) — ewentualny **Upgrade** na pełne konto płatnicze; łatwiej wtedy **zwiększyć limity**.  
2. [IAM → Quotas](https://console.cloud.google.com/iam-admin/quotas) — parametry jak **CPUs (all regions)**; dla klastra z kilkoma `e2-standard-4` może być potrzebny limit ok. **24 vCPU** i wystarczająco **SSD**.  
3. Sprawdzenie z CLI:

```bash
gcloud compute project-info describe --format="yaml(quotas)" | grep -A 3 "CPUS_ALL_REGIONS"
gcloud compute project-info describe --format="yaml(quotas)" | grep -A 3 "SSD_TOTAL_GB"
```

4. **Budgets & alerts:** [Budgets](https://console.cloud.google.com/billing/budgets).  
5. Trzymaj zasoby w **jednej strefie** (np. `europe-west3-c`), jeśli nie ma wymogu multi-zone — mniej niespodzianek z ruchem i limitami.

---

## 11. PVC i pełne sprzątanie środowiska

PVC **nie zawsze** usuwają się raz z Helm — dyski mogą dalej generować opłaty.

### Sprawdzenie wolumenów

```bash
kubectl get pvc -n camunda
```

Oczekuj **BOUND** dla wolumenów Zeebe / Elasticsearch.

### Krok A: usunięcie release Helm

```bash
helm uninstall camunda-platform -n camunda
```

### Krok B: usunięcie PVC i namespace

```bash
kubectl delete pvc --all -n camunda
kubectl delete namespace camunda
```

### Krok C: usunięcie klastra GKE (jeśli demontujesz całe środowisko)

```bash
gcloud container clusters delete camunda-stable \
  --zone=europe-west3-c \
  --project=my-camunda8-project
```

(Wstaw **nazwę klastra** i **strefę** z `gcloud container clusters list`.)

### Krok D: sprawdzenie „ogonów” w Compute Engine

```bash
gcloud compute disks list --project=my-camunda8-project
gcloud compute addresses list --project=my-camunda8-project
gcloud compute forwarding-rules list --project=my-camunda8-project
```

---

## 12. Desktop Modeler

1. Instalacja (Ubuntu): `sudo snap install camunda-modeler`  
2. **Camunda 8 Self-Managed** przy port-forward do gateway: **Cluster URL** `http://localhost:26500`, **Auth: None** (jeśli brak TLS/OAuth na gateway).  
3. **Deploy** diagramów z Modeler; instancje procesów w **Operate** (jeśli włączone i dostępne z hosta).

---

## 13. Job workery i powiązane repozytoria

| Gdzie | Co |
|-------|-----|
| To repozytorium | **`zeebe-tutorial/`** (Go + Python), **`docker-compose.yaml`** lokalnie |
| Ten rep (`camunda8-tutorial`) | **`kubernetes/`** — kilka job workerów; **`ci/cloudbuild-workers.yaml`** |

Dla Pythona wskazane są **pyzeebe** lub oficjalny klient Go Zeebe; starych snippetów **zeebe-grpc** bez weryfikacji z aktualnym **gateway.proto** nie kopiuj.

---

## 14. Helm: lista release'ów i eksport values

```bash
helm list -n camunda
helm status camunda-platform -n camunda
```

Zapis values do katalogu **`helm/`** tego repozytorium:

```bash
cd /path/to/camunda8-tutorial
chmod +x scripts/fetch-helm-values.sh
./scripts/fetch-helm-values.sh camunda camunda-platform
```

Szczegóły — **`helm/README.md`**.

---

## 15. CI/CD i obrazy (Artifact Registry)

W **`camunda8-tutorial`**:

- **`ci/cloudbuild-workers.yaml`** — build **wszystkich** szkoleniowych obrazów workerów i push do `REGION-docker.pkg.dev/PROJECT/c8-tutorial-docker/...`.
- **`kubernetes/form-log-worker/build-push.sh`** — lokalny `docker build` + push tylko Python form-log worker.
- Repozytorium Docker **`c8-tutorial-docker`** w Artifact Registry utwórz w razie potrzeby:  
  `gcloud artifacts repositories create c8-tutorial-docker --repository-format=docker --location=europe-west3 --project=my-camunda8-project`
- IAM w GCP (reader dla węzłów GKE do tego repo) — konsola lub `gcloud` (poniżej).

### Build obrazu i dostęp do registry (całość obrazu)

Ten sam obraz w rejestrze ma postać:

`REGION-docker.pkg.dev/PROJECT_ID/REPO_ID/IMAGE_NAME:TAG`  
(przykład: `europe-west3-docker.pkg.dev/my-camunda8-project/c8-tutorial-docker/c8jw-python:latest`).

**1) Build w Cloud Build (zalecane dla GKE)**  

- Źródła trafiają do Cloud Build; kontener `docker build`/`push` wykonuje się **w GCP** w imieniu konta serwisowego **`PROJECT_NUMBER@cloudbuild.gserviceaccount.com`**.  
- Na repozytorium potrzebna mu rola **`roles/artifactregistry.writer`** (często już nadana przez `gcloud artifacts repositories add-iam-policy-binding`).  
- Tobie jako deweloperowi na projekt często potrzebna jest **`roles/cloudbuild.builds.editor`**, by wykonywać `gcloud builds submit`.  
- Polecenie z katalogu głównego **`camunda8-tutorial`**:

```bash
gcloud builds submit --config=ci/cloudbuild-workers.yaml --project=my-camunda8-project .
```

Po sukcesie obraz jest już **w registry**; lokalny Docker do push **nie jest konieczny**.

**2) Build na własnej maszynie i `docker push`**  

- Lokalnie: zbuduj **`kubernetes/form-log-worker`** (i ewentualnie pozostałe katalogi w `kubernetes/`), potem tag pod pełną ścieżkę Artifact Registry i `docker push`.  
- **Uwierzytelnianie Docker → Artifact Registry** — to nie GCR! Potrzebne:

```bash
gcloud auth configure-docker europe-west3-docker.pkg.dev
gcloud auth login   # jeśli dawno nie było logowania
```

To dodaje w `~/.docker/config.json` **`credHelpers`** dla `…-docker.pkg.dev` (przy `docker` token pobiera `gcloud`).

Skrypt „zbuduj i wypchnij” (z katalogu workera w **`camunda8-tutorial`**):

```bash
cd /path/to/camunda8-tutorial/kubernetes/form-log-worker
chmod +x build-push.sh
./build-push.sh
# lub: IMAGE_TAG=v1 GCP_PROJECT_ID=my-camunda8-project ./build-push.sh
```

Ręcznie (ten sam obraz):

```bash
cd camunda8-tutorial/kubernetes/form-log-worker
docker build -t europe-west3-docker.pkg.dev/my-camunda8-project/c8-tutorial-docker/c8jw-python:v1 .
docker push europe-west3-docker.pkg.dev/my-camunda8-project/c8-tutorial-docker/c8jw-python:v1
```

**3) Self-hosted CI runner** (agent z dowolną nazwą/labelką typu **`active-runner`**: GitHub Actions, GitLab Runner, Jenkins itd.)  

Build odbywa się **na maszynie runnera** (`docker build` + `docker push`), nie w Cloud Build i nie na laptopie. Dostęp do Artifact Registry to **konto GCP dostępne na tym agencie**, z **`roles/artifactregistry.writer`** (lub szerszą rolą na projekt).

Typowe warianty:

| Sposób | Idea |
|--------|------|
| **Osobne SA + klucz** | Utwórz SA `ci-artifact-push@PROJECT.iam.gserviceaccount.com`, nadaj `artifactregistry.writer`, w sekrecie pipeline włóż JSON key; w kroku job: `gcloud auth activate-service-account --key-file=…` i `gcloud auth configure-docker europe-west3-docker.pkg.dev`, potem `docker push`. |
| **VM w GCP z przypisanym SA** | Instancja runnera ma własne Compute SA z **writer** na repo — na VM często wystarczy `gcloud auth application-default login` / metadata, plus `configure-docker` dla **pkg.dev**. |
| **GitHub Actions + WIF** | Bez długotrwałego klucza: [google-github-actions/auth](https://github.com/google-github-actions/auth) + workload identity federation do GCP, dalej `docker` push z uzyskanym dostępem. |

Ważne: `credHelpers` w `~/.docker/config.json` na runnerze muszą obejmować **`…-docker.pkg.dev`**, jak w pkt. **2** — inaczej push do Artifact Registry zachowa się tak, jak gdyby w pliku były tylko **`gcr.io`**.

**4) GKE pobiera obraz przy `imagePull` (to nie runner ani laptop)**  

- Kubelet na węźle używa **poświadczeń węzła** (często **domyślne Compute Engine SA** `PROJECT_NUMBER-compute@developer.gserviceaccount.com`) lub osobnego SA / Workload Identity.  
- Do **repozytorium** (lub projektu) temu SA potrzebna jest **`roles/artifactregistry.reader`**.  
- **`~/.docker/config.json`** na runnerze lub PC **nie wpływa na pod** — kubelet ma własne mechanizmy i IAM.

**Wyjaśnienie: VM-węzeł GKE a Kaniko w podzie — to różne rzeczy**

| Co | Gdzie wykonuje się | Rola |
|----|-------------------|------|
| **Pod roboczy** pobiera gotowy obraz `image: …pkg.dev/...` | **VM węzła** GKE. **Kubelet** na tym węźle pobiera warstwy z Artifact Registry w imieniu **SA węzła** (lub przez **Workload Identity**, gdy pod ma Kubernetes SA powiązane z GCP SA). | To **nie** Kaniko i nie „runner w podzie” w sensie CI — zwykły **imagePull**. |
| **Build obrazu w klastrze** | Często osobny **Pod** (np. **Kaniko**, **Buildah**, czasem **Docker-in-Docker**). | Kontener **buduje** Dockerfile i robi **push** do registry; potrzebuje **poświadczeń writer** (secret z `config.json`, SA + WIF itd.). Na niego **nie** rozciąga się logika „węzeł sam reader” — jawne przekazanie dostępu w manifeście Job/Pod. |
| **Self-hosted runner jako Pod na GKE** | Pod na **tej samej VM-węźle** co pozostałe pody. | Wewnątrz często montuje się **docker.sock** węzła (ostrożnie z bezpieczeństwem), albo uruchamia **Kaniko** osobnym krokiem, albo wysyła build do **Cloud Build**. Tożsamość do **push** i tak trzeba skonfigurować w jobie, jak w pkt. **3**. |

Podsumowanie: **pull** obrazów workerów w klastrze to **węzeł (VM) + IAM reader** na repozytorium **`c8-tutorial-docker`**. **Kaniko** dotyczy **build i push z wnętrza klastra**, osobna konfiguracja sekretów/SA — nie mylić z imagePull obciążenia roboczego.

Jeśli w `~/.docker/config.json` są tylko **`gcr.io`** / **`eu.gcr.io`**, bez wpisu dla **`…-docker.pkg.dev`**, lokalny `docker push` do Artifact Registry nie zadziała — zob. pkt. **2** (`gcloud auth configure-docker …`).

Po buildzie tagi obrazów są ustawione w **`kubernetes/*/k8s/statefulset.yaml`**.

### Wdrożenie workerów do klastra (`camunda8-tutorial`)

Manifesty: **`kubernetes/c8-tutorial-workers/namespace.yaml`** i **`kubernetes/*/k8s/statefulset.yaml`** (cztery StatefulSet po jednej replice i headless Service na każdy). Lista workerów — **`kubernetes/README.md`**.

Upewnij się, że **obrazy są zbudowane** (`gcloud builds submit --config=ci/cloudbuild-workers.yaml …`) i że węzły GKE mają **`roles/artifactregistry.reader`** na repozytorium **`c8-tutorial-docker`**.

```bash
kubectl config use-context <twój-gke-context>
chmod +x scripts/deploy-workers-to-gke.sh
./scripts/deploy-workers-to-gke.sh
```

**ImagePullBackOff:** brak obrazu w registry lub brak uprawnień pull u SA węzła — sprawdź `kubectl describe pod -n c8-tutorial-workers …`.

**CrashLoop / błąd Zeebe:** sprawdź, że w **`camunda`** działa **`svc/camunda-platform-zeebe-gateway:26500`** i że nazwa serwisu zgadza się z wersją Helm (`kubectl get svc -n camunda`).

---

## 16. Przydatne polecenia i typowe sytuacje

```bash
gcloud container clusters list --project=my-camunda8-project
kubectl get pods -n camunda
kubectl get svc -n camunda
helm list -n camunda
```

**Długi start:** zaraz po `helm install` pody mogą być **Pending** / **0/1 Running** — patrz zdarzenia:

```bash
kubectl describe pod -n camunda <pod-name>
kubectl get events -n camunda --sort-by='.lastTimestamp'
```

Przykładowe wyjście podczas oczekiwania (nazwy/wiek mogą się różnić):

```text
kubectl get pods -n camunda -w
NAME                                           READY   STATUS    RESTARTS   AGE
camunda-platform-elasticsearch-master-0        0/1     Pending   0          2m
camunda-platform-zeebe-0                       0/1     Running   0          2m
camunda-platform-zeebe-gateway-xxxx            0/1     Running   0          2m
```

**409 przy tworzeniu klastra** — nazwa zajęta; użyj **`get-credentials`** do istniejącego lub innej nazwy.

---

*Trzymaj tę ściągę zgodnie z praktyką projektu: po usunięciu klastra sprawdzaj **dyski i IP**, żeby nie zbierać niepotrzebnych opłat.*
