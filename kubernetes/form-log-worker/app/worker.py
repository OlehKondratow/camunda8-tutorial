"""
Testowy Zeebe job worker: przyjmuje job, rozbiera zmienne jako JSON i zapisuje je na stdout (kubectl logs).

Zmienne środowiskowe:
  ZEEBE_ADDRESS — wymagane, np. camunda-platform-zeebe-gateway.camunda.svc.cluster.local:26500
  JOB_TYPES     — po przecinku, domyślnie c8jw-python
  LOG_LEVEL     — INFO, DEBUG, …
"""
from __future__ import annotations

import asyncio
import json
import logging
import os
import sys
from typing import Any

from pyzeebe import ZeebeWorker, create_insecure_channel

_log = logging.getLogger("form-log-worker")

DEFAULT_JOB_TYPES = "c8jw-python"


def _parse_job_types() -> list[str]:
    raw = os.environ.get("JOB_TYPES", DEFAULT_JOB_TYPES).strip()
    return [t.strip() for t in raw.split(",") if t.strip()]


def _log_variables(job_type: str, variables: dict[str, Any]) -> None:
    """Czytelny zapis w logu poda (jak pola formularza / procesu)."""
    _log.info("======== job task_type=%s ========", job_type)
    if not variables:
        _log.info("variables: <empty>")
        return
    try:
        pretty = json.dumps(variables, ensure_ascii=False, indent=2, default=str)
        _log.info("variables (JSON):\n%s", pretty)
    except (TypeError, ValueError) as e:
        _log.warning("could not serialize variables: %s; raw=%r", e, variables)
        return
    for key in sorted(variables.keys()):
        val = variables[key]
        _log.info("field[%s] = %r (%s)", key, val, type(val).__name__)


def _configure_logging() -> None:
    level_name = os.environ.get("LOG_LEVEL", "INFO").upper()
    level = getattr(logging, level_name, logging.INFO)
    logging.basicConfig(
        level=level,
        format="%(asctime)s %(levelname)s %(message)s",
        stream=sys.stdout,
        force=True,
    )


async def main() -> None:
    _configure_logging()
    gateway = os.environ.get("ZEEBE_ADDRESS", "").strip()
    if not gateway:
        _log.error("ZEEBE_ADDRESS is required")
        raise SystemExit(1)

    types = _parse_job_types()
    channel = create_insecure_channel(grpc_address=gateway)
    worker = ZeebeWorker(channel)

    for task_type in types:

        def register(tt: str = task_type) -> None:
            @worker.task(task_type=tt)
            async def _handler(**variables: Any) -> dict[str, Any]:
                payload = dict(variables)
                _log.info("job activated task_type=%s keys=%s", tt, list(payload.keys()))
                _log_variables(tt, payload)
                _log.info("completing job task_type=%s with logged=true", tt)
                return {"logged": True, "worker": "c8jw-python"}

        register()

    _log.info(
        "listening gateway=%s job_types=%s",
        gateway,
        types,
    )
    await worker.work()


if __name__ == "__main__":
    asyncio.run(main())
