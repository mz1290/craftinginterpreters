#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "common.h"
#include "compiler.h"
#include "memory.h"
#include "scanner.h"

//#ifdef DEBUG_PRINT_CODE
#include "debug.h"
//#endif


// Storage for current and previous tokens
typedef struct {
    Token current;
    Token previous;
    bool  hadError;
    bool  panicMode;
} Parser;

// clox precedence lovels from lowest to highest
typedef enum {
    PREC_NONE,
    PREC_ASSIGNMENT,  // =
    PREC_OR,          // or
    PREC_AND,         // and
    PREC_EQUALITY,    // == !=
    PREC_COMPARISON,  // < > <= >=
    PREC_TERM,        // + -
    PREC_FACTOR,      // * /
    PREC_UNARY,       // ! -
    PREC_CALL,        // . ()
    PREC_PRIMARY
} Precedence;

// Simple typedef for a function type that takes no arguments and returns
// nothing
typedef void (*ParseFn)(bool);

typedef struct {
    ParseFn prefix;
    ParseFn infix;
    Precedence precedence;
} ParseRule;

typedef struct {
    Token name; // Variable name
    int depth;  // Scope depth
} Local;

typedef struct {
    Local locals[UINT8_COUNT];
    int localCount;
    // The number of blocks surrounding the current bit of code being compiled.
    // 0 = global, 1 = first top level block, 2 = withhin that, 3.. etc
    // This allows us to keep track of which block each local var belongs to so
    // that we know which locals to discard when the block ends.
    int scopeDepth;
} Compiler;

// Pattern used throughout clox. Using a single global variable allows us to
// pass the state around from function to function in Compiler.
Parser parser;

// Not the best implementation, but since we are not doing any concurrent
// features it gets the job dones as a Global. A better appraoch would be to
// have each function receive a pointer to it's compiler.
Compiler* current = NULL;

Chunk* compilingChunk;

static Chunk* currentChunk() {
    return compilingChunk;
}

static void errorAt(Token* token, const char* message) {
    // Check if we are already in panic. We suppress any errors detected after
    // inital panic.
    if (parser.panicMode) return;

    // Signal panic mode
    parser.panicMode = true;

    // Print where error occurred
    fprintf(stderr, "[line %d] error", token->line);

    // Print lexeme if human readable
    if (token->type == TOKEN_EOF) {
        fprintf(stderr, " at end");
    } else if (token->type == TOKEN_ERROR) {
        // Nothing.
    } else {
        fprintf(stderr, " at \"%.*s\"", token->length, token->start);
    }

    // Print error message
    fprintf(stderr, ": %s\n", message);

    // Signal error occurred during compilation
    parser.hadError = true;
}

// Report error to user at previous token location
static void error(const char* message) {
    errorAt(&parser.previous, message);
}

// Report error to user at current token location
static void errorAtCurrent(const char* message) {
    errorAt(&parser.current, message);
}

// Read the next token. Continue reading tokens and reporting errors until a
// non-error token is read OR EOF. This ensures the Compiler only sees valid
// tokens.
static void advance() {
    // Storing here allows us to access lexeme after token match
    parser.previous = parser.current;

    for (;;) {
        parser.current = scanToken();
        if (parser.current.type != TOKEN_ERROR) break;

        errorAtCurrent(parser.current.start);
    }
}

// Read the next token and validate that the token matches expected type. If
// not, report error.
static void consume(TokenType type, const char* message) {
    if (parser.current.type == type) {
        advance();
        return;
    }

    errorAtCurrent(message);
}

static bool check(TokenType type) {
    return parser.current.type == type;
}

static bool match(TokenType type) {
    if (!check(type)) return false;
    advance();
    return true;
}

// Writes the given byte (opcode or operand from instruction) to the chunk
static void emitByte(uint8_t byte) {
    // Send previous token's line so runtime errors are associated with that
    // line.
    writeChunk(currentChunk(), byte, parser.previous.line);
}

