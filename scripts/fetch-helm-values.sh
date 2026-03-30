#!/usr/bin/env bash
# Zapisuje wartości Helm release Camunda Platform z bieżącego kontekstu kube.
#
# Użycie:
#   ./scripts/fetch-helm-values.sh [namespace] [release_name]
# Domyślnie: namespace=camunda, release=camunda-platform
#
# Wymagane: helm 3, kubectl (kontekst na właściwy GKE).

set -euo pipefail

NS="${1:-camunda}"
REL="${2:-camunda-platform}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT_DIR="${ROOT}/helm"
mkdir -p "${OUT_DIR}"

if ! helm status "${REL}" -n "${NS}" >/dev/null 2>&1; then
  echo "Release '${REL}' w namespace '${NS}' nie znaleziony. Lista release'ów:"
  helm list -n "${NS}"
  exit 1
fi

# Tylko nadpisania użytkownika (jak -f / --set przy install).
helm get values "${REL}" -n "${NS}" -o yaml > "${OUT_DIR}/${REL}-user-values.yaml"

# Pełna obliczona konfiguracja (domyślne chartu + overrides).
helm get values "${REL}" -n "${NS}" --all -o yaml > "${OUT_DIR}/${REL}-all-values.yaml"

echo "Zapisano:"
echo "  ${OUT_DIR}/${REL}-user-values.yaml"
echo "  ${OUT_DIR}/${REL}-all-values.yaml"

if [[ "${FETCH_MANIFEST:-}" == "1" ]]; then
  helm get manifest "${REL}" -n "${NS}" > "${OUT_DIR}/${REL}-manifest.yaml"
  echo "  ${OUT_DIR}/${REL}-manifest.yaml"
fi
