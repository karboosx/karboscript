[![Tests](https://github.com/karboosx/karboscript/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/karboosx/karboscript/actions/workflows/tests.yml)

# KarboScript

Simple programming language made for fun to learn how parsers and compilers works. Don't have much capability but have the basic stuff like functions, if, loops.

## Examples
Fibonaci:
```
function main() {
    $a = 1;
    $b = 1;

    while ($b < 500) {

        $c = $b;
        $b = $a + $b;
        $a = $c;
        out ($b);
    }
}
```

Using while:
```
function main() {
    $a = 1;
    $b = 1;

    while ($a < 5) {
        $b = 1;
        while ($b < 5) {
            test ("test", $b);
            $b=$b+1;
        }

        $a=$a+1;
    }
}

```

Ifs
```
function main()
{
    if (10 == 10) {
        out("10 == 10");
    }
    if (500 < 200) {
        out("500 < 200");
    }
    if (12 > 10) {
        out("12 > 10");
    }
}
```

Arguments for functions
```
function main()
{
    out(1000 + test(800), test(500));
}

function test($test)
{
    return $test + 200;
}
```

For loop
```
function main() {
    for $i=0; $i<10; $i=$i+1; {
        out($i);
    }
}
```

Read line from stdin
```
function main() {
    out("Enter name: ");
    $name = readLine();
    out("Your name is:", $name);
}
```