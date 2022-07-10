#include "scanner.h"

typedef struct {
    const char* start;
    const char* current;
    int line;
} Scanner;

Scanner scanner;

void initScanner(const char* source) {
    scanner.start = source;
    scanner.current = source;
    scanner.line = 1;

    if (DEBUG_LOX & DF_SCANNING) {
        Scanner tmp = scanner;

        Scanner scannerCopy;
        scannerCopy.start = source;
        scannerCopy.current = source;
        scannerCopy.line = 1;

        scanner = scannerCopy;

        Token current;
        for (;;) {
            current = scanToken();
            printToken(current);
            if (current.type == TOKEN_EOF) break;
        }

        scanner = tmp;
    }
}

static bool isAlpha(char c) {
  return (c >= 'a' && c <= 'z') ||
         (c >= 'A' && c <= 'Z') ||
          c == '_';
}

static bool isDigit(char c) {
    return c >= '0' && c <= '9';
}

static bool isAtEnd() {
    // Valid source code is required to be null terminated
    return *scanner.current == '\0';
}

// Consume the next character and return it.
static char advance() {
    scanner.current++;
    return scanner.current[-1];
}

// Return the Scanner's current character but does not consume it.
static char peek() {
    return *scanner.current;
}

// Returns the character past the Scanner's current one.
static char peekNext() {
    if (isAtEnd()) return '\0';
    return scanner.current[1];
}

static Token makeToken(TokenType type) {
    Token token;

    token.type = type;

    // Pointer to first character of token's lexeme
    token.start = scanner.start;
    // The caller believes they have reached lexeme end when calling. Therefore,
    // the Scanner's current pointer is the end of token (\0) and length can
    // be simply computed from current less start using pointer arithmetic.
    token.length = (int)(scanner.current - scanner.start);

    token.line = scanner.line;

    return token;
}

static bool match(char expected) {
    if (isAtEnd()) return false;

    // Scanner's current character does not match. Do not advance and return
    // false.
    if (*scanner.current != expected) return false;

    // Scanner's current character matches expected. Advance and return true.
    scanner.current++;

    return true;
}

static Token errorToken(const char* message) {
    Token token;

    token.type = TOKEN_ERROR;

    // The only difference from make token is that instead of pointing to the
    // token's source code lexeme, we point to an error message.
    token.start = message;
    token.length = (int)strlen(message);

    token.line = scanner.line;

    return token;
}

// Skip past any whitespace to first meaningful character
static void skipWhitespace() {
    for (;;) {
        char c = peek();
        switch (c) {
        case ' ':
        case '\r':
        case '\t':
            advance();
            break;
        case '\n':
            // Consume newline character but also increase our line count to
            // represent that we've advanced source code lines.
            scanner.line++;
            advance();
            break;
        case '/': // Skip comments
            if (peekNext() == '/') {
                // A comment goes until the end of the line. We want consume all
                // characters *before* the end of the line. The '\n' must be
                // scanned so the Scanner can advance the line count internally.
                while (peek() != '\n' && !isAtEnd()) advance();
            } else {
                return;
            }
            break;
        default:
            return;
        }
    }
}

// Utility function to test the rest of potential identifier keyword
static TokenType checkKeyword(int start, int length, const char* rest,
    TokenType type) {

    // Verify:
    // 1. Lexeme length matches keyword length
    // 2. Remaining characters of lexeme match keyword exactly
    if (scanner.current - scanner.start == start + length &&
        memcmp(scanner.start + start, rest, length) == 0) {
        return type;
    }

    return TOKEN_IDENTIFIER;
}

