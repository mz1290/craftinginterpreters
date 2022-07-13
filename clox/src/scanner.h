#ifndef clox_scanner_h
#define clox_scanner_h


typedef enum {
    // Single-character tokens.
    TOKEN_LEFT_PAREN, TOKEN_RIGHT_PAREN,
    TOKEN_LEFT_BRACE, TOKEN_RIGHT_BRACE,
    TOKEN_COMMA, TOKEN_DOT, TOKEN_MINUS, TOKEN_PLUS,
    TOKEN_SEMICOLON, TOKEN_SLASH, TOKEN_STAR,

    // One or two character tokens.
    TOKEN_BANG, TOKEN_BANG_EQUAL,
    TOKEN_EQUAL, TOKEN_EQUAL_EQUAL,
    TOKEN_GREATER, TOKEN_GREATER_EQUAL,
    TOKEN_LESS, TOKEN_LESS_EQUAL,

    // Literals.
    TOKEN_IDENTIFIER, TOKEN_STRING, TOKEN_NUMBER,

    // Keywords.
    TOKEN_AND, TOKEN_CLASS, TOKEN_ELSE, TOKEN_FALSE,
    TOKEN_FOR, TOKEN_FUN, TOKEN_IF, TOKEN_NIL, TOKEN_OR,
    TOKEN_PRINT, TOKEN_RETURN, TOKEN_SUPER, TOKEN_THIS,
    TOKEN_TRUE, TOKEN_VAR, TOKEN_WHILE,

    // Reports unterminated strings and unrecognized characters. The scanner
    // will produce a "error" token and pass to the compiler. The compiler then
    // knows an error occured and begins error recovery before reporting.
    TOKEN_ERROR,

    TOKEN_EOF
} TokenType;

typedef struct {
    TokenType type;
    // lexeme start
    // A lexeme is the character sequence exactly as it appears in the user’s
    // source code.
    const char* start;
    int length;
    int line;
} Token;

void initScanner(const char*);
void printToken(Token tok);

// Scan and return the next token from the source code.
Token scanToken();

#endif