package main

import (
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
	pluginDir, err := findNvimPluginDir()
	if err != nil {
		log.Println(red(err))
	}
	fmt.Println(blue(pluginDir) + " targeted")
	vimDir, err := findVimPluginDir()
	if err != nil {
		log.Println(red(err))
	}
	fmt.Println(blue(vimDir) + " targeted")
}

// findNvimPluginDir returns the directory in which neovim stores its plugins
// or the empty string if it was unable to find the directory
func findNvimPluginDir() (string, error) {
	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if len(xdgHome) < 1 {
		fmt.Println(yellow("Environment Variable XDG_CONFIG_HOME" +
			" is not set or has no value. Please set this if" +
			" you use neovim. Inferring default value for" +
			" XDG_CONFIG_HOME..."))
		xdgHome = path.Join(os.Getenv("HOME"), ".config")
	}

	fmt.Println("XDG_CONFIG_HOME=" + xdgHome)
	bundlePath := path.Join(xdgHome, "nvim", "bundle")
	bundleDir, err := os.Open(bundlePath)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to find %s directory:", bundlePath)
	}
	defer bundleDir.Close()
	fmt.Println(green(bundlePath + " exists"))
	return bundlePath, nil
}

// findNvimPluginDir returns the directory in which neovim stores its plugins
// or the empty string if it was unable to find the directory
func findVimPluginDir() (string, error) {
	home := os.Getenv("HOME")
	fmt.Println("HOME=" + home)

	bundlePath := path.Join(home, ".vim", "bundle")
	bundleDir, err := os.Open(bundlePath)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to find %s directory:", bundlePath)
	}
	defer bundleDir.Close()
	fmt.Println(green(bundlePath + " exists"))
	return bundlePath, nil
}
