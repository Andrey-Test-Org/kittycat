package com.kittycat.dashboard;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;

public final class StatsClient {
    private final URI base;
    private final HttpClient http;
    private final Cache cache;

    public StatsClient(URI base, Cache cache) {
        this.base = base;
        this.cache = cache;
        this.http = HttpClient.newBuilder().connectTimeout(Duration.ofSeconds(2)).build();
    }

    public String fetchStats() throws Exception {
        String cached = cache.get("stats");
        if (cached != null) return cached;

        HttpRequest req = HttpRequest.newBuilder(base.resolve("/stats"))
            .timeout(Duration.ofSeconds(5))
            .GET()
            .build();
        HttpResponse<String> resp = http.send(req, HttpResponse.BodyHandlers.ofString());
        if (resp.statusCode() != 200) {
            throw new RuntimeException("dashboard /stats returned " + resp.statusCode());
        }
        cache.put("stats", resp.body());
        return resp.body();
    }
}
