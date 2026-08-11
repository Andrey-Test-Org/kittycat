import os
import sys

import pytest

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from anagram import is_anagram


@pytest.mark.parametrize(
    "a, b, expected",
    [
        ("", "", True),
        ("a", "a", True),
        ("a", "b", False),
        ("listen", "silent", True),
        ("Listen", "Silent", True),
        ("conversation", "voices rant on", True),
        ("   ", "", True),
        ("abc", "abcd", False),
        ("aabbcc", "abcabc", True),
        ("aabbcc", "aabbcd", False),
        ("rail safety", "fairy tales", True),
        ("hello", "world", False),
        ("A" * 1000 + "B" * 1000, "B" * 1000 + "A" * 1000, True),
    ],
)
def test_is_anagram(a: str, b: str, expected: bool) -> None:
    assert is_anagram(a, b) is expected


def test_is_anagram_is_symmetric() -> None:
    assert is_anagram("dusty", "study") == is_anagram("study", "dusty")


def test_is_anagram_case_insensitive() -> None:
    assert is_anagram("Tea", "EAT") is True


def test_is_anagram_ignores_spaces() -> None:
    assert is_anagram("the eyes", "they see") is True
