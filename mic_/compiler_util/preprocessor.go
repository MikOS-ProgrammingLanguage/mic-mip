package compiler_util

import (
	"fmt"
	"os"
	"strings"
)

var Yoinked_files []string
var Yoinked_files_src []string

func Preprocess(text, pth *string) *string {
	*pth = strings.ReplaceAll(*pth, "\\", "/")
	new_pth := *pth
	temp_pth := strings.Split(*pth, "/")
	*pth = ""
	for i, val := range temp_pth {
		if i < len(temp_pth)-1 {
			*pth += val + "/"
		}
	}
	*pth = strings.ReplaceAll(*pth, "\n", "")
	mip_src_path_, err := os.ReadFile("mik.conf")
	if err != nil {
		panic(err)
	}
	mip_src_path := strings.ReplaceAll(string(mip_src_path_), "\n", "")

	var ditch_nil string = ""
	var new_txt *string = &ditch_nil

	*new_txt += "@section(\"" + new_pth + "\")"
	for _, val := range strings.Split(*text, "\n") {
		if (strings.HasPrefix(val, "#yoink <") || strings.HasPrefix(val, "#yoink-src <")) && strings.HasSuffix(val, ">") {
			var yoink_name string = ""
			if strings.HasPrefix(val, "#yoink <") {
				yoink_name = strings.Split(strings.Split(val, "#yoink <")[1], ">")[0]
				if StringInSlice(yoink_name, Yoinked_files) {
					*new_txt += "\n"
					continue
				}
				// add a section, append the data to the text and close the section
				*new_txt += fmt.Sprintf("\n@section(\"%s\")\n", yoink_name)
				data, err := os.ReadFile(*pth + yoink_name)
				if err != nil {
					panic(err)
				}
				*new_txt += string(data)
				*new_txt += "\n@secend"

				Yoinked_files = append(Yoinked_files, yoink_name)
			} else if strings.HasPrefix(val, "#yoink-src <") {
				yoink_name = strings.Split(strings.Split(val, "#yoink-src <")[1], ">")[0]
				if StringInSlice(yoink_name, Yoinked_files_src) {
					*new_txt += "\n"
				}
				// add a section, append the data to the text and close the section
				*new_txt += fmt.Sprintf("@section(\"%s\")\n", yoink_name)
				data, err := os.ReadFile(string(mip_src_path) + "mik-src/pkg/" + yoink_name + "/main_" + yoink_name + ".milk")
				if err != nil {
					panic(err)
				}
				*new_txt += string(data)
				*new_txt += "\n@secend"

				Yoinked_files_src = append(Yoinked_files_src, yoink_name)
			} else {
				NewError("", "", "", true)
			}

		} else {
			*new_txt += val + "\n"
		}
	}
	*new_txt += "@secend"
	return new_txt
}

// nachher sollte gier noch der mip_src_path_hin
