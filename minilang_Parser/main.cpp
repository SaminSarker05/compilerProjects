#include <iostream>
#include "lexer.h"

int main(int argc, char** argv) {
    if (argc < 2) {
        std::cerr << "Usage: " << argv[0] << " <source-file>" << std::endl;
        return 1;
    }
    Lexer lexer(argv[1]);
    while(true) {
        try
        {
            Token t = lexer.getNextToken();
            t.print();
            if (t.type == Token::EndOfFile) return 0;
        }
        catch(const LexerException& e)
        {
            std::cerr << e.what() << std::endl;
            return 0;
        }
    }
    return 0;
}