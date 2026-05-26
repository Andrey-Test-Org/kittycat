package com.kittycat;

public final class Fibonacci {
    private Fibonacci() {}

    public static long fibonacci(int n) {
        if (n < 0) throw new IllegalArgumentException("n must be non-negative");
        if (n < 2) return n;
        long a = 0, b = 1;
        for (int i = 2; i <= n; i++) {
            long next = a + b;
            a = b;
            b = next;
        }
        return b;
    }

    public static void main(String[] args) {
        int count = args.length > 0 ? Integer.parseInt(args[0]) : 10;
        for (int i = 0; i <= count; i++) {
            System.out.println("fib(" + i + ") = " + fibonacci(i));
        }
    }
}
