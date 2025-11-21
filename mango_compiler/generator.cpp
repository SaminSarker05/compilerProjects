#include "generator.hpp"
#include "astnode.hpp"
#include "parser.hpp"
#include <llvm/Support/TargetSelect.h>

// helper function to map language types to llvm types
static llvm::Type* type_of(llvm::LLVMContext& context, const IdentifierNode& ident) {
    // create llvm type object, use context to avoid duplication
    if (ident.name == "int") {
        return llvm::Type::getInt64Ty(context);
    } else if (ident.name == "double") {
        return llvm::Type::getDoubleTy(context);
    } else {
        return llvm::Type::getVoidTy(context);
    }
}

// entrypoint for code generation
void GenCodeContext::generate_code(BlockNode& root) {
    std::cout << "Generating IR code..." << std::endl;

    // create main function signature/type, with no args and void return
    std::vector<llvm::Type*> arg_types;
    llvm::FunctionType* ftype = llvm::FunctionType::get(llvm::Type::getVoidTy(*context), arg_types, false);
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
    std::cout << "AssignmentNode::codegen" << std::endl;
    // ensure lhs variable exists
    if (context.get_locals().find(lhs.name) == context.get_locals().end()) {
        std::cerr << "Undefined variable: " << lhs.name << std::endl;
        return nullptr;
    }
    llvm::Value* rhs_value = rhs.codeGen(context);
    if (!rhs_value) { return nullptr; }

    // store rhs value into memory location of lhs variable
    return context.builder->CreateStore(rhs_value, context.get_locals()[lhs.name]);
}

llvm::Value* MethodCallNode::codeGen(GenCodeContext& context) {
    std::cout << "MethodCallNode::codeGen" << std::endl;
    llvm::Function* func = context.module->getFunction(name.name.c_str());
    if (func == nullptr) {
        std::cerr << "Undefined function: " << name.name << std::endl;
        return nullptr;
    }
    // generate code for each expression in argument list
    std::vector<llvm::Value*> arg_values;
    for (ExprNode* arg : args) {
        arg_values.push_back(arg->codeGen(context));
    }
    // create a call instruction
    return context.builder->CreateCall(func, arg_values, "calltmp");
}

llvm::Value* BinaryOpNode::codeGen(GenCodeContext& context) {
    std::cout << "BinaryOpNode::codeGen" << std::endl;
    llvm::Value* l = lhs.codeGen(context);
    llvm::Value* r = rhs.codeGen(context);

    if (!l || !r) { return nullptr; }

    // create binary operation
    switch(op) {
    case PLUS:
        return context.builder->CreateAdd(l, r, "addtmp");
    case MINUS:
        return context.builder->CreateSub(l, r, "subtmp");
    case MUL:
        return context.builder->CreateMul(l, r, "multmp");
    case DIV:
        return context.builder->CreateSDiv(l, r, "divtmp");
    case LT:
        return context.builder->CreateICmpSLT(l, r, "lttmp");
    case GT:
        return context.builder->CreateICmpSGT(l, r, "gttmp");
    case EQUAL: // == check
        return context.builder->CreateICmpEQ(l, r, "eqtmp");
    default:
        std::cerr << "Unsupported binary operator: " << op << std::endl;
        return nullptr;
    }
}

llvm::Value* BlockNode::codeGen(GenCodeContext& context) {
    std::cout << "BlockNode::codeGen" << std::endl;
    llvm::Value* last = nullptr;
    // generate code for each statement in block
    for (StmtNode* stmt : stmts) {
        last = stmt->codeGen(context);
    }
    return last;
}

llvm::Value* ExprStmtNode::codeGen(GenCodeContext& context) {
    std::cout << "ExprStmtNode::codeGen" << std::endl;
    return expr.codeGen(context);
}

llvm::Value* FuncDeclNode::codeGen(GenCodeContext& context) {
    std::cout << "FuncDeclNode::codeGen" << std::endl;
    // build list of parameter types
    std::vector<llvm::Type*> arg_types;
    for (VarDeclNode* arg : params) {
        arg_types.push_back(type_of(*context.context, arg->type));
    }

    // create function type object and actual function
    llvm::FunctionType* ftype = llvm::FunctionType::get(
        type_of(*context.context, return_type), arg_types, false
    );
    llvm::Function* func = llvm::Function::Create(ftype, llvm::Function::InternalLinkage, 
        name.name.c_str(), context.module.get());
    
    // create basic block for function
    llvm::BasicBlock* bb = llvm::BasicBlock::Create(*context.context, "entry", func);
    context.push_block(bb);  // enter function scope

    // generate code for each parameter, function body, and return instruction
    for (VarDeclNode* arg : params) {
        arg->codeGen(context);
    }
    body.codeGen(context);
    context.builder->CreateRetVoid();

    // exit function scope
    context.pop_block();
    return func;
    
}

llvm::Value* VarDeclNode::codeGen(GenCodeContext& context) {
    std::cout << "VarDeclNode::codeGen" << std::endl;

    // get llvm type of variable
    llvm::Type* var_type = type_of(*context.context, type);
    if (var_type == nullptr) {
        std::cerr << "Unsupported variable type: " << type.name << std::endl;
        return nullptr;
    }

    // allocate variable on stack
    llvm::AllocaInst* alloca = context.builder->CreateAlloca(var_type, nullptr, name.name.c_str());
    
    // add variable to local variable scope
    context.get_locals()[name.name] = alloca;
    // if assignment expression, generate code for rhs
    if (assignmentExpr) {
        AssignmentNode assign(name, *assignmentExpr);
        assign.codeGen(context);
    }
    return alloca;
}
