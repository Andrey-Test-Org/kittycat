package com.kittycat;

public final class MathUtils {
    private MathUtils() {}

    public static long gcd(long a, long b) {
        a = Math.abs(a);
        b = Math.abs(b);
        while (b != 0) {
            long t = a % b;
            a = b;
            b = t;
        }
        return a;
    }

    public static boolean isPrime(long n) {
        if (n < 2) return false;
        if (n < 4) return true;
        if (n % 2 == 0) return false;
        for (long i = 3; i * i <= n; i += 2) {
            if (n % i == 0) return false;
        }
        return true;
    }

    public static long factorial(int n) {
        if (n < 0) throw new IllegalArgumentException("n must be non-negative");
        long result = 1;
        for (int i = 2; i <= n; i++) result *= i;
        return result;
    }
}
