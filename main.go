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

	// Find plugin git repositories
	pluginCandidates, err := pluginDir.Readdir(0)
	if err != nil {
		panic(errors.Wrapf(err, "Unable to read the contents of plugin directory %s", pluginDir.Name()))
	}
	plugins := make([]string, 0)
	for _, candidate := range pluginCandidates {
		fullPath := pluginDir.Name() + string(os.PathSeparator) + candidate.Name()
		fmt.Printf("Checking %15s...", candidate.Name())
		if ok, err := isGitRepository(fullPath); err != nil {
			fmt.Fprintln(os.Stderr, yellow(err))
		} else if ok {
			plugins = append(plugins, fullPath)
			fmt.Println(green("OK"))
		} else {
			fmt.Println(red("NO"))
		}
	}
	fmt.Printf(green("%d plugins found\n"), len(plugins))
	numberPlugins := len(plugins)
	updatesCompleted := make(chan int, 0)
	for _, repository := range plugins {
		go updateGitRepo(repository, updatesCompleted)
	}
	for i := 0; i < numberPlugins; i++ {
		<-updatesCompleted
	}
	fmt.Println(green("Plugins updated."))

}

// isGitRepository checks whether a path points to the root directory of a git
// repository.
func isGitRepository(path string) (bool, error) {
	candidateDir, err := getDirectory(path)
	defer candidateDir.Close()
	if err != nil {
		return false, errors.Wrapf(err, "Unable to determine if %s is a git repository", path)
	}
	contents, err := candidateDir.Readdir(0)
	if err != nil {
		return false, errors.Wrapf(err, "Unable to determine if %s is a git repository", path)
	}
	for _, file := range contents {
		if file.Name() == ".git" {
			return true, nil
		}
	}
	return false, nil

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
