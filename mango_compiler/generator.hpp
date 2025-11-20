#include <stack>
#include <map>
#include <string>
#include <memory>

#include <llvm/IR/Module.h>
#include <llvm/IR/Function.h>
#include <llvm/IR/Type.h>
#include <llvm/IR/Value.h>
#include <llvm/IR/LLVMContext.h>
#include <llvm/IR/BasicBlock.h>
#include <llvm/IR/IRBuilder.h>
#include <llvm/ExecutionEngine/GenericValue.h>
#include <llvm/ExecutionEngine/ExecutionEngine.h>

class BlockNode;

// represent a single code block/scope in the program
class GenCodeBlock {
public:
    llvm::BasicBlock *bb;
    // map local variable names to llvm values
    std::map<std::string, llvm::Value*> locals;
};

// main orchestrator and scope manager
class GenCodeContext {
    // FIFO stack of code blocks
    std::stack<GenCodeBlock *> blocks;
    llvm::Function* main_func;
public:
    std::unique_ptr<llvm::Module> module;
    std::unique_ptr<llvm::LLVMContext> context;
    std::unique_ptr<llvm::IRBuilder<>> builder;

    void generate_code(BlockNode& root);
    llvm::GenericValue run_code();

    GenCodeContext() {
        context = std::make_unique<llvm::LLVMContext>();
        module = std::make_unique<llvm::Module>("mango", *context);
        builder = std::make_unique<llvm::IRBuilder<>>(*context);
        main_func = nullptr;
    }

    // return current scope and map of local variables
    std::map<std::string, llvm::Value*>& get_locals() {
        if (blocks.empty()) {
            throw std::runtime_error("no current scope");
        }
        return blocks.top()->locals;
    }

    // return current basic block
    llvm::BasicBlock* get_curr_bb() { return blocks.top()->bb; }

    // enter a new scope
    void push_block(llvm::BasicBlock *block) { 
        blocks.push(new GenCodeBlock());
        blocks.top()->bb = block;
        builder->SetInsertPoint(blocks.top()->bb);
    }
    // leave current scope
    void pop_block() {
        GenCodeBlock* block = blocks.top();
        blocks.pop();
        delete block;
    }
};
