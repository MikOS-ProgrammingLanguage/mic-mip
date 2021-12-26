package mip_util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"mik/mic_/compiler_util"
	"os"
	"strconv"
	"strings"

	cp "github.com/otiai10/copy"
)

var src_path string = strings.ReplaceAll(strings.ReplaceAll(readConf(), "\\", "/"), "\n", "")
var mik_src_path string = src_path + "mik-src/"
var req_satisfied_loc string = mik_src_path + "req_satisfied.conf"
var TEMP_CNT int = -1

func strip(str string) string {
	return strings.Join(strings.Fields(str), "")
}
func mk_str(arr []string) string {
	ret_str := ""
	for i, val := range arr {
		if i < 1 {
			continue
		}
		ret_str += val
	}
	return ret_str
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func readConf() string {
	cntnt, err := os.ReadFile("mik.conf")
	if err != nil {
		panic(err)
	}
	return string(cntnt)
}

// Lists all the packages in mik-src
func ListAll() {
	// Loads the content of the req_satisfied
	content, err := os.ReadFile(req_satisfied_loc)
	if err != nil {
		panic(err)
	}

	split := strings.Split(string(content), "\n")

	if len(split) == 1 && split[0] == "" {
		compiler_util.NewInfo("No packages found", "", "", true)
	} else {
		if len(split) == 1 {
			fmt.Println("You currently have 1 package installed")
		} else {
			fmt.Println("You currently have " + strconv.Itoa(len(split)) + " packages installed")
		}
		for i := 0; i < len(split); i++ {
			fmt.Println("\t|____ " + strings.Split(split[i], ":::")[1])
		}
	}
}

// add_pkg adds a pkg to mik-src with a specified path
func AddPkg(path *string) {
	checkErr(os.RemoveAll(fmt.Sprintf("%stemp/", mik_src_path)))
	checkErr(os.Mkdir(fmt.Sprintf("%stemp/", mik_src_path), os.ModePerm))
	TEMP_CNT += 1
	var PKG_NAME string
	var CURRENT_TEMP string = fmt.Sprintf("tmp%d", TEMP_CNT)
	var IGNORE_FILES []string
	var DEPENDENCIES []string

	*path = strings.ReplaceAll(strings.ReplaceAll(*path, "\\", "/"), "\n", "")

	// copy dir to to temp
	checkErr(cp.Copy(*path, fmt.Sprintf("%stemp/%s", mik_src_path, CURRENT_TEMP)))
	out, err := os.ReadFile(fmt.Sprintf("%stemp/%s/milk.pkg", mik_src_path, CURRENT_TEMP))
	if err != nil {
		compiler_util.NewCritical("No milk.pkg found", "At "+fmt.Sprintf("%s", *path), "", true)
	} else {
		pkg_txt := string(out)
		pkg_args := strings.Split(pkg_txt, "\n")
		for _, val := range pkg_args {
			val = strip(val)
			val_2 := strings.Split(val, ":")

			if len(val_2) < 2 {
				compiler_util.NewError("No vars specified after <arg>:", "", "", true)
			}
			switch val_2[0] {
			case "package-name":
				PKG_NAME = mk_str(val_2)
			case "ignore-file":
				IGNORE_FILES = append(IGNORE_FILES, mk_str(val_2))
			case "depends":
				DEPENDENCIES = append(DEPENDENCIES, mk_str(val_2))
			default:
				compiler_util.NewError("Invalid Argument: "+val_2[0], "", "", true)
			}
		}
		// remove all files in ignore
		for _, val := range IGNORE_FILES {
			checkErr(os.Remove(fmt.Sprintf("%stemp/%s/%s", mik_src_path, CURRENT_TEMP, val)))
		}
		compiler_util.NewInfo("Done Ignoring", "", "", false)

		// make dependencies
		for _, val := range DEPENDENCIES {
			InstallGit(&val)
		}

		// make a pkg structure
		checkErr(os.MkdirAll(mik_src_path+"pkg/"+PKG_NAME, os.ModePerm))
		checkErr(cp.Copy(fmt.Sprintf("%stemp/%s", mik_src_path, CURRENT_TEMP), fmt.Sprintf("%spkg/%s", mik_src_path, PKG_NAME)))

		// clear temp
		checkErr(os.RemoveAll(fmt.Sprintf("%stemp/%s/", mik_src_path, CURRENT_TEMP)))

		// preprocess all not ignored files ending in .milk
		files_, err := ioutil.ReadDir(fmt.Sprintf("%spkg/%s/", mik_src_path, PKG_NAME))
		checkErr(err)
		var files []string

		// get all files that end in .milk
		for _, val := range files_ {
			name := val.Name()
			if strings.HasSuffix(name, ".milk") {
				files = append(files, name)
			}
		}

		// make a yoink string
		var yoink_str string
		for _, val := range files {
			yoink_str += fmt.Sprintf("#yoink <%s>\n", val)
		}

		// preprocess yoink string and write it
		pth := fmt.Sprintf("%spkg/%s/", mik_src_path, PKG_NAME)
		preprocessed_txt := compiler_util.Preprocess(&yoink_str, &pth)
		checkErr(os.WriteFile(fmt.Sprintf("%spkg/%s/main_%s.milk", mik_src_path, PKG_NAME, PKG_NAME), []byte(*preprocessed_txt), 0644))
		checkErr(os.WriteFile(req_satisfied_loc, []byte(PKG_NAME+":::"+PKG_NAME), os.ModeAppend))

		compiler_util.NewSuccess(fmt.Sprintf("Sucessfully added package: %s\nYou can now use it with '#yoink-src<%s>", PKG_NAME, PKG_NAME), "", "", false)
	}
}

// install installs to mik-src via. a github link
func InstallGit(url *string) {
	// clear git

}

// remove removes a pkg from mik-src
func RemovePkg(pkg_name *string) {
	if *pkg_name == "." {
		compiler_util.NewWarning("Do you want to remove all packages? (Y/N)", "", "", false)
		reader := bufio.NewReader(os.Stdin)

		answer, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}

		if string(answer) == "Y" || string(answer) == "y" {
			compiler_util.NewInfo("Deleting all packages...", "", "", false)
			// delete all pkgs
		} else {
			compiler_util.NewInfo("Aborted.", "", "", false)
		}
	} else {
		// check if pkg is in req_satisfied.conf
		cntnt, err := os.ReadFile(req_satisfied_loc)
		checkErr(err)

		var pkg_exists bool = false
		var new_cntnt string = ""

		for _, val := range strings.Split(string(cntnt), "\n") {
			if *pkg_name == strings.Split(string(val), ":::")[1] {
				pkg_exists = true
				new_cntnt = strings.ReplaceAll(string(cntnt), *pkg_name+":::"+strings.Split(string(val), ":::")[1], "")
				break
			}
		}

		if pkg_exists {
			checkErr(os.RemoveAll(mik_src_path + "pkg/" + *pkg_name))
			checkErr(os.WriteFile(req_satisfied_loc, []byte(new_cntnt), 0644))
			compiler_util.NewSuccess("Succesfully removed package: "+*pkg_name, "", "", true)
		} else {
			compiler_util.NewCritical("The package you tried to remove was not found", "Try to use -list to see all installed packages", "", true)
		}
	}
}
