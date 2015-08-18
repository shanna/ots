package ots

/*
#cgo pkg-config: glib-2.0 libxml-2.0
#include "ots.h"
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

  C.ots_dictionary_dir_set(cdir)
}

func Languages() (languages []string, err error) {
  files, err := ioutil.ReadDir(C.GoString(C.ots_dictionary_dir_get()))

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

  C.ots_article_parse(article.pointer, ctext, C.size_t(len(text)));

  return article, nil
}

func (a Article) Keywords() []string {
  title := C.GoString(a.pointer.title)
  return strings.Split(title, ",")
}

type Summary struct {
  Sentences []Sentence
}

type Sentence struct {
  Sentence string
  Score    float64
}

//export summary_append
func summary_append(csummary unsafe.Pointer, csentence *C.char, cscore C.float) {
  summary  := (*Summary)(csummary)
  sentence := C.GoString(csentence)
  score    := float64(cscore)

  summary.Sentences = append(summary.Sentences, Sentence{sentence, score})
}

func (a Article) Sentences(sentences int) *Summary {
  C.ots_highlight_doc_lines(a.pointer, C.int(sentences))

  summary := &Summary{Sentences: []Sentence{}} // TODO: NewSummary
  C.ots_article_summary(a.pointer, unsafe.Pointer(summary))
  return summary
}

func (a Article) Percentage(percentage int) *Summary {
  C.ots_highlight_doc(a.pointer, C.int(percentage))

  summary := &Summary{Sentences: []Sentence{}} // TODO: NewSummary
  C.ots_article_summary(a.pointer, unsafe.Pointer(summary))
  return summary
}
