
#ifndef LEXER_H
/* A lexical scanner */
#define LEXER_H

#include <vector>
#include <fstream>
#include <iostream>

class LexerException : public std::runtime_error {
public:
    LexerException(int line, int col) : std::runtime_error("Lexer Error at line " + std::to_string(line) +
                                  ", column " + std::to_string(col)) {}
};

class Token {
public:
    enum Type {
        Identifier, Keyword, Integer, AssignOp, Op, Delimiter, EndOfFile
    };

    Type type;
    std::string text;
    int line;
    int column;

    Token(Type t=EndOfFile, const std::string& txt="", int l=0, int c=0)
        : type(t), text(txt), line(l), column(c) {}

    std::string typeToString() const {
        switch (type) {
            case Identifier: return "Identifier";
            case Keyword: return "Keyword";
            case Integer: return "Integer";
            case AssignOp: return "AssignOp";
            case Op: return "Operator";
            case Delimiter: return "Delimiter";
            case EndOfFile: return "EOF";
            default: return "Unknown";
        }
    }

    void print() const {
        std::cout << typeToString() << "('" << text << "') at " << line << ":" << column << std::endl;
    }
};

class Lexer {
    std::ifstream in;
    std::vector <Token> tokens;
    int line = 1, col = 0;
public:
    Lexer(const char *filename);
    ~Lexer();
    void printTokens();
    Token getNextToken();
};

#endif /* LEXER_H */