#ifndef ASTNODE_HPP
#define ASTNODE_HPP

#include <vector>
#include <string>
#include <iostream>
#include <llvm/IR/Value.h>


// forward declarations
class GenCodeContext;
class ExprNode;
class StmtNode;
class VarDeclNode;

typedef std::vector<ExprNode*> ExprList;
typedef std::vector<StmtNode*> StmtList;
typedef std::vector<VarDeclNode*> VarDeclList;

// ast node base class
class Node {
public:
    virtual ~Node() = default;
};

class ExprNode : public Node {  // produces values
public:
    virtual llvm::Value* codeGen(GenCodeContext& context) = 0;
};

// performs action, does not necessarily produce a value
class StmtNode : public Node {
public:
    virtual llvm::Value* codeGen(GenCodeContext& context) = 0;
};

class IntegerNode : public ExprNode {
public:
    long long value;
    IntegerNode(long long val) : value(val) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

class DoubleNode : public ExprNode {
public:
    double value;
    DoubleNode(double val) : value(val) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

// e.g x
class IdentifierNode : public ExprNode { // names given to variables
public:
    std::string name;
    IdentifierNode(std::string n) : name(n) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

// e.g bar(1, 2);
class MethodCallNode : public ExprNode {
public:
    const IdentifierNode& name;
    ExprList args;
    MethodCallNode(const IdentifierNode& n, ExprList a) : name(n), args(a) {}
    // allow function calls without arguments
    MethodCallNode(const IdentifierNode& n) : name(n), args() {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

// e.g a + b;
class BinaryOpNode : public ExprNode {
public:
    ExprNode& lhs;
    ExprNode& rhs;
    int op;
    BinaryOpNode(ExprNode& l, ExprNode& r, int o) : lhs(l), rhs(r), op(o) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

// e.g x = 5;
class AssignmentNode : public ExprNode {
public:
    IdentifierNode& lhs;
    ExprNode& rhs;
    AssignmentNode(IdentifierNode& l, ExprNode& r) : lhs(l), rhs(r) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

class BlockNode : public StmtNode {
public:
    StmtList stmts;
    BlockNode() {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

// e.g int x; or int x = 5;
class VarDeclNode : public StmtNode {
public:
    const IdentifierNode& type;
    IdentifierNode& name;
    ExprNode* assignmentExpr;
    // allow declarations without assignment/initialization
    VarDeclNode(const IdentifierNode& t, IdentifierNode& n) : 
        type(t), name(n), assignmentExpr(nullptr) {}
    VarDeclNode(const IdentifierNode& t, IdentifierNode& n, ExprNode* a) : 
        type(t), name(n), assignmentExpr(a) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

// e.g int add(int a, int b) { return a + b; }
class FuncDeclNode : public StmtNode {
public:
    const IdentifierNode& return_type;
    const IdentifierNode& name;
    VarDeclList params;
    BlockNode& body;
    FuncDeclNode(const IdentifierNode& t, const IdentifierNode& n, VarDeclList p, BlockNode& b) : 
        return_type(t), name(n), params(p), body(b) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

class ExprStmtNode : public StmtNode {
public:
    ExprNode& expr;
    ExprStmtNode(ExprNode& e) : expr(e) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

#endif
