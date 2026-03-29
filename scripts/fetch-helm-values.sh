#!/usr/bin/env bash
# Сохраняет Helm values релиза Camunda Platform из текущего kube-контекста.
#
# Использование:
#   ./scripts/fetch-helm-values.sh [namespace] [release_name]
# По умолчанию: namespace=camunda, release=camunda-platform
#
# Требуется: helm 3, kubectl (контекст на нужный GKE).

set -euo pipefail

NS="${1:-camunda}"
REL="${2:-camunda-platform}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT_DIR="${ROOT}/helm"
mkdir -p "${OUT_DIR}"

if ! helm status "${REL}" -n "${NS}" >/dev/null 2>&1; then
  echo "Релиз '${REL}' в namespace '${NS}' не найден. Список релизов:"
  helm list -n "${NS}"
  exit 1
fi

# Только переопределения пользователя (как -f / --set при install).
helm get values "${REL}" -n "${NS}" -o yaml > "${OUT_DIR}/${REL}-user-values.yaml"

# Полная вычисленная конфигурация (defaults чарта + overrides).
helm get values "${REL}" -n "${NS}" --all -o yaml > "${OUT_DIR}/${REL}-all-values.yaml"

echo "Записано:"
echo "  ${OUT_DIR}/${REL}-user-values.yaml"
echo "  ${OUT_DIR}/${REL}-all-values.yaml"

if [[ "${FETCH_MANIFEST:-}" == "1" ]]; then
  helm get manifest "${REL}" -n "${NS}" > "${OUT_DIR}/${REL}-manifest.yaml"
  echo "  ${OUT_DIR}/${REL}-manifest.yaml"
fi
