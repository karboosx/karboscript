[![Tests](https://github.com/karboosx/karboscript/actions/workflows/tests.yml/badge.svg?branch=master)](https://github.com/karboosx/karboscript/actions/workflows/tests.yml)

# KarboScript

Simple programming language made for fun to learn how parsers and compilers works. Don't have much capability but have the basic stuff like functions, if, loops.

## Main program

Each program has to have `main()` function:
```
function main() {
    //add your code here
}
```

This is the starting point of every script.


## Buildin functions

We have to our disposal couple of buildin functions:

| function name | arguments | return | example |
|---------------|-----------|--------|---------|
| out() | any variable... | nothing | out(1,2,3); |
| readline() | nothing | string | name = readline(); |

## Examples

Fibonaci:
```c
function main() {
    a = 1;
    b = 1;

    while (b < 500) {

        c = b;
        b = a + b;
        a = c;
        out (b);
    }
}
```

while loop:
```c
function main() {
    a = 1;
    b = 1;

    while (a < 5) {
        b = 1;
        while (b < 5) {
            test ("test", b);
            b=b+1;
        }

        a=a+1;
    }
}

```

If statement
```c
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

Arguments for function
```c
function main()
{
    out(1000 + test(800), test(500));
}

function test(test)
{
    return test + 200;
}
```

For loop
```c
function main() {
    for i=0; i<10; i=i+1; {
        out(i);
    }
}
```

Read line from stdin
```c
function main() {
    out("Enter name: ");
    name = readLine();
    out("Your name is:", name);
}
```
