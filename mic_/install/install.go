package install

import (
	"mik/mic_/compiler_util"
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
	err2 := os.Chdir("/etc")
	if err2 != nil {
		panic(err2)
	}
	err3 := os.WriteFile("mik.conf", data, 0644)
	if err3 != nil {
		panic(err3)
	}
	compiler_util.NewSuccess("Updated your path succesfully", "", "", false)
	compiler_util.NewInfo("Your new Path is now :", string(data), "", false)
}
