package main

import (
	"flag"
	"fmt"
	"mik/mic_/compiler_util"
	"mik/mip_/mip_util"
	"os"
)

func main() {
	// arg parse arguments
	installPtr := flag.String("install", "", "install a package from git")
	listPtr := flag.Bool("list", false, "list all packages")
	removePtr := flag.String("remove", "", "remove a package")
	addPkgPtr := flag.String("add_pkg", "", "adds a package at the specified location")
	initNew := flag.String("init", "", "initializes repo")
	milkPtr := flag.String("milk", "", "milks milk.pkg")
	addPtr := flag.String("add", "", "adds a git pkg to .pkgs")
	listcPtr := flag.Bool("listc", false, "list packages in current dir")

	flag.Parse()

	// Checks if any flag was specified
	if *installPtr != "" || *listPtr || *removePtr != "" || *addPkgPtr != "" || *initNew != "" || *milkPtr != "" || *addPtr != "" || *listcPtr {
		// Check which is active and if more than one is active
		if *installPtr != "" {
			// installs a specified github url
			mip_util.InstallGit(installPtr)
			os.Exit(0)
		} else if *listPtr {
			// lists all the installed packages
			mip_util.ListAll()
			os.Exit(0)
		} else if *removePtr != "" {
			// removes a specified package by name
			mip_util.RemovePkg(removePtr)
		} else if *addPkgPtr != "" {
			// adds a package (at specified path) to source
			ditch_nil := ""
			ditch_nil2 := "pkg/"
			ditch_nil3 := ""
			mip_util.AddPkgSrc(addPkgPtr, &ditch_nil, &ditch_nil2, &ditch_nil3)
			os.Exit(0)
		} else if *initNew != "" {
			// inits a package
			var dir string = ""
			if *initNew == "." {
				dir = ""
			} else {
				dir = *initNew + "/"
			}
			mip_util.Init(dir)
		} else if *milkPtr != "" {
			// milks a milk.pkg file
			var dir string = ""
			if *milkPtr == "." {
				dir = ""
			} else {
				dir = *milkPtr + "/"
			}
			mip_util.Milk(dir)
		} else if *addPtr != "" {
			mip_util.Add(addPtr)
		} else if *listcPtr {
			mip_util.ListCurrent()
		} else {
			// Too many args
			compiler_util.NewError("To many arguments specified", "", "", false)
		}
	} else {
		// No args specified
		fmt.Println("Usage:\n\t-install <github_link> to install a pkg from github\n\t-list to list all the installed packages\n\t-remove <pkg_name> to remove a pkg\n\t-add_pkg <dir_path> to add a dir to mik-srd")
	}
}
