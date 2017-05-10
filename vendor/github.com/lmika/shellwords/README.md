# shellwords

A simple string tokenizer for Go, inspired by Ruby's 
[Shellwords](http://ruby-doc.org/stdlib-2.0/libdoc/shellwords/rdoc/Shellwords.html) module.

Shellwords supports tokens separated by whitespace or delimited by single or double quotes.

To get:

```
go get github.com/lmika/shellwords
```

To use:

```go
package main

import "fmt"
import "shellwords"

func main() {
    shellwords.Split("these 'a three' tokens")  // ["these", "a three", "tokens"]
}
```

Documentation at http://godoc.org/github.com/lmika/shellwords
