"""
Zeebe job workers (schemat c8jw-*):
  - c8jw-golang   — gałąź „Obsłużyć automatycznie” w application.bpmn (Type = c8jw-golang)
  - c8jw-python   — zeebe-tutorial/bpmn/examples/example-task.bpmn

Zmienne dla c8jw-golang (przykład): {"name": "Elis", "amount": 3000}
"""
from __future__ import annotations

import asyncio
import logging
import os
from typing import Union

from pyzeebe import ZeebeWorker, create_insecure_channel

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("example-zeebe-worker")


async def main() -> None:
    gateway = os.environ.get("ZEEBE_ADDRESS", "127.0.0.1:26500")
    channel = create_insecure_channel(grpc_address=gateway)
    worker = ZeebeWorker(channel)

    @worker.task(task_type="c8jw-golang")
    async def golang_branch_task(name: str = "applicant", amount: Union[int, float] = 0) -> dict:
        n = (name or "applicant").strip() or "applicant"
        amt = float(amount)
        msg = f"automatic processing for {n}, amount={amt}"
        logger.info("c8jw-golang job name=%s amount=%s", n, amt)
        return {"processed": True, "message": msg, "mode": "automatic"}

    @worker.task(task_type="c8jw-python")
    async def python_example_task(name: str = "world") -> dict:
        n = (name or "world").strip() or "world"
        logger.info("c8jw-python name=%s", n)
        return {"message": f"Hello, {n}!", "ok": True}

    logger.info("workers started gateway=%s types=c8jw-golang,c8jw-python", gateway)
    await worker.work()


if __name__ == "__main__":
    asyncio.run(main())
