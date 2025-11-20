#ifndef ASTNODE_HPP
#define ASTNODE_HPP

#include <vector>
#include <string>
#include <iostream>
#include <llvm/IR/Value.h>


// forward declarations
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
    virtual llvm::Value* codeGen(GenCodeContext& context) { }
};

class ExprNode : public Node {  // produces values
};

// performs action, does not necessarily produce a value
class StmtNode : public Node {
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

class IdentifierNode : public ExprNode { // names given to variables
public:
    std::string name;
    IdentifierNode(std::string n) : name(n) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

class MethodCallNode : public ExprNode {
public:
    const IdentifierNode& name;
    ExprList args;
    MethodCallNode(const IdentifierNode& n, ExprList a) : name(n), args(a) {}
    // allow function calls without arguments
    MethodCallNode(const IdentifierNode& n) : name(n), args() {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

class BinaryOpNode : public ExprNode {
public:
    ExprNode& lhs;
    ExprNode& rhs;
    int op;
    BinaryOpNode(ExprNode& l, ExprNode& r, int o) : lhs(l), rhs(r), op(o) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

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

class FuncDeclNode : public StmtNode {
public:
    const IdentifierNode& type;
    const IdentifierNode& name;
    VarDeclList params;
    BlockNode& body;
    FuncDeclNode(const IdentifierNode& t, const IdentifierNode& n, VarDeclList p, BlockNode& b) : 
        type(t), name(n), params(p), body(b) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

class ExprStmtNode : public StmtNode {
public:
    ExprNode& expr;
    ExprStmtNode(ExprNode& e) : expr(e) {}
    virtual llvm::Value* codeGen(GenCodeContext& context);
};

#endif
