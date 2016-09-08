package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"log"
	"os"
)

var red (func(...interface{}) string) = color.New(color.FgRed).SprintFunc()
var yellow (func(...interface{}) string) = color.New(color.FgYellow).SprintFunc()
var green (func(...interface{}) string) = color.New(color.FgGreen).SprintFunc()
var blue (func(...interface{}) string) = color.New(color.FgCyan).SprintFunc()

const MISSING_FLAG_ERROR string = `The plugin-path argument is required.
Usage:
go-epidemic --plugin-path=$HOME/.config/nvim/bundle
  or
go-epidemic --plugin-path=$HOME/.vim/bundle
`

func main() {
	// Remove timestamp prefix
	log.SetFlags(0)

	// Parse command line arguments
	var pluginPath string
	flag.StringVar(&pluginPath, "plugin-path", "", "REQUIRED: the path to your pathogen plugins")
	flag.Parse()

	// Handle args
	if pluginPath == "" {
		log.Fatalln(red(MISSING_FLAG_ERROR))
	}

	// Find plugins
	pluginDir, err := getDirectory(pluginPath)
	if err != nil {
		log.Fatalln(red(err))
	}
	fmt.Println(blue(pluginDir.Name()) + " targeted")
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
