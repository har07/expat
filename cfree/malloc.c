#include <stdlib.h>
#include <string.h>

char* GetError(){
    char dummy[] = "Invalid character";
    char *msg;
    msg = malloc(sizeof(dummy));
    strcpy(msg, dummy);
    return msg;
}