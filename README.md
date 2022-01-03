# 8_bit_hubble_golang

Replacement Branch: certain C functionality is replaced (C-Golang bindings)

The goal is not to 100% imitate the C version functionality, but to understand how to communicate Go code and C code. If the outputs obtained from C version and Golang versions are identical in appearence (i.e: metadata can differ), the replacement is considered successfull. To understand the changes, a little context may be necessary (check original 8_bit_hubble project)

Links used:
- https://go.dev/blog/cgo
- https://zchee.github.io/golang-wiki/cgo/
- https://medium.com/@ben.mcclelland/an-adventure-into-cgo-calling-go-code-with-c-b20aa6637e75

Commits are self-explanatory, and code more also. But in summary:
- Golang provides as part of its standard library the "image" package (and others related to it), which provides functionality to create images (like PNG images). So the goal is to replace the "libpng" usage from C code, and use insthead that go package.
- So, to do that: a Golang package is created with the required functionality to create an image using as input the C data (a dynamic-allocated array of hex values, each one representing the color of the pixel)
- Then the go code is modified to make it exportable (see CGO documentation)
- Then the C code is modified to call the Golang functions
- The cgo generates both new header files, and a static library. This static library is the one to be linked to the C code, so the C functionality can execute the Golang code.
- Profit!

Core command:
`go build -buildmode=c-archive -o libimageWriter.a imageWriter.go`

The previous command takes the ".go" code, and generates both a ".h" header file to be included in the C code, and the static library to be linked to the C executable.
