CXX = g++
CXXFLAGS = -std=c++17 -Wall -Wextra -Werror -g -I.
PARSER = bison
PARSERFLAGS = -Wall -Werror -d
LEXER_SRC = lex.yy.c
LEXER = flex
LEXERFLAGS =

lex.yy.c: lexer.l
	$(LEXER) $(LEXERFLAGS) -o $(LEXER_SRC) lexer.l

.PHONY: clean
clean:
	rm -f *.o $(LEXER_SRC)
