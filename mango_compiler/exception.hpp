#ifndef EXCEPTION_HPP
#define EXCEPTION_HPP

#include <stdexcept>
#include <string>

class LexerException : public std::runtime_error {
public:
    LexerException() : std::runtime_error("Lexer Error") {}
};

#endif
