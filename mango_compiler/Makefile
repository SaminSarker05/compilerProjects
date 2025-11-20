CXX = g++
CXXFLAGS = -std=c++17 -Wall -Wextra -Werror -g -I.
PARSER = bison
LEXER = flex

parser: parser.cpp lexer.cpp
	g++ -o parser parser.cpp lexer.cpp main.cpp

parser.cpp: parser.y
	bison -d -o parser.cpp parser.y

parser.hpp: parser.cpp

lexer.cpp: lexer.l parser.hpp
	flex -o lexer.cpp $^

.PHONY: clean
clean:
	rm -f *.o lexer.cpp parser.cpp parser.hpp parser
