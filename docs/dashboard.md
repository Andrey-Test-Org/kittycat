# kitty-dashboard

End-to-end sketch of the multi-language dashboard demo.

## Components

- `python/dashboard/` — backend HTTP server, aggregator, models
- `java/src/main/java/com/kittycat/dashboard/` — JVM client with caching
- `src/Dashboard.kt` — shared aggregation in Kotlin
- `web/dashboard/` — vanilla HTML/CSS/JS frontend
- `scripts/` — build pipeline and seed SQL
- `config/dashboard.yaml` — runtime config

## Data flow

```
[kitty_event SQL] -> [python aggregator] -> [/stats HTTP] -> [JS dashboard]
                                            \-> [Java StatsClient (cached)]
```

## Local dev

```
bash scripts/build.sh
python3 -m python.dashboard.server
open web/dashboard/index.html
```

## Status

This PR scaffolds the structure. Real persistence, auth, and CI wiring land in follow-ups.
