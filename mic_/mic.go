package main

import (
	"flag"
	"fmt"
	"mik/mic_/compiler_util"
	"mik/mic_/install"
	jsonconf "mik/mic_/json_conf"
	"mik/mic_/test"
	"os"
	"strings"
	"time"
)

var help_str string = ""
var test_var bool = false

func main() {
	if test_var {
		// O(n)
		test.TestLex(1)
		// wtf doesn't scale at all lol
		test.TestParse(1)
	}

	// load the source path from mik.conf
	/*src_path, err := os.ReadFile("../mik.conf")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(src_path))*/

	inptPtr := flag.String("i", "", "input flag")
	outpPtr := flag.String("o", "mik", "output flag")
	configPtr := flag.String("cnfg", "", "Path of a config file (mikConf.json)")
	cTarget := flag.Bool("c", true, "Specifies the compiling target.")
	asmTarget := flag.Bool("asm", false, "Sets target flag to asm")
	wasmTarget := flag.Bool("wasm", false, "Sets target to wasm")
	installPtr := flag.Bool("install", false, "wether the config should be created or not")
	ignorePtr := flag.String("ign", "", "Ignore a section put in \"\" and seperate by comma")

	flag.Parse()

	var ign_sections []string
	ign_sections = append(ign_sections, strings.Split(*ignorePtr, ":")...)

	var conf jsonconf.Config

	if *configPtr != "" {
		conf = jsonconf.Configure(*configPtr)
	} else {
		conf.CLibs = nil
	}

	// sets compiling target
	var tar string = "c"
	if *cTarget {
		tar = "c"
	} else if *asmTarget {
		tar = "asm"
	} else if *wasmTarget {
		tar = "wasm"
	}

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
		temp_txt, err = os.ReadFile(*inptPtr)
		if err != nil {
			panic(err)
		}

		// parse the memory address of temp_txt
		start := time.Now()
		txt = string(temp_txt)
		new_txt := compiler_util.Preprocess(&txt, inptPtr)
		tokens := compiler_util.Lex(new_txt, "lexer_test.mik")
		_, func_name, var_name := jsonconf.MakeCLib(conf, *inptPtr)

		ast := compiler_util.Parse(tokens, func_name, var_name, ign_sections)

		// generate code
		comp_success := compiler_util.Generate(&ast, tar, *outpPtr, *inptPtr, conf)
		if comp_success {
			compiler_util.NewSuccess("Succesfully compiled", fmt.Sprintf("In %s. To %s.c", time.Since(start).String(), *outpPtr), "", false)
		}
	}
}
