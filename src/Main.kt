fun sumToN(n: Int): Long {
    require(n >= 0) { "n must be non-negative" }
    return n.toLong() * (n + 1) / 2
}

fun isPrime(n: Long): Boolean {
    if (n < 2) return false
    if (n < 4) return true
    if (n % 2 == 0L) return false
    var i = 3L
    while (i * i <= n) {
        if (n % i == 0L) return false
        i += 2
    }
    return true
}

fun fibonacci(n: Int): Long {
    require(n >= 0) { "n must be non-negative" }
    if (n < 2) return n.toLong()
    var a = 0L
    var b = 1L
    (2..n).forEach { i ->
        val next = a + b
        a = b
        b = next
    }
    return b
}

fun main(args: Array<String>) {
    val count = args.firstOrNull()?.toIntOrNull() ?: 10
    for (i in 0..count) {
        println("fib($i) = ${fibonacci(i)}")
    }
}
