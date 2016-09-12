package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"strings"
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
	numberPlugins := len(plugins)
	fmt.Printf(green("%d plugins found\n"), numberPlugins)
	if numberPlugins < 1 {
		fmt.Println(yellow("No plugins found; nothing to update"))
		os.Exit(0)
	}
	updatesCompleted := make(chan int, 0)
	for _, repository := range plugins {
		go updateGitRepo(repository, updatesCompleted)
	}
	for i := 0; i < numberPlugins; i++ {
		<-updatesCompleted
	}
	fmt.Println(green("Plugins up to date."))

}

// updateGitRepo runs the command `git pull` within the given directory
// and signals that it is done by sending a value on the done channel.
func updateGitRepo(path string, done chan int) {
	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, yellow("Unable to change directories to %s"), path))
		done <- 0
		return
	}
	// Update remote branches
	gitUpdate := exec.Command("git", "remote", "update", "origin")
	err = gitUpdate.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, yellow("Unable to run `git remote update origin` in directory %s"), path))
		done <- 0
		return
	}
	// Check whether current branch is out of date with upstream
	gitStatus := exec.Command("git", "status")
	currentStatus, err := gitStatus.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, yellow("Unable to get stdout from `git status` in directory %s"), path))
		done <- 0
		return
	}
	currentStatusString := string(currentStatus[:])
	if !strings.Contains(currentStatusString, "behind") {
		fmt.Printf(blue("Nothing to update for %s\n"), path)
		done <- 0
		return

	}
	// Pull if out of date
	gitPull := exec.Command("git", "pull")
	err = gitPull.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, yellow("Unable to run `git pull` in directory %s"), path))
		done <- 0
		return
	}
	fmt.Printf(blue("%s updated\n"), path)
	done <- 0
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
