package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"log"
	"os"
	"path"
)

var red (func(...interface{}) string) = color.New(color.FgRed).SprintFunc()
var yellow (func(...interface{}) string) = color.New(color.FgYellow).SprintFunc()
var green (func(...interface{}) string) = color.New(color.FgGreen).SprintFunc()
var blue (func(...interface{}) string) = color.New(color.FgCyan).SprintFunc()

func main() {
	// Parse command line arguments
	var vimVersion string
	var pluginPath string
	var pluginDirectoryFunction (func() (*os.File, error))
	flag.StringVar(&vimVersion, "vim-version", "neovim", "values are \"vim\" or \"neovim\"")
	flag.StringVar(&pluginPath, "plugin-path", "", "the path to your pathogen plugins")
	flag.Parse()

	// Handle args
	switch vimVersion {
	case "vim":
		pluginDirectoryFunction = findVimPluginDir
	case "neovim":
		pluginDirectoryFunction = findNvimPluginDir
	}
	if pluginPath != "" {
		pluginDirectoryFunction = func() (*os.File, error) {
			return getDirectory(pluginPath)
		}
	}

	// Find plugins
	pluginDir, err := pluginDirectoryFunction()
	if err != nil {
		log.Println(red(err))
	}
	fmt.Println(blue(pluginDir.Name()) + " targeted")
}

// getDirectory returns the directory at the given path if it can be opened
// and confirmed to be a directory.
func getDirectory(directoryPath string) (*os.File, error) {
	directory, err := os.Open(directoryPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to open %s:", directory)
	}
	if fileInfo, err := directory.Stat(); err != nil {
		return nil, errors.Wrapf(err, "Unable to access metadata for %s:", directory)
	} else if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", directory)
	}
	return directory, nil
}

// findNvimPluginDir returns the directory in which neovim stores its plugins
// or the empty string if it was unable to find the directory.
func findNvimPluginDir() (*os.File, error) {
	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if len(xdgHome) < 1 {
		fmt.Println(yellow("$XDG_CONFIG_HOME undefined. Inferring default value..."))
		xdgHome = path.Join(os.Getenv("HOME"), ".config")
	}

	bundlePath := path.Join(xdgHome, "nvim", "bundle")
	return getDirectory(bundlePath)
}

// findNvimPluginDir returns the directory in which neovim stores its plugins
// or the empty string if it was unable to find the directory.
func findVimPluginDir() (*os.File, error) {
	home := os.Getenv("HOME")
	bundlePath := path.Join(home, ".vim", "bundle")
	return getDirectory(bundlePath)
}
