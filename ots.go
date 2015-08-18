package ots

/*
#cgo pkg-config: glib-2.0 libxml-2.0

#include <stdlib.h>
#include <stdio.h>

// Don't blame me, blame libots. It looks like a constant but isn't.
char *DICTIONARY_DIR = NULL;

#include "libots.h"

// Not thread safe.
static void ots_set_dictionary_dir(const char *path) {
  if (DICTIONARY_DIR)
    free(DICTIONARY_DIR);

  DICTIONARY_DIR = (char *)malloc(strlen(path) + 2); // TODO: Expansion and cleanup in Go.
  sprintf(DICTIONARY_DIR, "%s/", path);
}

// Not thread safe.
static char *ots_get_dictionary_dir(void) {
  return DICTIONARY_DIR;
}

static void ots_parse(OtsArticle *article, const char *text, size_t size) {
  ots_parse_stream(text, size, article);
  ots_grade_doc(article);
}

extern void summary_append(void *, char *, float);

static void ots_article_summary2(OtsArticle *article, void *summary) {
  for (GList *line = article->lines; line != NULL; line = g_list_next(line)) {
    OtsSentence *sentence = (OtsSentence *)line->data;
    if (!sentence->selected) continue;

    size_t size;
    unsigned char *content = ots_get_line_text(sentence, TRUE, &size);

    summary_append(summary, content, sentence->score);

    // reset this so subsequent calls work right.
    // sentence->selected = FALSE;
  }
}
*/
import "C"

import (
  "fmt"
  "io/ioutil"
  "path/filepath"
  "runtime"
  "sort"
  "strings"
  "unsafe"
)

func init() {
  _, filename, _, _ := runtime.Caller(0)
  cdir := C.CString(filepath.Join(filepath.Dir(filename), "dictionaries"))
  defer C.free(unsafe.Pointer(cdir))

  C.ots_set_dictionary_dir(cdir)
}

func Languages() (languages []string, err error) {
  files, err := ioutil.ReadDir(C.GoString(C.ots_get_dictionary_dir()))

  for _, f := range files {
    if extension := filepath.Ext(f.Name()); extension == ".xml" {
      languages = append(languages, strings.TrimSuffix(f.Name(), extension))
    }
  }

  sort.Strings(languages)
  return languages, err
}

type Article struct {
  pointer  *C.OtsArticle
  Language string
}

func Parse(text string, language string) (*Article, error) {
  article := &Article{
    pointer:  C.ots_new_article(),
    Language: language,
  }

  ctext     := C.CString(text)
  clanguage := C.CString(language)
  defer C.free(unsafe.Pointer(ctext))
  defer C.free(unsafe.Pointer(clanguage))

  if C.ots_load_xml_dictionary(article.pointer, clanguage) != C.TRUE {
    return nil, fmt.Errorf("No dictionary for language: %s", language)
  }

  C.ots_parse(article.pointer, ctext, C.size_t(len(text)));

  return article, nil
}

func (a Article) Keywords() []string {
  title := C.GoString(a.pointer.title)
  return strings.Split(title, ",")
}

type Summary struct {
  Sentence string
  Score    float32
}

//export summary_append
func summary_append(csummary unsafe.Pointer, csentence *C.char, cscore C.float) {
  /*
  summary  := ([]Summary)(csummary)
  sentence := C.GoString(csentence)
  score    := float64(cscore)

  summary = append(summary, Summary{sentence, score})
  */
}

func (a Article) Sentences(sentences int) []Summary {
  C.ots_highlight_doc_lines(a.pointer, C.int(sentences))

  summary := []Summary{}
  C.ots_article_summary2(a.pointer, unsafe.Pointer(summary))
  return summary
}

/*
func (a Article) Percentage(percentage int) string {
  C.ots_highlight_doc(a.pointer, C.int(percentage))
}
*/
