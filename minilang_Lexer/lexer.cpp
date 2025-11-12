/*
A lexical scanner for a simple programming language "minilang".

   The lexical scanner reads characters one byte at a time from a file stream
and recognizes lexemes using regex. A boolean flag is used to track
lines that have comments and when those commented sections of code begin. If
the current column is within a comment or we encounter a blank space, \r, or
\t we skip the character.

    A lookahead mechanism using the std::istream::peak is used to implement the
maximal munch policy and determine when a token needs to be classified. The
lookahead finalizes the current word if the next char is invalid, a delimiter,
EOF, etc. The actual token classification is done using regex expressions,
combining union and kleene star/plus operators for matches.

    If a token cannot be classified, a LexerException is thrown with the
current line and column number.

to run:
    make
    ./lexer <source-file>

example:
    ./lexer test/input/sample1.minilang
*/

#include "lexer.h"
#include <cctype>
#include <sstream>
#include <unordered_set>

#include <string>
#include <regex>

Lexer::Lexer(const char *filename) : in(filename) {
    if (!in) {
        throw std::runtime_error("Failed to open file");
    }
}

Lexer::~Lexer() {
    in.close();
}

Token Lexer::getNextToken() {
    bool comment = false;  // flag if current line is comment 
    char next_char;
    char c;
    int word_start = 0;
    std::string word = "";

    // regular expressions for valid language elements
    std::regex identifiers("[a-zA-Z_][a-zA-Z0-9_]*");
    std::regex keywords("func|var|let|if|else|while|print|return");
    std::regex integers("[0-9]+");
    std::regex operators("==|!=|<=|>=|<|>|/|\\+|\\-|\\*");
    std::regex delimiters("[(){};:,]");

    // set to identify invalid chars such as '^' and ']'
    std::regex invalidChars("[^a-zA-Z0-9_(){};:,=!<>/\\+\\-\\*]");
    
    while (in.get(c)) {
        this->col += 1;
        if (c == '\n') {  // increase line # when newline char encountered
            this->line += 1;
            this->col = 0;
            comment = false;
            continue;
        }
        if (comment || c == ' ' || c == '\r' || c == '\t') { continue; }
        if (c == '#') {
            comment = true;
            continue;
        }

        // if unrecognized character throw immediate error
        if (std::regex_match(std::string(1, c), invalidChars)) {
            throw LexerException(this->line, this->col);
        }

        // if delimiter matched return token
        if (std::regex_match(std::string(1, c), delimiters)) {
            if (c == ':' && in.peek() == '=') {
                in.get();  // consume next char 
                this->col += 1;
                return Token(Token::AssignOp, ":=", this->line, this->col - 1);
                // disambiguate between ':' delimiter and ':=' assignment op
            } else {
                word = c;
                return Token(Token::Delimiter, word, this->line, this->col);
            }
        }

        // mark the start of a token
        if (word.length() == 0) {
            word_start = this->col;
        }

        /*  add character to running word then check 
            if next char marks end of byte consumption,
            if so try to classify word using regex
        */
        word += c;
        next_char = in.peek();  // lookahead to stop eating bytes
        if (next_char == ' ' || next_char == '\n' || next_char == '\r'
            || next_char == '\t' || next_char == EOF
            || std::regex_match(std::string(1, next_char), delimiters)
            || std::regex_match(std::string(1, next_char), invalidChars)) {
            // try to classify token and return; if we cannot return error
            Token token;
            if (std::regex_match(word, keywords)) {
                token = Token(Token::Keyword, word, this->line, word_start);
                this->tokens.push_back(token);
            } else if (std::regex_match(word, identifiers)) {
                token = Token(Token::Identifier, word, this->line, word_start);
            } else if (std::regex_match(word, integers)) {
                token = Token(Token::Integer, word, this->line, word_start);
            } else if (std::regex_match(word, operators)) {
                token = Token(Token::Op, word, this->line, word_start);
            } else { // if token type never set then not valid
                throw LexerException(this->line, this->col);
            }

            // store lexeme in token vector of Lexer class
            tokens.push_back(token);
            return token;
        }
    }

    // final end of file token to mark end of input
    return Token(Token::EndOfFile, "", this->line, this->col);
}

void Lexer::printTokens() {
    for (const auto& token : tokens) {
        token.print();
    }
}
