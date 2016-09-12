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

const usageMessage string = `
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
	flag.Parse()

	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
	}
	pluginPath := os.Args[1]

	// Find plugins
	pluginDir, err := getDirectory(pluginPath)
	defer pluginDir.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, red(err))
		os.Exit(1)
	}
	fmt.Println(blue(pluginDir.Name()) + " targeted")
}

// printUsage prints usage information to stderr and exits with an error status.
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n%s", os.Args[0], fmt.Sprintf(usageMessage, os.Args[0]))
	os.Exit(2)
}

// isDir returns whether a given path points to a directory. If checking this information
// produces an error, that is returned as well.
func isDirectory(path string) (bool, error) {
	stats, err := os.Stat(path)
	if err != nil {
		return false, errors.Wrapf(err, "Could not stat %s to check if directory", path)
	}
	return stats.IsDir(), nil
}

// getDirectory returns the directory at the given path if it can be opened
// and confirmed to be a directory.
func getDirectory(directoryPath string) (*os.File, error) {
	if ok, err := isDirectory(directoryPath); err != nil {
		return nil, errors.Wrapf(err, "Unable to check whether %s is a directory", directoryPath)
	} else if !ok {
		return nil, fmt.Errorf("%s is not a directory", directoryPath)
	}
	// if you get here, it's a directory
	directory, err := os.Open(directoryPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to open %s:", directoryPath)
	}
	return directory, nil
}
