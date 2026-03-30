#!/usr/bin/env bash
# Wdrożenie workerów: nazwy obrazów / StatefulSet / app / typ joba Zeebe — c8jw-* (zob. kubernetes/README.md).

set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
K8S="$ROOT/kubernetes"

need() {
  [[ -f "$1" ]] || { echo "Brak pliku $1"; exit 1; }
}

need "$K8S/c8-tutorial-workers/namespace.yaml"
need "$K8S/form-log-worker/k8s/statefulset.yaml"
need "$K8S/go-example-worker/k8s/statefulset.yaml"
need "$K8S/echo-demo-worker/k8s/statefulset.yaml"
need "$K8S/java-example-worker/k8s/statefulset.yaml"

echo "kubectl context: $(kubectl config current-context)"
kubectl get ns camunda >/dev/null 2>&1 || { echo "Namespace camunda nie istnieje — Camunda Platform nie wdrożona?"; exit 1; }

# Migracja ze starych Deployment (nazwy zgodne ze StatefulSet — usunąć poprzedni rodzaj zasobu).
kubectl delete deployment c8jw-python c8jw-golang c8jw-python-echo c8jw-java \
  -n c8-tutorial-workers --ignore-not-found

kubectl apply -f "$K8S/c8-tutorial-workers/namespace.yaml"
kubectl apply -f "$K8S/form-log-worker/k8s/statefulset.yaml"
kubectl apply -f "$K8S/go-example-worker/k8s/statefulset.yaml"
kubectl apply -f "$K8S/echo-demo-worker/k8s/statefulset.yaml"
kubectl apply -f "$K8S/java-example-worker/k8s/statefulset.yaml"

echo "Oczekiwanie na rollout…"
for s in c8jw-python c8jw-golang c8jw-python-echo c8jw-java; do
  kubectl rollout status "statefulset/$s" -n c8-tutorial-workers --timeout=180s
done

kubectl get sts,svc,pods -n c8-tutorial-workers
for app in c8jw-python c8jw-golang c8jw-python-echo c8jw-java; do
  echo "--- logs $app ---"
  kubectl logs -n c8-tutorial-workers -l "app=$app" --tail=20 || true
done
