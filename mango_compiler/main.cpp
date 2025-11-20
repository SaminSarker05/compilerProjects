#include <iostream>
#include "astnode.hpp"

extern int yyparse();
extern BlockNode* program_block;

int main() {
    yyparse();  // waits for input
    // print out address representing root of AST
    std::cout << program_block << std::endl;
    return 0;
}
