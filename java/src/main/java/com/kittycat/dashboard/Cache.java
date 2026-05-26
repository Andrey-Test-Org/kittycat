package com.kittycat.dashboard;

import java.time.Duration;
import java.time.Instant;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public final class Cache {
    private static final class Entry {
        final String value;
        final Instant expiresAt;
        Entry(String value, Instant expiresAt) {
            this.value = value;
            this.expiresAt = expiresAt;
        }
    }

    private final Map<String, Entry> data = new ConcurrentHashMap<>();
    private final Duration ttl;

    public Cache(Duration ttl) {
        this.ttl = ttl;
    }

    public String get(String key) {
        Entry e = data.get(key);
        if (e == null) return null;
        if (Instant.now().isAfter(e.expiresAt)) {
            data.remove(key);
            return null;
        }
        return e.value;
    }

    public void put(String key, String value) {
        data.put(key, new Entry(value, Instant.now().plus(ttl)));
    }

    public int size() {
        return data.size();
    }
}
