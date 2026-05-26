"""Basic descriptive statistics."""
from __future__ import annotations
import argparse
from typing import Sequence


def mean(values: Sequence[float]) -> float:
    if not values:
        raise ValueError("values must not be empty")
    return sum(values) / len(values)


def median(values: Sequence[float]) -> float:
    if not values:
        raise ValueError("values must not be empty")
    s = sorted(values)
    mid = len(s) // 2
    if len(s) % 2:
        return s[mid]
    return (s[mid - 1] + s[mid]) / 2


def _parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(description="Print mean and median of input numbers.")
    p.add_argument("numbers", nargs="+", type=float)
    return p.parse_args()


def main() -> None:
    args = _parse_args()
    print(f"mean   = {mean(args.numbers)}")
    print(f"median = {median(args.numbers)}")


if __name__ == "__main__":
    main()