// Convenience function for writing opcode followed by one-byte operand
static void emitBytes(uint8_t byte1, uint8_t byte2) {
    emitByte(byte1);
    emitByte(byte2);
}

static void emitReturn() {
    emitByte(OP_RETURN);
}

static uint8_t makeConstant(Value value) {
    int constant = addConstant(currentChunk(), value);

    if (constant > UINT8_MAX) {
        error("too many constants in one chunk");
        return 0;
    }

    return (uint8_t)constant;
}

static void emitConstant(Value value) {
    // Add value to constant table using makeConstant(), then emit OP_CONSTANT
    // instruction and push onto stack at runtime.
    emitBytes(OP_CONSTANT, makeConstant(value));
}

static void initCompiler(Compiler* compiler) {
    compiler->localCount = 0;
    compiler->scopeDepth = 0;
    current = compiler;
}

static void endCompiler() {
    emitReturn();

    if (DEBUG_LOX & DF_CODE) {
        if (!parser.hadError) {
            disassembleChunk(currentChunk(), "code");
        }
    }
//#ifdef DEBUG_PRINT_CODE
//    if (!parser.hadError) {
//        disassembleChunk(currentChunk(), "code");
//    }
//#endif
}

static void beginScope() {
    current->scopeDepth++;
}

static void endScope() {
    current->scopeDepth--;

    while (current->localCount > 0 &&
        current->locals[current->localCount - 1].depth > current->scopeDepth) {
        emitByte(OP_POP); // Tell the VM to remove values from stack at runtime
        current->localCount--;
    }
}

// These forward declaration enable clox's grammar recursive nature. For example
// binary() is defined *before* the rules table so the table can store a pointer
// to it. This prevents binary() function body from being able to access the
// table directly. To allow binary() to access the table, we define function
// getRule() *after* the table and forward declare here so binary() can access
// the table when needed.
static void expression();
static void statement();
static void declaration();
static ParseRule* getRule(TokenType type);
static void parsePrecedence(Precedence precedence);

static uint8_t identifierConstant(Token* name) {
    return makeConstant(OBJ_VAL(copyString(name->start, name->length)));
}

// Check if two identifiers are equal
static bool identifiersEqual(Token* a, Token* b) {
    if (a->length != b->length) return false;
    return memcmp(a->start, b->start, a->length) == 0;
}

static int resolveLocal(Compiler* compiler, Token* name) {
    // We always walk backwards because it ensures we enforce the *shadowing*
    // feature of Lox. Multiple variables of the same name can exist in the
    // surrounding scopes, but we only care about the last one.
    for (int i = compiler->localCount - 1; i >= 0; i--) {
        Local* local = &compiler->locals[i];

        if (identifiersEqual(name, &local->name)) {
            // Verify that the variable has been fully defined
            if (local->depth == -1) {
                error("can't read local variable in its own initializer");
            }

            return i;
        }
    }

    // Local not found, must be global
    return -1;
}

static void addLocal(Token name) {
    // VM can only support up to 256 local variables in a given scope
    if (current->localCount == UINT8_COUNT) {
        error("too many local variables in function");
        return;
    }

    // Initialize the next available Local in the compiler’s array of variables
    Local* local = &current->locals[current->localCount++];

    local->name = name;

    // Signals that the variable is *currently* uninitialized. Once the variable
    // initializer gets compiled, this value will change.
    local->depth = -1;
}

// This is how the compiler recognizes the existence of local variables. If
// global, then we return immediately.
static void declareVariable() {
    if (current->scopeDepth == 0) return;

    Token* name = &parser.previous;

    // Check for invalid redclarations of same variable name
    for (int i = current->localCount - 1; i >= 0; i--) {
        Local* local = &current->locals[i];

        // If we reach the beginning of the locals array OR encounter a variable
        // owned by a different scope, we know we've check all existing
        // variables in the current scope, exit.
        if (local->depth != -1 && local->depth < current->scopeDepth) {
            break; 
        }

        if (identifiersEqual(name, &local->name)) {
            error("already a variable with this name in this scope");
        }
    }

    addLocal(*name);
}

