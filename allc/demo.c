#include <stdio.h>
#include <expat.h>

#define BUFFSIZE        8192

char Buff[BUFFSIZE];

int Depth;

static void XMLCALL
start(void *data, const XML_Char *el, const XML_Char **attr)
{
  int i;
  (void)data;

  for (i = 0; i < Depth; i++)
    printf("  ");

  printf("%s", el);

  for (i = 0; attr[i]; i += 2) {
    printf(" %s='%s'", attr[i], attr[i + 1]);
  }

  printf("\n");
  Depth++;
}

static void XMLCALL
end(void *data, const XML_Char *el)
{
  (void)data;
  (void)el;

  int i;
  for (i = 0; i < Depth; i++)
    printf("  ");

  printf("end element %s\n", el);

  Depth--;
}

int
demo(char *data, int len)
{
  XML_Parser p = XML_ParserCreate(NULL);

  if (! p) {
    fprintf(stderr, "Couldn't allocate memory for parser\n");
    return -1;
  }

  XML_SetElementHandler(p, start, end);

  if (XML_Parse(p, data, len, 1) == XML_STATUS_ERROR) {
    fprintf(stderr,
            "Parse error at line %lu:\n%s\n",
            XML_GetCurrentLineNumber(p),
            XML_ErrorString(XML_GetErrorCode(p)));
    return -1;
  }
  XML_ParserFree(p);
  return 0;
}
