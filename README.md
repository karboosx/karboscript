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

## Usage

Execute `script.ks` file
```
# ./karboscript script.ks
```

Show opcodes for `script.ks` file
```
# ./karboscript --opcode script.ks
```

## Buildin functions

We have to our disposal couple of buildin functions:

| function name | arguments | return | example |
|---------------|-----------|--------|---------|
| out() | any variable... | nothing | out(1,2,3); |
| readLine() | nothing | string | name = readLine(); |
| readInt() | nothing | int | name = readInt(); |

## Syntax

# Declare function

```c
function <name>([<type> <argument_name>], ...) [return_type] {
    [body]
}
```

# Declare variable
```c
<type> <var_name> = <expression>;
```
For example: `string test = "hello world";`

# Array
```c
array <var_name> = [<expression>, ...];
```

Access array element
```c
<array_name>[<index>]
```

Assign array element
```c
<array_name>[<index>] = <expression>;
```

Push new element to array
```c
<array_name>[] = <expression>;
```

# Call function
```c
<function_name>(<argument>, ...);
```
For example: `func(1, 2, 3, variable);`

# Loops
While
```c
    while (<expresion>) {
        [body]
    }
```

For
```c
    for <init_statement>; <compare_expresion>; <inrement_statement>; {
        [body]
    }
```

From to
```c
    from <starting_value_expresion> to <ending_value_expresion> as <variable_name>; {
        [body]
    }
```

## Examples

Fibonaci:
```c
function main() {
    int a = 1;
    int b = 1;

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
    int a = 1;
    int b = 1;

    while (a < 5) {
        b = 1;
        while (b < 5) {
            out ("test", b);
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

function test(int test)
{
    return test + 200;
}
```

For loop
```c
function main() {
    for int i=0; i<10; i=i+1; {
        out(i);
    }
}
```

For increment (loop from one expresion to another with interval of 1)
```c
function main() {
    from 0 to 10 as i {
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

Return type
```c
function main() {
    out(test());
}

function test() string {
    return "test";
}
```

Array declaration
```c
function main() {
    array a = [1, 2, 3, 4, 5];
    out(a[0]);
}
```

Array assignment
```c
function main() {
    array a = [1, 2, 3, 4, 5];
    a[0] = 10;
    out(a[0]);
}
```

# Intermediary code
The code is compiled to intermediary code which is then executed by virtual machine. For example this code:
```c
function main() {
    int a = 1;
    int b = 1;
    int max = 500;
    while (b < max) {
        out (b);

        int c = b;
        int b = a + b;
        a = c;
    }
}
```

Will be compiled to:
```
"_function.main: function"
"add_scope"
"push_exp ( 1 )"
"sub_scope"
"set_local_var_exp ( int a )"
"add_scope"
"push_exp ( 1 )"
"sub_scope"
"set_local_var_exp ( int b )"
"add_scope"
"push_exp ( 500 )"
"sub_scope"
"set_local_var_exp ( int max )"
"_while.d: while_start"
"add_scope"
"push_exp_var ( b )"
"push_exp_var ( max )"
"exp_call ( < )"
"sub_scope"
"while ( _while.13 )"
"add_scope"
"push_exp_var ( b )"
"sub_scope"
"push_function_arg"
"call_function ( out 1 )"
"add_scope"
"push_exp_var ( b )"
"sub_scope"
"set_local_var_exp ( int c )"
"add_scope"
"push_exp_var ( a )"
"push_exp_var ( b )"
"exp_call ( + )"
"sub_scope"
"set_local_var_exp ( int b )"
"add_scope"
"push_exp_var ( c )"
"sub_scope"
"set_local_var_exp (  a )"
"jmp ( _while.d )"
"_while.13: while_else"
"function_return"
"call_function ( main 0 )"
"exit ( main 0 )"
```