static uint8_t parseVariable(const char* errorMessage) {
    consume(TOKEN_IDENTIFIER, errorMessage);

    declareVariable();
    // If the variable is local (>0) then return a dummy index since we don't
    // nstore local variable names in our constant table.
    if (current->scopeDepth > 0) return 0;

    return identifierConstant(&parser.previous);
}

static void markInitialized() {
    current->locals[current->localCount - 1].depth = current->scopeDepth;
}

// This is where the compiler recognizes that a variable is available for use.
static void defineVariable(uint8_t global) {
    // We leverage the VM's stack to deal with local variables. If the variable
    // is local scope, don't do anything. The value we need is sitting at the
    // top of the VM stack.
    if (current->scopeDepth > 0) {
        markInitialized();
        return;
    }

    emitBytes(OP_DEFINE_GLOBAL, global);
}

static void binary(bool canAssign) {
    TokenType operatorType = parser.previous.type;
    ParseRule* rule = getRule(operatorType);
    parsePrecedence((Precedence)(rule->precedence + 1));

    switch (operatorType) {
        case TOKEN_BANG_EQUAL:    emitBytes(OP_EQUAL, OP_NOT); break;
        case TOKEN_EQUAL_EQUAL:   emitByte(OP_EQUAL); break;
        case TOKEN_GREATER:       emitByte(OP_GREATER); break;
        case TOKEN_GREATER_EQUAL: emitBytes(OP_LESS, OP_NOT); break;
        case TOKEN_LESS:          emitByte(OP_LESS); break;
        case TOKEN_LESS_EQUAL:    emitBytes(OP_GREATER, OP_NOT); break;
        case TOKEN_PLUS:          emitByte(OP_ADD); break;
        case TOKEN_MINUS:         emitByte(OP_SUBTRACT); break;
        case TOKEN_STAR:          emitByte(OP_MULTIPLY); break;
        case TOKEN_SLASH:         emitByte(OP_DIVIDE); break;
        default: return; // Unreachable.
  }
}

static void literal(bool canAssign) {
    switch (parser.previous.type) {
        case TOKEN_FALSE: emitByte(OP_FALSE); break;
        case TOKEN_NIL: emitByte(OP_NIL); break;
        case TOKEN_TRUE: emitByte(OP_TRUE); break;
        default: return; // Unreachable.
    }
}

// Assumes the initial '(' has already been consumed.
static void grouping(bool canAssign) {
    expression();
    consume(TOKEN_RIGHT_PAREN, "expected ')' after expression");
}

// Assumes the number literal has already been consumed and is store in Compiler
// 'previous'.
static void number(bool canAssign) {
    // Convert number literal to double
    double value = strtod(parser.previous.start, NULL);
    
    // Generate code to load value
    emitConstant(NUMBER_VAL(value));
}

// Takes the string’s characters directly from the lexeme,trimes the trailing
// quotes, creates a string object, wraps it in a Value, and stuffs it into the
// constant table.
static void string(bool canAssign) {
    emitConstant(OBJ_VAL(copyString(parser.previous.start + 1,
        parser.previous.length - 2)));
}

// Takes the given identifier token and adds it's lexeme to the chunk’s constant
// table as a string.
static void namedVariable(Token name, bool canAssign) {
    uint8_t getOp, setOp;

    // Try to find local varible with given name
    int arg = resolveLocal(current, &name);

    if (arg != -1) {
        // We found a local variable
        getOp = OP_GET_LOCAL;
        setOp = OP_SET_LOCAL;
    } else {
        // We found a global variable
        arg = identifierConstant(&name);
        getOp = OP_GET_GLOBAL;
        setOp = OP_SET_GLOBAL;
    }

    if (canAssign && match(TOKEN_EQUAL)) {
        expression();
        emitBytes(setOp, (uint8_t)arg);
    } else {
        emitBytes(getOp, (uint8_t)arg);
    }
}

