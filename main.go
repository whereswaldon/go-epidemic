package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"os"
)

var red (func(...interface{}) string) = color.New(color.FgRed).SprintFunc()
var yellow (func(...interface{}) string) = color.New(color.FgYellow).SprintFunc()
var green (func(...interface{}) string) = color.New(color.FgGreen).SprintFunc()
var blue (func(...interface{}) string) = color.New(color.FgCyan).SprintFunc()

const USAGE_MESSAGE string = `
%[1]s path

path -- the path to where your pathogen vim plugins are stored.
	This is generally the bundle directory within your vim
	configuration.
	Ex: (neovim) $HOME/.config/nvim/bundle
	    (vim)    $HOME/.vim/bundle

Ex: (neovim) %[1]s $HOME/.config/nvim/bundle
    (vim)    %[1]s $HOME/.vim/bundle
`

func main() {
	// Override usage text
	flag.Usage = printUsage

	// Parse command line arguments
	pluginPath := flag.Arg(0)
	if pluginPath == "" {
		printUsage()
	}

	// Find plugins
	pluginDir, err := getDirectory(pluginPath)
	if err != nil {
		fmt.Println(os.Stderr, red(err))
		os.Exit(1)
	}
	fmt.Println(blue(pluginDir.Name()) + " targeted")
}

// printUsage prints usage information to stderr and exits with an error status.
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n%s", os.Args[0], fmt.Sprintf(USAGE_MESSAGE, os.Args[0]))
	os.Exit(2)
}

// getDirectory returns the directory at the given path if it can be opened
// and confirmed to be a directory.
func getDirectory(directoryPath string) (*os.File, error) {
	directory, err := os.Open(directoryPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to open %s:", directoryPath)
	}
	if fileInfo, err := directory.Stat(); err != nil {
		return nil, errors.Wrapf(err, "Unable to access metadata for %s:", directory.Name())
	} else if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", directory.Name())
	}
	return directory, nil
}
