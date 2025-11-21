#include <iostream>
#include "astnode.hpp"
#include "generator.hpp"

extern int yyparse();
extern BlockNode* program_block;

int main() {
    yyparse();  // waits for input
    // print out address representing root of AST
    std::cout << program_block << std::endl;

    // walk over AST to generate IR and run code
    GenCodeContext context;
    context.generate_code(*program_block);
    context.run_code();

    return 0;
}
