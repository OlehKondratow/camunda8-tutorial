"""
Minimalny testowy worker: job type c8jw-python-echo — zwraca zmienne w echo.
Proces w BPMN: Service Task z Type = c8jw-python-echo.
"""
from __future__ import annotations

import asyncio
import logging
import os
import sys
from typing import Any

from pyzeebe import ZeebeWorker, create_insecure_channel

_log = logging.getLogger("echo-demo-worker")

JOB_TYPE = "c8jw-python-echo"


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

    channel = create_insecure_channel(grpc_address=gateway)
    worker = ZeebeWorker(channel)

    @worker.task(task_type=JOB_TYPE)
    async def _echo(**variables: Any) -> dict[str, Any]:
        _log.info("c8jw-python-echo job keys=%s", list(variables.keys()))
        return {"echo": variables, "worker": "c8jw-python-echo"}

    _log.info("listening gateway=%s job_type=%s", gateway, JOB_TYPE)
    await worker.work()


if __name__ == "__main__":
    asyncio.run(main())
