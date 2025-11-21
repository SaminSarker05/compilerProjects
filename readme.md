```
## Project Descriptions

### mango_compiler/
C-like compiler built with C++, Flex, Bison, and LLVM that compiles to LLVM IR and executes with JIT.

Tech: C++17, LLVM 18, Flex, Bison, Docker
Features: Lexical analysis, parsing, AST construction, LLVM IR generation, JIT execution
Supports: Variables, arithmetic operations, comparisons, function declarations

Build:
  docker build -t mango .
  docker run --rm -it -v $(pwd):/usr/src/app mango
  make
  echo "int x = 5; int y = x + 10;" | ./mango

### gorilla_interpreter/
A tree-walking interpreter written in Go, inspired by "Writing An Interpreter in Go". [IN WORK]

Tech: Go 1.20
Current Features: Lexer with lookahead, token classification, REPL interface
Supports: Variables, functions, conditionals, arithmetic, comparison operators

Run:
  go run main.go
  👾 ~ let x = 5;
  👾 ~ let result = add(x, 10);

### lisp_c_compiler/
Really small compiler that transforms Lisp-style S-expressions to C-style function calls.

Tech: Go
Pipeline: Lexer → Parser → AST → Transformer (visitor pattern) → Code Generator
Example: (multiply (add 1 2) (divide 6 3)) → multiply(add(1, 2), divide(6, 3));

Run:
  go run main.go

### minilang_lexer/
Lexical scanner for a simple language with maximal munch tokenization.

Tech: C++17, Regex
Features: Lookahead mechanism, comment handling, error reporting with line/column
Supports: Keywords, identifiers, integers, operators, delimiters

Build:
  make
  ./lexer test/input/sample1.minilang

## Note

Each project is self-contained.
```
