"""Primality check and prime generator."""
from __future__ import annotations
import sys
from typing import Iterator


def is_prime(n: int) -> bool:
    if n < 2:
        return False
    if n < 4:
        return True
    if n % 2 == 0:
        return False
    i = 3
    while i * i <= n:
        if n % i == 0:
            return False
        i += 2
    return True


def primes() -> Iterator[int]:
    n = 2
    while True:
        if is_prime(n):
            yield n
        n += 1


if __name__ == "__main__":
    count = int(sys.argv[1]) if len(sys.argv) > 1 else 20
    gen = primes()
    for _ in range(count):
        print(next(gen))
