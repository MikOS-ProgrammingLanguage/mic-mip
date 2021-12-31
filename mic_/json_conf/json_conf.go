package jsonconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	CLibs []Lib `json:"CLibs"`
}
type Lib struct {
	Dir       string   `json:"dir"`
	Name      string   `json:"name"`
	Pth       string   `json:"path"`
	FuncNames []string `json:"FuncNames"`
	VarNames  []string `json:"VarNames"`
}

func MakeCLib(conf_s Config, inpt string) (*string, []string, []string) {
	new_inpt_pth := ""
	for i, val := range inpt {
		if i < len(inpt)-1 {
			new_inpt_pth += string(val)
		}
	}

	ditch := ""
	var code_ *string = &ditch
	var func_names []string
	var var_names []string
	*code_ += "// LIB CODE\n"
	for _, val := range conf_s.CLibs {
		var dir_flg bool = false
		var name string
		var pth string
		if val.Dir == "true" {
			dir_flg = true
		}
		name = val.Name
		pth = val.Pth
		func_names = append(func_names, val.FuncNames...)
		var_names = append(var_names, val.VarNames...)

		// walk the path and add the code
		if !dir_flg {
			if strings.HasSuffix(pth, ".c") || strings.HasSuffix(pth, ".h") {
				data, err := os.ReadFile(new_inpt_pth + "/" + pth)
				if err != nil {
					fmt.Println("LibFileNotFoundError! ", "The lib file "+name+" was not found")
					os.Exit(1)
				}

				*code_ += string(data)
			} else {
				fmt.Println("IllegalLibError! ", "The specified lib is not .c or .h")
				os.Exit(1)
			}
		} else {
			file, err := os.ReadDir(new_inpt_pth + pth)
			if err != nil {
				fmt.Println("IllegalLibError! ", "The specified lib is no directory")
				os.Exit(1)
			}
			// make lib from dir

			for _, d := range file {
				if strings.HasSuffix(d.Name(), ".c") || strings.HasSuffix(d.Name(), ".h") {
					// read file contents
					data, err := ioutil.ReadFile(new_inpt_pth + pth + d.Name())
					if err != nil {
						fmt.Println("LibFileNotFoundError! ", "A file in the lib folder could not be opened")
						os.Exit(1)
					}

					*code_ += string(data)
				}
			}
		}
	}
	*code_ += "\n// LIB CODE END\n"
	return code_, func_names, var_names
}

func Configure(fpath string) Config {
	var conf Config

	jsonFile, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	byte_, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byte_, &conf)
	return conf
}
