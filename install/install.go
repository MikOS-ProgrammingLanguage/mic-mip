package install

import (
	"mic/compiler_util"
	"os"
	"strings"
)

// Creates a config at the current location
func Install_at_current_path() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	data := []byte(strings.ReplaceAll(path, "\\", "/"))
	split := strings.Split(string(data), "/")
	rem_split := split[0 : len(split)-1]
	new_string := ""
	// puts the strings in the array back together
	for i := 0; i < len(rem_split); i++ {
		new_string += rem_split[i]
		new_string += "/"
	}
	os.WriteFile("../mik.conf", []byte(new_string), 0644)
	compiler_util.NewSuccess("Updated your path succesfully", "", "", false)
	compiler_util.NewInfo("Your new Path is now :", new_string, "", false)
}
