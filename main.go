package main

import (
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"path"
)

var yellow (func(...interface{}) string) = color.New(color.FgYellow).SprintFunc()
var green (func(...interface{}) string) = color.New(color.FgGreen).SprintFunc()

func main() {
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
		log.Fatalln(err)
	}
	defer configDir.Close()
	fmt.Println(green(configPath + " exists"))
	bundlePath := path.Join(configPath, "bundle")
	bundleDir, err := os.Open(bundlePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer bundleDir.Close()
	fmt.Println(green(bundlePath + " exists"))
}
