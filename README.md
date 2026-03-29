# imageset

`imageset` is a Go module for DayZ `.imageset` files.

It provides:

* parser from file, reader, bytes, or string
* canonical text formatter/writer
* semantic validation with structured diagnostics
* `lintkit/lint` provider registration (`RegisterLintRules`) and codes catalog
* identifier normalization helpers
* symbolic and numeric flags parsing
* generated lint rules snapshot and docs: [rules.yaml](rules.yaml),
  [RULES.md](RULES.md)

## Install

```bash
go get github.com/woozymasta/imageset
```

## Quick Example

```go
package main

import (
    "log"

    "github.com/woozymasta/imageset"
)

func main() {
    doc, err := imageset.ParseFile("ui.imageset")
    if err != nil {
        log.Fatal(err)
    }

    if err := imageset.Validate(doc); err != nil {
        log.Fatal(err)
    }

    err = imageset.WriteFile("ui_out.imageset", doc, &imageset.FormatOptions{
        UseCamelCaseNames: false,
    })
    if err != nil {
        log.Fatal(err)
    }
}
```
