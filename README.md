# OTS

ots is an interface to libots - The [Open Text Summarizer](http://libots.sourceforge.net/).

## Dependencies

  * libxml2
  * glib2.0

## Installation

### Debian flavors of Linux

```
  sudo apt-get install pkg-config libxml2-dev libglib2.0-dev
  go get github.com/shanna/ots
```

### OSX

```
  brew install pkg-config libxml2 glib
  go get github.com/shanna/ots
```

## Usage

```go
  import 'github.com/shanna/ots'

  article := ots.Parse("I think I need some ice cream to cool me off. It is too hot down under", "en")
  article = ots.Parse("j'ai besoin de la crème glacée. il fait trop chaud en australie.", "fr")

  article.Keywords()
  article.Percent(50)
  article.Sentences(1)

  ots.Languages() #=> list of supported language dictionaries baked-in to libots
```

## See

  * [https://github.com/deepfryed/ots](https://github.com/deepfryed/ots)
  * [https://github.com/ssoper/summarize](https://github.com/ssoper/summarize)

## License

MIT

