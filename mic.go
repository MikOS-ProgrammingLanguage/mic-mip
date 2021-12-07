package main

import (
	"flag"
	"fmt"
	"mic/compiler_util"
	"mic/install"
	"os"
)

var help_str string = ""

func main() {
	// load the source path from mik.conf
	src_path, err := os.ReadFile("../mik.conf")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(src_path))

	inptPtr := flag.String("i", "", "input flag")
	outpPtr := flag.String("o", "mik", "output flag")
	installPtr := flag.Bool("install", false, "wether the config should be created or not")

	flag.Parse()

	if *inptPtr == "" && *installPtr {
		install.Install_at_current_path()
	} else if *inptPtr == "" && *outpPtr == "mik" && !*installPtr {
		fmt.Println(help_str)
	} else {
		// start parsing in the generator
		var txt string
		var temp_txt []byte
		var err error

		// load source_code
		temp_txt, err = os.ReadFile("test/main_test.mik")
		if err != nil {
			panic(err)
		}

		// parse the memory address of temp_txt
		txt = string(temp_txt)
		_ = compiler_util.Lex(&txt)
	}
}
