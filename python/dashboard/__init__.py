"""kitty-dashboard backend package."""
from .models import KittyEvent, KittyStats
from .aggregator import aggregate

__all__ = ["KittyEvent", "KittyStats", "aggregate"]
