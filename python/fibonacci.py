"""Iterative Fibonacci."""
from __future__ import annotations
import sys


def fibonacci(n: int) -> int:
    if n < 0:
        raise ValueError("n must be non-negative")
    a, b = 0, 1
    for _ in range(n):
        a, b = b, a + b
    return a


def main(argv: list[str]) -> None:
    count = int(argv[1]) if len(argv) > 1 else 10
    for i in range(count + 1):
        print(f"fib({i}) = {fibonacci(i)}")


if __name__ == "__main__":
    main(sys.argv)
