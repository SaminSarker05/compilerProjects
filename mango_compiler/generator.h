#include <stack>
#include <llvm/Module.h>
#include <llvm/Function.h>
#include <llvm/Type.h>
#include <llvm/DerivedTypes.h>
#include <llvm/LLVMContext.h>
#include <llvm/PassManager.h>
#include <llvm/Instructions.h>
#include <llvm/CallingConv.h>
#include <llvm/Bitcode/ReaderWriter.h>
#include <llvm/Analysis/Verifier.h>
#include <llvm/Assembly/PrintModulePass.h>
#include <llvm/Support/IRBuilder.h>
#include <llvm/ModuleProvider.h>
#include <llvm/Target/TargetSelect.h>
#include <llvm/ExecutionEngine/GenericValue.h>
#include <llvm/ExecutionEngine/JIT.h>
#include <llvm/Support/raw_ostream.h>

// represent a single code block/scope in the program
class GenCodeBlock {
public:
    llvm::BasicBlock *bb;
    // map local variable names to llvm values
    std::map<std::string, llvm::Value*> symbol_table;
};

// state/scope manager and main orchestration
class GenCodeContext {
public:
    llvm::Module* module;
    llvm::BasicBlock* curr_bb() { return blocks.top()->bb; }
    void push_block(GenCodeBlock *block) { 
        blocks.push(new GenCodeBlock());
        blocks.top()->bb = block;
        builder.SetInsertPoint(blocks.top()->bb);
    }
    void pop_block() {
        GenCodeBlock* block = blocks.top();
        blocks.pop();
        delete block;
    }
private:
    // FIFO stack of code blocks
    std::stack<GenCodeBlock *> blocks;
    llvm::Function* main_func;
    llvm::IRBuilder<> builder;
};
