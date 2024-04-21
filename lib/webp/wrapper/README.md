# CGO Wrapper for libwebp

Wraps the libwebp source code in a way that CGO can compile it.

## Updating to a new libwebp

1. Update the libwebp git submodule to the wanted version tag.
2. run `go generate .` inside this directory.
