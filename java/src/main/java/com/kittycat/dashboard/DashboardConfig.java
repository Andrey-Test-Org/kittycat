package com.kittycat.dashboard;

import java.net.URI;
import java.time.Duration;

public final class DashboardConfig {
    public final URI backend;
    public final Duration cacheTtl;
    public final int port;

    public DashboardConfig(URI backend, Duration cacheTtl, int port) {
        this.backend = backend;
        this.cacheTtl = cacheTtl;
        this.port = port;
    }

    public static DashboardConfig defaults() {
        return new DashboardConfig(URI.create("http://127.0.0.1:8765"), Duration.ofSeconds(30), 8080);
    }
}
