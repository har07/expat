typedef struct parseError {
    char *message;
    int line;
    int col;
} ParseError;