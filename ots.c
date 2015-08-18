#include "ots.h"

// Don't blame me, blame libots. It looks like a constant but isn't.
char *DICTIONARY_DIR = NULL;

// Not thread safe.
void ots_dictionary_dir_set(const char *path) {
  if (DICTIONARY_DIR)
    free(DICTIONARY_DIR);

  DICTIONARY_DIR = (char *)malloc(strlen(path) + 2); // TODO: Expansion and cleanup in Go.
  sprintf(DICTIONARY_DIR, "%s/", path);
}

// Not thread safe.
char *ots_dictionary_dir_get(void) {
  return DICTIONARY_DIR;
}

void ots_article_parse(OtsArticle *article, const char *text, size_t size) {
  ots_parse_stream(text, size, article);
  ots_grade_doc(article);
}

void ots_article_summary(OtsArticle *article, void *summary) {
  for (GList *line = article->lines; line != NULL; line = g_list_next(line)) {
    OtsSentence *sentence = (OtsSentence *)line->data;
    if (!sentence->selected) continue;

    size_t size;
    unsigned char *content = ots_get_line_text(sentence, TRUE, &size);
    summary_append(summary, content, sentence->score);
    sentence->selected = FALSE;
  }
}
