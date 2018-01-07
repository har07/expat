#include <expat.h>
#include <string.h>
#include <stdio.h>
#include "_cgo_export.h"

#define POOLSIZE        8192

int lastId;

typedef struct parserInstance {
    int id;
    XML_Parser parser;
    int start_handler;
    int end_handler;
} ParserInstance;

typedef enum {ELEMENT_START, ELEMENT_END} EventType;

ParserInstance PI;

static void XMLCALL _handle_start_(void *data, const XML_Char *el, const XML_Char **attr)
{
    int i;
    int count = 0;
    (void)data;

    XML_Char *tag;
    tag = malloc((strlen(el)+1)*sizeof(XML_Char));
    strcpy(tag, el);

    // create copy of attributes to be passed to Go handler
    for (i = 0; attr[i]; i++) {
        count++;
    }
    XML_Char **attribcpy = malloc((count+1)*sizeof(XML_Char*));
    for (i = 0; i<count; i++) {
        attribcpy[i] = malloc((strlen(attr[i]))*sizeof(XML_Char));
        strcpy(attribcpy[i], attr[i]);
    }

    GStartElementHandler(PI.id, tag, attribcpy);
}

static void XMLCALL _handle_end_(void *data, const XML_Char *el)
{
    (void)data;
    (void)el;

    XML_Char *tag;
    tag = malloc(strlen(el)+1);
    strcpy(tag, el);

    GEndElementHandler(PI.id, tag);
}

static void XMLCALL _handle_char_data_(void *data, const XML_Char *s, int len)
{
    XML_Char *text;
    text = malloc(len*sizeof(XML_Char));
    strcpy(text, s);

    GCharDataHandler(PI.id, text, len);
}

static void XMLCALL _handle_default_(void *data, const XML_Char *s, int len)
{
    XML_Char *text;
    text = malloc(len*sizeof(XML_Char));
    strcpy(text, s);

    GDefaultHandler(PI.id, text, len);
}

int Create(XML_Char *encoding, int namespace){
    XML_Parser p;
    if(namespace){
        p = XML_ParserCreate(encoding);
    } else {
        p = XML_ParserCreateNS(encoding, ':');
    }
    XML_SetElementHandler(p, _handle_start_, _handle_end_);
    XML_SetCharacterDataHandler(p, _handle_char_data_);
    XML_SetDefaultHandler(p, _handle_default_);
    ParserInstance pi;
    pi.id = lastId++;
    pi.parser = p;
    
    PI = pi;
    return pi.id;
}

int Feed(int id, XML_Char *chunk, int len, int finish){
    if (XML_Parse(PI.parser, chunk, len, 1) == XML_STATUS_ERROR) {
        // ParseError err;
        // const char * temp = XML_ErrorString(XML_GetErrorCode(PI.parser));
        // char *msg;
        // msg = malloc(strlen(temp)+1);
        // strcpy(msg, temp);
        // err.message = msg;
        // err.line = XML_GetCurrentLineNumber(PI.parser);
        // return err;
        return XML_GetErrorCode(PI.parser);
  }
  return 0;
}

char* GetError(int id, int code){
    const char * temp = XML_ErrorString(XML_GetErrorCode(PI.parser));
    char *msg;
    msg = malloc(strlen(temp)+1);
    strcpy(msg, temp);
    return msg;
}

int GetCurrentLineNumber(int id){
    return XML_GetCurrentLineNumber(PI.parser);
}

int GetCurrentColumnNumber(int id){
    return XML_GetCurrentColumnNumber(PI.parser);
}

int GetCurrentAttributeCount(int id){
    return XML_GetSpecifiedAttributeCount(PI.parser);
}

void SetHandlers(int id, int start, int end){
    PI.start_handler = start;
    PI.end_handler = end;
}

void Free(int id){
    XML_ParserFree(PI.parser);
}

