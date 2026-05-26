# kittycat

Multi-language sandbox for small algorithms and demos.

## Project structure

```
kittycat/
├── src/          # Kotlin: Main.kt (Fibonacci, factorial, gcd, isPrime, sumToN)
├── python/       # Python scripts (fibonacci, factorial, primes)
├── java/         # Java Maven module (com.kittycat: HelloWorld, Fibonacci, MathUtils)
├── web/          # HTML + CSS landing site (index, about, dark-mode toggle)
└── prompts/      # Misc prompt files
```

## Run

### Kotlin

```
kotlinc src/Main.kt -include-runtime -d kittycat.jar
java -jar kittycat.jar 10
```

First arg = how many Fibonacci numbers to print (default 10).

### Python

```
python3 python/fibonacci.py 10
python3 python/factorial.py 6
python3 python/primes.py 20
```

### Java

```
cd java
mvn -q compile
mvn -q exec:java -Dexec.mainClass=com.kittycat.Fibonacci -Dexec.args=10
```

### Web

Open `web/index.html` directly in a browser, or serve `web/` with any static file server.
