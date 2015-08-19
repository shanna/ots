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
	// Set the dictionary path relative to the installed library path.
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
	pointer *C.OtsArticle

	// Language used to parse the article text.
	Language string
}

func Parse(text string, language string) (*Article, error) {
	article := &Article{
		pointer:  C.ots_new_article(),
		Language: language,
	}

	ctext := C.CString(text)
	clanguage := C.CString(language)
	defer C.free(unsafe.Pointer(ctext))
	defer C.free(unsafe.Pointer(clanguage))

	if C.ots_load_xml_dictionary(article.pointer, clanguage) != C.TRUE {
		return nil, fmt.Errorf("No dictionary for language: %s", language)
	}

	C.ots_article_parse(article.pointer, ctext, C.size_t(len(text)))
	return article, nil
}

func (a Article) Keywords() []string {
	title := C.GoString(a.pointer.title)
	return strings.Split(title, ",") // TODO: Sort?
}

// I can't figure out how to pass the sentences slice array Go -> C -> Go so I'm
// wrapping it in this struct which is easy to pass using an unsafe.Pointer.
type summary struct {
	Sentences Sentences
}

type Sentence struct {
	Text  string
	Score float64
}

// Sortable []Sentence collection by the Sentence.Score field.
// See sort.Sort() and the sort.Interface.
type Sentences []Sentence

func (s Sentences) Len() int           { return len(s) }
func (s Sentences) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Sentences) Less(i, j int) bool { return s[i].Score < s[j].Score }

//export summary_append
func summary_append(csummary unsafe.Pointer, csentence *C.char, cscore C.float) {
	s := (*summary)(csummary)
	s.Sentences = append(s.Sentences, Sentence{C.GoString(csentence), float64(cscore)})
}

func (a Article) sentences() Sentences {
	s := &summary{}
	C.ots_article_summary(a.pointer, unsafe.Pointer(s))
	sort.Sort(s.Sentences)
	return s.Sentences
}

func (a Article) Sentences(sentences int) Sentences {
	C.ots_highlight_doc_lines(a.pointer, C.int(sentences))
	return a.sentences()
}

func (a Article) Percentage(percentage int) Sentences {
	C.ots_highlight_doc(a.pointer, C.int(percentage))
	return a.sentences()
}
