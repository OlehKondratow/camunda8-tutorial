# Contributing

Thank you for improving this tutorial repository. The goal is to keep examples **small, runnable, and consistent** with Camunda 8 job-worker semantics.

## Before you open a PR

1. **Run CI checks locally** (mirrors `.github/workflows/ci.yml`):
   - **Go:** `cd zeebe-tutorial && go vet ./... && go build -o /dev/null ./cmd/example-worker`
   - **Python:** `python3 -m compileall kubernetes/echo-demo-worker/app kubernetes/form-log-worker/app zeebe-tutorial/python/example-zeebe-worker/app`
   - **Java:** `mvn -B -f kubernetes/java-example-worker/pom.xml package -DskipTests`
2. If you change **BPMN / DMN / forms**, update the relevant section in `README.md` or `zeebe-tutorial/bpmn-and-dmn-bundles.md` so deploy order and process IDs stay documented.
3. Prefer **English** for code comments, log messages, and identifiers; **Polish** is fine for tutorial prose in `README.md` / `docs/`.

## Conventions

- **Job types** follow the `c8jw-*` prefix used across Go, Python, Java, and Kubernetes manifests (see `kubernetes/README.md`).
- Avoid breaking renames of env vars consumed by `docker-compose.yaml` without updating Compose and docs together.
- New cloud or CI steps should remain secrets-free: no project IDs, keys, or registry URLs that only you can access — use placeholders and point to existing docs.

## Pull request description

A good PR explains **what** changed and **why**, and mentions how you verified (e.g. “Compose up”, “GKE deploy”, or “CI only”). Screenshots are welcome for Modeler/Operate changes.
