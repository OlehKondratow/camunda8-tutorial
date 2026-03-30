#!/usr/bin/env bash
# Build obrazu c8jw-python i push do Google Artifact Registry.
#
# Wymagane: Docker, gcloud; użytkownik — rola artifactregistry.writer (lub szersza).
# Przed pierwszym push: gcloud auth login && gcloud auth configure-docker <host>
#
# Zmienne środowiskowe (opcjonalnie):
#   GCP_PROJECT_ID   — domyślnie my-camunda8-project
#   ARTIFACT_REGION  — domyślnie europe-west3
#   ARTIFACT_REPO    — domyślnie c8-tutorial-docker
#   IMAGE_NAME       — domyślnie c8jw-python
#   IMAGE_TAG        — domyślnie git short SHA lub manual-YYYYMMDDhhmm
#
# Przykłady:
#   ./build-push.sh
#   IMAGE_TAG=v1 ./build-push.sh
#   GCP_PROJECT_ID=other ./build-push.sh

set -euo pipefail

GCP_PROJECT_ID="${GCP_PROJECT_ID:-my-camunda8-project}"
ARTIFACT_REGION="${ARTIFACT_REGION:-europe-west3}"
ARTIFACT_REPO="${ARTIFACT_REPO:-c8-tutorial-docker}"
IMAGE_NAME="${IMAGE_NAME:-c8jw-python}"
if [[ -z "${IMAGE_TAG:-}" ]]; then
  if git rev-parse --short HEAD >/dev/null 2>&1; then
    IMAGE_TAG="$(git rev-parse --short HEAD)"
  else
    IMAGE_TAG="manual-$(date +%Y%m%d%H%M)"
  fi
fi

REGISTRY_HOST="${ARTIFACT_REGION}-docker.pkg.dev"
FULL_IMAGE="${REGISTRY_HOST}/${GCP_PROJECT_ID}/${ARTIFACT_REPO}/${IMAGE_NAME}:${IMAGE_TAG}"
LATEST="${REGISTRY_HOST}/${GCP_PROJECT_ID}/${ARTIFACT_REPO}/${IMAGE_NAME}:latest"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Project:       ${GCP_PROJECT_ID}"
echo "Registry:      ${REGISTRY_HOST}"
echo "Image:         ${FULL_IMAGE}"
echo "Also tag:      ${LATEST}"
echo "Build context: ${SCRIPT_DIR}"
echo

gcloud config set project "${GCP_PROJECT_ID}" >/dev/null
echo "Configuring Docker credential helper for ${REGISTRY_HOST}…"
gcloud auth configure-docker "${REGISTRY_HOST}" --quiet

echo "Building…"
docker build -t "${FULL_IMAGE}" -t "${LATEST}" "${SCRIPT_DIR}"

echo "Pushing…"
docker push "${FULL_IMAGE}"
docker push "${LATEST}"

echo
echo "OK: pushed ${FULL_IMAGE}"
