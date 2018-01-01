#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char** get_strarr(){
    char **arr;
    int count = 8;
    int i;

    arr = malloc((count+1) * sizeof(char*));
    for(i=0; i<count; i++){
        arr[i] = malloc(5*sizeof(char));
        sprintf(arr[i], "foo%d", i);
    }
    // arr[count] = malloc(sizeof(char));
    arr[count] = NULL;
    return arr;
}

int demo(){
    char **attr = get_strarr();
    int i;
    // (void)attr;
    // (void)i;
    // for(i=0; i<8; i++){
    for(i=0; attr[i]; i++){
        printf("%s\n", attr[i]);
        free(attr[i]);
    }
    free(attr[i]);
    free(attr);
}