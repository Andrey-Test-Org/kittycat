from __future__ import annotations
from collections import defaultdict
from typing import Iterable, Dict

from .models import KittyEvent, KittyStats


def aggregate(events: Iterable[KittyEvent]) -> Dict[str, KittyStats]:
    buckets: Dict[str, KittyStats] = {}
    weight_sums: Dict[str, int] = defaultdict(int)
    counts: Dict[str, int] = defaultdict(int)

    for ev in events:
        s = buckets.setdefault(ev.cat_id, KittyStats(cat_id=ev.cat_id))
        s.total_events += 1
        weight_sums[ev.cat_id] += ev.weight_g
        counts[ev.cat_id] += 1
        if s.last_seen is None or ev.at > s.last_seen:
            s.last_seen = ev.at
        if ev.kind not in s.kinds:
            s.kinds.append(ev.kind)

    for cat_id, s in buckets.items():
        s.avg_weight_g = weight_sums[cat_id] / counts[cat_id]
    return buckets
