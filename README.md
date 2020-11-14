# Bast
Compiler for BAST(Basic Tiny) - A Tiny Toy Basic Dialect

# Introduction
This compiler takes code written in bast and compiles(transpiles) to C code.You need to compile the C code with C compiler (gcc/tcc/clang are tested to work) to produce executable.

This provided the opportunity to learn about lexing,parsing and to learn more Golang. I believe that this toy language will not be used other than for educational purpose(for people like past-me).

# Installation 
You can install `bast` compiler by typing in the terminal 

```bash
go get "github.com/varbhat/bast"
```

# Usage

* `bast --help` will print the help.

* `bast -in=filename.bast -out=filename.c` will compile bast source `filename.bast` to C source `filename.c`. 

* You need to compile the emitted C file with C compiler like GCC/Clang/TCC.

* `cc ./filename.c -o filename` (where cc is gcc/clang/tcc)

* You can combine these steps :

```bash
bast -in=filename.bast -out=filename.c && cc filename.c -o filename
```

And then run the executable binary `filename` it produced.

# Language 
Language is small Dialect of BASIC same what AZHenly implemented.I only wrote compiler.
Grammer file can be found at [grammer.txt](https://github.com/varbhat/bast/blob/master/grammar.txt)

It supports:
  - Numerical variables
  - Basic arithmetic
  - If statements
  - While loops
  - Print text and numbers
  - Input numbers
  - Labels and goto
  - Comments

# Example code

```
PRINT "How many fibonacci numbers do you want?"
INPUT nums
PRINT ""

LET a = 0
LET b = 1
WHILE nums > 0 REPEAT
    PRINT a
    LET c = a + b
    LET a = b
    LET b = c
    LET nums = nums - 1
ENDWHILE
```
You can find more examples at [Examples](https://github.com/varbhat/bast/tree/master/examples)

# Thanks
Thanks to [AZHenly](https://github.com/AZHenley) for their post  [Let's make a Teeny Tiny compiler](http://web.eecs.utk.edu/~azh/blog/teenytinycompiler1.html) which acted as reference. 
But , instructions in that post are for Python and the compiler is implemented in Python. I who wanted to know more about compiler , read it as reference , wrote the compiler from scratch in Golang.

Although the language for both the compilers are same,`bast` compiler is written from scratch in Golang ,implemented by instructions of their post.

# License
This Software has been licensed under [MIT](https://github.com/varbhat/bast/blob/master/LICENSE)
