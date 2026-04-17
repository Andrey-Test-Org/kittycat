fun fibonacci(n: Int): Long {
    require(n >= 0) { "n must be non-negative" }
    if (n < 2) return n.toLong()
    var a = 0L
    var b = 1L
    for (i in 2..n) {
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
