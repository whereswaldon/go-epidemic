# go-epidemic
A vim/neovim pathogen plugin autoupdater

## Install
Requires `go`. If you don't have `go` installed, please follow the directions [here](https://golang.org/doc/install).
Make sure that `$GOPATH/bin` is in your `$PATH`.
Run `go get github.com/whereswaldon/go-epidemic`.

## Usage
Run `go-epidemic <path-to-your-plugin-directory>`.
Usually, this will look like `go-epidemic ~/.vim/bundle` or `go-epidemic ~/.config/nvim/bundle`
Set this up on a schedule with the scheduling utility of your choice to update your
plugins at set intervals.

## License
MIT
