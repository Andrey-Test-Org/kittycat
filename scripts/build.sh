#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

echo "[build] Kotlin"
kotlinc "$ROOT/src/Main.kt" "$ROOT/src/Dashboard.kt" -include-runtime -d "$ROOT/out/kittycat.jar"

echo "[build] Java"
(cd "$ROOT/java" && mvn -q -DskipTests package)

echo "[build] Python (lint only)"
python3 -m compileall -q "$ROOT/python"

echo "[build] done"