static void variable(bool canAssign) {
    namedVariable(parser.previous, canAssign);
}

// Leading '-' is sitting in previous.
static void unary(bool canAssign) {
    TokenType operatorType = parser.previous.type;

    // Compile the operand.
    parsePrecedence(PREC_UNARY);

    // Emit the operator instruction.
    switch (operatorType) {
        case TOKEN_BANG: emitByte(OP_NOT); break;
        case TOKEN_MINUS: emitByte(OP_NEGATE); break;
        default: return; // Unreachable.
    }
}

ParseRule rules[] = {
    [TOKEN_LEFT_PAREN]    = {grouping, NULL,   PREC_NONE},
    [TOKEN_RIGHT_PAREN]   = {NULL,     NULL,   PREC_NONE},
    [TOKEN_LEFT_BRACE]    = {NULL,     NULL,   PREC_NONE}, 
    [TOKEN_RIGHT_BRACE]   = {NULL,     NULL,   PREC_NONE},
    [TOKEN_COMMA]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_DOT]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_MINUS]         = {unary,    binary, PREC_TERM},
    [TOKEN_PLUS]          = {NULL,     binary, PREC_TERM},
    [TOKEN_SEMICOLON]     = {NULL,     NULL,   PREC_NONE},
    [TOKEN_SLASH]         = {NULL,     binary, PREC_FACTOR},
    [TOKEN_STAR]          = {NULL,     binary, PREC_FACTOR},
    [TOKEN_BANG]          = {unary,    NULL,   PREC_NONE},
    [TOKEN_BANG_EQUAL]    = {NULL,     binary, PREC_EQUALITY},
    [TOKEN_EQUAL]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_EQUAL_EQUAL]   = {NULL,     binary, PREC_EQUALITY},
    [TOKEN_GREATER]       = {NULL,     binary, PREC_COMPARISON},
    [TOKEN_GREATER_EQUAL] = {NULL,     binary, PREC_COMPARISON},
    [TOKEN_LESS]          = {NULL,     binary, PREC_COMPARISON},
    [TOKEN_LESS_EQUAL]    = {NULL,     binary, PREC_COMPARISON},
    [TOKEN_IDENTIFIER]    = {variable, NULL,   PREC_NONE},
    [TOKEN_STRING]        = {string,   NULL,   PREC_NONE},
    [TOKEN_NUMBER]        = {number,   NULL,   PREC_NONE},
    [TOKEN_AND]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_CLASS]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_ELSE]          = {NULL,     NULL,   PREC_NONE},
    [TOKEN_FALSE]         = {literal,  NULL,   PREC_NONE},
    [TOKEN_FOR]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_FUN]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_IF]            = {NULL,     NULL,   PREC_NONE},
    [TOKEN_NIL]           = {literal,  NULL,   PREC_NONE},
    [TOKEN_OR]            = {NULL,     NULL,   PREC_NONE},
    [TOKEN_PRINT]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_RETURN]        = {NULL,     NULL,   PREC_NONE},
    [TOKEN_SUPER]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_THIS]          = {NULL,     NULL,   PREC_NONE},
    [TOKEN_TRUE]          = {literal,  NULL,   PREC_NONE},
    [TOKEN_VAR]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_WHILE]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_ERROR]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_EOF]           = {NULL,     NULL,   PREC_NONE},
};

