#pragma once

#include <stdlib.h>
#include <stdio.h>

#include "libots.h"

void ots_dictionary_dir_set(const char *path);
char *ots_dictionary_dir_get(void);

void ots_article_parse(OtsArticle *article, const char *text, size_t size);
void ots_article_summary(OtsArticle *article, void *summary);

// Go
extern void summary_append(void *, char *, float);

