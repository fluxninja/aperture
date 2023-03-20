import inspect
from typing import Callable, TypeVar

TWrappedReturn = TypeVar("TWrappedReturn")


async def run_fn(fn: Callable[..., TWrappedReturn], *args, **kwargs) -> TWrappedReturn:
    """Run a function or coroutine."""
    # We want to support both sync and async functions
    # https://stackoverflow.com/questions/44169998/how-to-create-a-python-decorator-that-can-wrap-either-coroutine-or-function
    if inspect.iscoroutinefunction(fn):
        return await fn(*args, **kwargs)
    else:
        return fn(*args, **kwargs)