// Detects a keyword and failes fast if identifier is not a reserved word.
// Implemented using tiny keyword tree using trie data structure.
static TokenType identifierType() {
    switch (scanner.start[0]) {
        case 'a': return checkKeyword(1, 2, "nd", TOKEN_AND);
        case 'c': return checkKeyword(1, 4, "lass", TOKEN_CLASS);
        case 'e': return checkKeyword(1, 3, "lse", TOKEN_ELSE);
        case 'f':
            // Verify that a second letter exists based on Scanner progress
            if (scanner.current - scanner.start > 1) {
                switch (scanner.start[1]) {
                    case 'a': return checkKeyword(2, 3, "lse", TOKEN_FALSE);
                    case 'o': return checkKeyword(2, 1, "r", TOKEN_FOR);
                    case 'u': return checkKeyword(2, 1, "n", TOKEN_FUN);
                }
            }
            break;
        case 'i': return checkKeyword(1, 1, "f", TOKEN_IF);
        case 'n': return checkKeyword(1, 2, "il", TOKEN_NIL);
        case 'o': return checkKeyword(1, 1, "r", TOKEN_OR);
        case 'p': return checkKeyword(1, 4, "rint", TOKEN_PRINT);
        case 'r': return checkKeyword(1, 5, "eturn", TOKEN_RETURN);
        case 's': return checkKeyword(1, 4, "uper", TOKEN_SUPER);
        case 't':
            // Verify that a second letter exists based on Scanner progress
            if (scanner.current - scanner.start > 1) {
                switch (scanner.start[1]) {
                    case 'h': return checkKeyword(2, 2, "is", TOKEN_THIS);
                    case 'r': return checkKeyword(2, 2, "ue", TOKEN_TRUE);
                }
            }
        break;
        case 'v': return checkKeyword(1, 2, "ar", TOKEN_VAR);
        case 'w': return checkKeyword(1, 4, "hile", TOKEN_WHILE);
    }

    return TOKEN_IDENTIFIER;
}

// Assumes caller validated the starting character letter before calling. clox
// identifiers must start with letters.
static Token identifier() {
    // Consume all alphanumeric characters until first non-alphanumeric char
    while (isAlpha(peek()) || isDigit(peek())) advance();

    // Produce the "proper" token type
    return makeToken(identifierType());
}

static Token number() {
    // Consume all digit character until first non-digit
    while (isDigit(peek())) advance();

    // Check if the first non-digit char is a '.' suggesting a fractional value
    // and the next char is indeed another digit
    if (peek() == '.' && isDigit(peekNext())) {
        // Consume the "."
        advance();

        // Consume all digit character until next non-digit
        while (isDigit(peek())) advance();
    }

    return makeToken(TOKEN_NUMBER);
}

static Token string() {
    // Consume all character until reaching the ending '"'
    while (peek() != '"' && !isAtEnd()) {
        // Multi-line strings are valid, we check for a new line and update
        // Scanner line count accordingly.
        if (peek() == '\n') scanner.line++;

        // Consume the current char, update Scanner current to next char.
        advance();
    }

    // Check if we brok out of loop above because of EOF and report error if so.
    if (isAtEnd()) return errorToken("unterminated string");

    // Consume the closing quote.
    advance();

    return makeToken(TOKEN_STRING);
}

Token scanToken() {
    // Advance the Scanner past any leading whitespace
    skipWhitespace();

    // Point to the current character (start of token to scan)
    scanner.start = scanner.current;

    if (isAtEnd()) return makeToken(TOKEN_EOF);

    // Read next character from the source code.
    char c = advance();

    // Identifier and keyword tokens
    if (isAlpha(c)) return identifier();

    // Literal token - numbers
    // If we have a digit, catch it here instead of switch processing
    if (isDigit(c)) return number();

    switch (c) {
        // Single character tokens
        case '(': return makeToken(TOKEN_LEFT_PAREN);
        case ')': return makeToken(TOKEN_RIGHT_PAREN);
        case '{': return makeToken(TOKEN_LEFT_BRACE);
        case '}': return makeToken(TOKEN_RIGHT_BRACE);
        case ';': return makeToken(TOKEN_SEMICOLON);
        case ',': return makeToken(TOKEN_COMMA);
        case '.': return makeToken(TOKEN_DOT);
        case '-': return makeToken(TOKEN_MINUS);
        case '+': return makeToken(TOKEN_PLUS);
        case '/': return makeToken(TOKEN_SLASH);
        case '*': return makeToken(TOKEN_STAR);
        // Two-character tokens
        // After consuming the first character we look for '='. If found, we
        // consume it and return the corresponding two-character token.
        //
        // If not found, we leave the current character alone so it can be part
        // of the next token scan. The one-character token for the first consume
        // is then returned instead.
        case '!':
            return makeToken(
                match('=') ? TOKEN_BANG_EQUAL : TOKEN_BANG);
        case '=':
            return makeToken(
                match('=') ? TOKEN_EQUAL_EQUAL : TOKEN_EQUAL);
        case '<':
            return makeToken(
                match('=') ? TOKEN_LESS_EQUAL : TOKEN_LESS);
        case '>':
            return makeToken(
                match('=') ? TOKEN_GREATER_EQUAL : TOKEN_GREATER);
        // Literal token - strings
        case '"': return string();
    }

    return errorToken("unexpected character");
}

void printToken(Token tok) {
    printf("type: %d lexeme: %.*s line: %d\n", tok.type, tok.length, tok.start, tok.line);
}