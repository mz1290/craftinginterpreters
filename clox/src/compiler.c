#include "compiler.h"


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
typedef void (*ParseFn)();

typedef struct {
    ParseFn prefix;
    ParseFn infix;
    Precedence precedence;
} ParseRule;

// Pattern used throughout clox. Using a single global variable allows us to
// pass the state around from function to function in Compiler.
Parser parser;
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
        fprintf(stderr, " at '%.*s'", token->length, token->start);
    }

    // Print error message
    fprintf(stderr, ": %s\n", message);

    // Signal error occurred during compilation
    parser.hadError = true;
}

// Report error to uaser at previous token location
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

// These forward declaration enable clox's grammar recursive nature. For example
// binary() is defined *before* the rules table so the table can store a pointer
// to it. This prevents binary() function body from being able to access the
// table directly. To allow binary() to access the table, we define function
// getRule() *after* the table and forward declare here so binary() can access
// the table when needed.
static void expression();
static ParseRule* getRule(TokenType type);
static void parsePrecedence(Precedence precedence);

static void binary() {
    TokenType operatorType = parser.previous.type;
    ParseRule* rule = getRule(operatorType);
    parsePrecedence((Precedence)(rule->precedence + 1));

    switch (operatorType) {
        case TOKEN_PLUS:          emitByte(OP_ADD); break;
        case TOKEN_MINUS:         emitByte(OP_SUBTRACT); break;
        case TOKEN_STAR:          emitByte(OP_MULTIPLY); break;
        case TOKEN_SLASH:         emitByte(OP_DIVIDE); break;
        default: return; // Unreachable.
  }
}

// Assumes the initial '(' has already been consumed.
static void grouping() {
    expression();
    consume(TOKEN_RIGHT_PAREN, "expected ')' after expression");
}

// Assumes the number literal has already been consumed and is store in Compiler
// 'previous'.
static void number() {
    // Convert number literal to double
    double value = strtod(parser.previous.start, NULL);
    
    // Generate code to load value
    emitConstant(value);
}

// Leading '-' is sitting in previous.
static void unary() {
    TokenType operatorType = parser.previous.type;

    // Compile the operand.
    parsePrecedence(PREC_UNARY);

    // Emit the operator instruction.
    switch (operatorType) {
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
    [TOKEN_BANG]          = {NULL,     NULL,   PREC_NONE},
    [TOKEN_BANG_EQUAL]    = {NULL,     NULL,   PREC_NONE},
    [TOKEN_EQUAL]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_EQUAL_EQUAL]   = {NULL,     NULL,   PREC_NONE},
    [TOKEN_GREATER]       = {NULL,     NULL,   PREC_NONE},
    [TOKEN_GREATER_EQUAL] = {NULL,     NULL,   PREC_NONE},
    [TOKEN_LESS]          = {NULL,     NULL,   PREC_NONE},
    [TOKEN_LESS_EQUAL]    = {NULL,     NULL,   PREC_NONE},
    [TOKEN_IDENTIFIER]    = {NULL,     NULL,   PREC_NONE},
    [TOKEN_STRING]        = {NULL,     NULL,   PREC_NONE},
    [TOKEN_NUMBER]        = {number,   NULL,   PREC_NONE},
    [TOKEN_AND]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_CLASS]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_ELSE]          = {NULL,     NULL,   PREC_NONE},
    [TOKEN_FALSE]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_FOR]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_FUN]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_IF]            = {NULL,     NULL,   PREC_NONE},
    [TOKEN_NIL]           = {NULL,     NULL,   PREC_NONE},
    [TOKEN_OR]            = {NULL,     NULL,   PREC_NONE},
    [TOKEN_PRINT]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_RETURN]        = {NULL,     NULL,   PREC_NONE},
    [TOKEN_SUPER]         = {NULL,     NULL,   PREC_NONE},
    [TOKEN_THIS]          = {NULL,     NULL,   PREC_NONE},
    [TOKEN_TRUE]          = {NULL,     NULL,   PREC_NONE},
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

    // Valid prefix function, execute it
    prefixRule();

    // Check for infix parser on the next token. Keep loopingg through infix
    // operators and their operands until we hit a token that isn’t an infix
    // operator or is too low precedence and stop.
    while (precedence <= getRule(parser.current.type)->precedence) {
        advance();
        ParseFn infixRule = getRule(parser.previous.type)->infix;
        infixRule();
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

bool compile(const char* source, Chunk* chunk) {
    initScanner(source);
    compilingChunk = chunk;

    // Initialize compiler panic settings
    parser.hadError = false;
    parser.panicMode = false;

    advance();

    // Parse a single expression
    expression();

    // We should have EOF aftersource code compilation
    consume(TOKEN_EOF, "expected end of expression");

    endCompiler();
    return !parser.hadError;
}