// 1. Lookup prefrix parser for current token.
// Note: the first token will always belong to some prefix expression.
// 2. Consume 0 or more tokens until prefix expression is done.
// 3. Check for infix parser on next token, if one is found it means the prefix
//    expression we just compiled could be an operand for it.
// Note: This is only the case IF call to parsePrecedence() has a precedence low
// enough to permit that infix operator.
static void parsePrecedence(Precedence precedence) {
    // Read next token
    advance();

    // Lookup token's prefix rule
    ParseFn prefixRule = getRule(parser.previous.type)->prefix;

    // If there is no prefix parser, then we have syntax error
    if (prefixRule == NULL) {
        error("expected expression");
        return;
    }

    bool canAssign = precedence <= PREC_ASSIGNMENT;

    // Valid prefix function, execute it
    prefixRule(canAssign);

    // Check for infix parser on the next token. Keep loopingg through infix
    // operators and their operands until we hit a token that isn’t an infix
    // operator or is too low precedence and stop.
    while (precedence <= getRule(parser.current.type)->precedence) {
        advance();
        ParseFn infixRule = getRule(parser.previous.type)->infix;
        infixRule(canAssign);
    }

    if (canAssign && match(TOKEN_EQUAL)) {
        error("invalid assignment target");
    }
}

// Returns the rule at the given index. This function exists solely to handle a
// declaration cycle in the C code.
static ParseRule* getRule(TokenType type) {
    return &rules[type];
}

static void expression() {
    parsePrecedence(PREC_ASSIGNMENT);
}

static void block() {
    while (!check(TOKEN_RIGHT_BRACE) && !check(TOKEN_EOF)) {
        declaration();
    }

    consume(TOKEN_RIGHT_BRACE, "Expect '}' after block.");
}

// Variable declaration parsing begins in varDeclaration().
//
// First, parseVariable() consumes the identifier token for the variable name,
// adds its lexeme to the chunk’s constant table as a string, and then returns
// the constant table index where it was added.
// 
// Second, after varDeclaration() compiles the initializer, it calls
// defineVariable() to emit the bytecode for storing the variable’s value in
// the global variable hash table.
static void varDeclaration() {
    uint8_t global = parseVariable("expected variable name");

    if (match(TOKEN_EQUAL)) {
        expression();
    } else {
        emitByte(OP_NIL);
    }
    consume(TOKEN_SEMICOLON,
            "expected \";\" after variable declaration");

    defineVariable(global);
}

static void expressionStatement() {
    expression();
    consume(TOKEN_SEMICOLON, "expected \";\" after expression");
    emitByte(OP_POP);
}

static void printStatement() {
    expression();
    consume(TOKEN_SEMICOLON, "expected \";\" after value");
    emitByte(OP_PRINT);
}

static void synchronize() {
  parser.panicMode = false;

    while (parser.current.type != TOKEN_EOF) {
        if (parser.previous.type == TOKEN_SEMICOLON) return;

        switch (parser.current.type) {
        case TOKEN_CLASS:
        case TOKEN_FUN:
        case TOKEN_VAR:
        case TOKEN_FOR:
        case TOKEN_IF:
        case TOKEN_WHILE:
        case TOKEN_PRINT:
        case TOKEN_RETURN:
            return;
        default:
            ; // Do nothing.
        }

        advance();
    }
}

static void declaration() {
    if (match(TOKEN_VAR)) {
        varDeclaration();
    } else {
        statement();
    }

    if (parser.panicMode) synchronize();
}

static void statement() {
    if (match(TOKEN_PRINT)) {
        printStatement();
    } else if (match(TOKEN_LEFT_BRACE)) {
        beginScope();
        block();
        endScope();
    } else {
        expressionStatement();
    }
}

bool compile(const char* source, Chunk* chunk) {
    initScanner(source);
    Compiler compiler;
    initCompiler(&compiler);
    compilingChunk = chunk;

    // Initialize compiler panic settings
    parser.hadError = false;
    parser.panicMode = false;

    advance();

    while (!match(TOKEN_EOF)) {
        declaration();
    }

    endCompiler();
    return !parser.hadError;
}