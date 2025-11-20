#include "generator.hpp"
#include "astnode.hpp"
#include "parser.hpp"

// helper function to map language types to llvm types
static Type* type_of(llvm::LLVMContext& context, const IdentifierNode& ident) {
    // create llvm type object, use context to avoid duplication
    if (type.name == "int") {
        return Type::getInt64Ty(context);
    } else if (type.name == "double") {
        return Type::getDoubleTy(context);
    } else {
        return Type::getVoidTy(context);
    }
}

// entrypoint for code generation
void GenCodeContext::generate_code(BlockNode& root) {
    std::cout << "Generating IR code..." << std::endl;

    // create main function signature/type, with no args and void return
    std::vector<llvm::Type*> arg_types;
    llvm::FunctionType* ftype = FunctionType::get(Type::getVoidTy(*context), arg_types, false);
    // define actual main function, visible only within module
    main_func = llvm::Function::Create(ftype, llvm::GlobalValue::InternalLinkage, 
        "main", module.get());

    // add first basic block and start code generation
    llvm::BasicBlock* bb = llvm::BasicBlock::Create(*context, "entry", main_func);
    push_block(bb);
    root.codeGen(*this); // recursively generate code
    builder->CreateRetVoid();   // add return statement

    pop_block();  // leave main scope
    std::cout << "Code generation complete." << std::endl;
    module->print(llvm::outs(), nullptr);  // print human readable IR
}

// run generated code using JIT
llvm::GenericValue GenCodeContext::run_code() {
    std::cout << "Running code..." << std::endl;

    // prepare llvm for execution on target machine
    llvm::InitializeNativeTarget();
    llvm::InitializeNativeTargetAsmPrinter();

    // create execution engine for JIT
    llvm::ExecutionEngine* ee = llvm::EngineBuilder(std::move(module)).create();
    if (!ee) {
        std::cerr << "Engine creation failed!" << std::endl;
        return llvm::GenericValue();
    }

    std::vector<llvm::GenericValue> args;  // run main with no args
    llvm::GenericValue result = ee->runFunction(main_func, args);
   
    std::cout << "Code running complete." << std::endl;
    return result;
}

llvm::Value* IntegerNode::codeGen(GenCodeContext& context) {
    std::cout << "IntegerNode::codeGen" << std::endl;
    // create llvm 64 bit signed integer constant
    return llvm::ConstantInt::get(llvm::Type::getInt64Ty(*context.context), value, true);
}

llvm::Value* DoubleNode::codeGen(GenCodeContext& context) {
    std::cout << "DoubleNode::codeGen" << std::endl;
    // create llvm 64 bit double constant
    return llvm::ConstantFP::get(llvm::Type::getDoubleTy(*context.context), value);
}

llvm::Value* IdentifierNode::codeGen(GenCodeContext& context) {
    std::cout << "IdentifierNode::codeGen" << std::endl;
    // lookup variable in symbol table
    if (context.get_locals().find(name) == context.get_locals().end()) {
        std::cerr << "Undefined variable: " << name << std::endl;
        return nullptr;
    }

    // load variable from memory
    return context.builder->CreateLoad(llvm::Type::getInt64Ty(*context.context), 
        context.get_locals()[name], name.c_str());
}

llvm::Value* AssignmentNode::codeGen(GenCodeContext& context) {
    
}

llvm::Value* MethodCallNode::codeGen(GenCodeContext& context) {
    llvm::Function* func = context.module->getFunction(name);
    if (func == nullptr) {
        std::cerr << "Undefined function: " << name << std::endl;
        return nullptr;
    }
    // generate code for each expression in argument list
    std::vector<llvm::Value*> arg_values;
    for (ExprNode* arg : args) {
        args.push_back(arg->codeGen(context));
    }
    // create a call instruction
    return context.builder->CreateCall(func, arg_values, "calldebug");
}

llvm::Value* BinaryOpNode::codeGen(GenCodeContext& context) {
    
}

llvm::Value* BlockNode::codeGen(GenCodeContext& context) {
    
}

llvm::Value* ExprStmtNode::codeGen(GenCodeContext& context) {
    
}

llvm::Value* FuncDeclNode::codeGen(GenCodeContext& context) {
    
}

llvm::Value* VarDeclNode::codeGen(GenCodeContext& context) {
    
}
