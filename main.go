package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"path"
)

var red (func(...interface{}) string) = color.New(color.FgRed).SprintFunc()
var yellow (func(...interface{}) string) = color.New(color.FgYellow).SprintFunc()
var green (func(...interface{}) string) = color.New(color.FgGreen).SprintFunc()
var blue (func(...interface{}) string) = color.New(color.FgCyan).SprintFunc()

const NVIM_BUNDLE_NOT_FOUND string = "Unable to find bundle directory for neovim."
const VIM_BUNDLE_NOT_FOUND string = "Unable to find bundle directory for vim."

func main() {
	pluginDir := findNvimPluginDir()
	fmt.Println(blue(pluginDir) + " targeted")
	vimDir := findVimPluginDir()
	fmt.Println(blue(vimDir) + " alternatively targeted")
}

func findNvimPluginDir() string {
	xdgHome := os.Getenv("XDG_CONFIG_HOME")
	if len(xdgHome) < 1 {
		log.Println(yellow("Environment Variable XDG_CONFIG_HOME is not set" +
			" or has no value. Please set this if you" +
			" use neovim. Inferring default value for" +
			" XDG_CONFIG_HOME..."))
		xdgHome = path.Join(os.Getenv("HOME"), ".config")
	}

	fmt.Println("XDG_CONFIG_HOME=" + xdgHome)

	configPath := path.Join(xdgHome, "nvim")
	configDir, err := os.Open(configPath)
	if err != nil {
		log.Println(red(err))
		log.Println(red(NVIM_BUNDLE_NOT_FOUND))
		return ""
	}
	defer configDir.Close()
	fmt.Println(green(configPath + " exists"))
	bundlePath := path.Join(configPath, "bundle")
	bundleDir, err := os.Open(bundlePath)
	if err != nil {
		log.Println(red(err))
		log.Println(red(NVIM_BUNDLE_NOT_FOUND))
		return ""
	}
	defer bundleDir.Close()
	fmt.Println(green(bundlePath + " exists"))
	return bundlePath
}

func findVimPluginDir() string {
	vimConfig := os.Getenv("HOME")
	fmt.Println("HOME=" + vimConfig)

	configPath := path.Join(vimConfig, ".vim")
	configDir, err := os.Open(configPath)
	if err != nil {
		log.Println(red(err))
		log.Println(red(VIM_BUNDLE_NOT_FOUND))
		return ""
	}
	defer configDir.Close()
	fmt.Println(green(configPath + " exists"))
	bundlePath := path.Join(configPath, "bundle")
	bundleDir, err := os.Open(bundlePath)
	if err != nil {
		log.Println(red(err))
		log.Println(red(VIM_BUNDLE_NOT_FOUND))
		return ""
	}
	defer bundleDir.Close()
	fmt.Println(green(bundlePath + " exists"))
	return bundlePath
}
