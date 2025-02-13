package ui

import "embed"

//go:embed "html" "static"
var Files embed.FS

/*//go:embed "html" "static"
This is a compiler directive that tells the Go compiler to embed the html and static directories into the binary.

The //go:embed directive must be placed immediately above a variable declaration of type embed.FS (or a compatible type like string or []byte for single files).

What It Does
The html and static directories (relative to the location of the Go source file) are embedded into the binary.

At compile time, the Go compiler reads the contents of these directories and includes them in the executable.*/
