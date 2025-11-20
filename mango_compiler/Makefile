CXX = g++
CXXFLAGS = -std=c++17 -Wall -Wextra -Werror -g -I.
PARSER = bison
PARSERFLAGS = -Wall -Werror -d
LEXER = flex
LEXERFLAGS =

lex.yy.c: lexer.l
	$(LEXER) $(LEXERFLAGS) -o lex.yy.c lexer.l

.PHONY: clean
clean:
	rm -f *.o lex.yy.c parser.cpp parser.hpp
