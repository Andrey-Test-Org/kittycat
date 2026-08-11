from __future__ import annotations
from dataclasses import dataclass, field
from datetime import datetime
from typing import List


@dataclass
class KittyEvent:
    cat_id: str
    kind: str
    weight_g: int
    at: datetime


@dataclass
class KittyStats:
    cat_id: str
    total_events: int = 0
    avg_weight_g: float = 0.0
    last_seen: datetime | None = None
    kinds: List[str] = field(default_factory=list)
