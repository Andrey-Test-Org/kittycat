from datetime import datetime, timedelta

import pytest

from python.dashboard.aggregator import aggregate
from python.dashboard.models import KittyEvent


def _ev(cat: str, kind: str, w: int, offset_min: int = 0) -> KittyEvent:
    return KittyEvent(cat_id=cat, kind=kind, weight_g=w, at=datetime(2026, 1, 1) + timedelta(minutes=offset_min))


def test_aggregate_groups_by_cat() -> None:
    out = aggregate([_ev("a", "nap", 3200), _ev("b", "nap", 4100), _ev("a", "play", 3210, 5)])
    assert set(out) == {"a", "b"}
    assert out["a"].total_events == 2
    assert out["b"].total_events == 1


def test_aggregate_avg_weight() -> None:
    out = aggregate([_ev("a", "x", 100), _ev("a", "x", 200), _ev("a", "x", 300)])
    assert out["a"].avg_weight_g == pytest.approx(200.0)


def test_aggregate_tracks_kinds_in_order() -> None:
    out = aggregate([_ev("a", "nap", 1), _ev("a", "play", 1), _ev("a", "nap", 1)])
    assert out["a"].kinds == ["nap", "play"]


def test_aggregate_empty() -> None:
    assert aggregate([]) == {}
