package compiler_util

import (
	"fmt"
	"os"
)

type colors struct {
	colorReset  string
	colorRed    string
	colorGreen  string
	colorYellow string
	colorBlue   string
	colorPurple string
	colorCyan   string
	colorWhite  string
}

func makeColorPalette() colors {
	color := colors{"\033[0m", "\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
	return color
}

/*NewError takes the error type, it's description,
it's position, and if it should quit. It then prints the error and quits if told to*/
func NewError(err_type, description, position string, quit bool) {
	c_palette := makeColorPalette()
	fmt.Println(string(c_palette.colorRed), "\r[ERROR]"+err_type+"! "+description+" At "+position, string(c_palette.colorReset))
	if quit {
		os.Exit(0)
	}
}

/*NewWarning takes the warning type, it's description,
it's position, and if it should quit. It then prints the error and quits if told to*/
func NewWarning(err_type, description, position string, quit bool) {
	c_palette := makeColorPalette()
	fmt.Println(string(c_palette.colorPurple), "\r[WARNING]"+err_type+" "+description+" At "+position, string(c_palette.colorReset))
	if quit {
		os.Exit(0)
	}
}

/*NewCritical takes the critical type, it's description,
it's position, and if it should quit. It then prints the error and quits if told to*/
func NewCritical(err_type, description, position string, quit bool) {
	c_palette := makeColorPalette()
	fmt.Println(string(c_palette.colorBlue), "\r[CRITICAL]"+err_type+"! "+description+" At "+position, string(c_palette.colorReset))
	if quit {
		os.Exit(0)
	}
}

/*NewInfo takes the info type, it's description,
it's position, and if it should quit. It then prints the error and quits if told to*/
func NewInfo(err_type, description, position string, quit bool) {
	c_palette := makeColorPalette()
	fmt.Println(string(c_palette.colorYellow), "\r[INFO]"+err_type+" "+description+" At "+position, string(c_palette.colorReset))
	if quit {
		os.Exit(0)
	}
}

/*NewSuccess takes the success type, it's description,
it's position, and if it should quit. It then prints the error and quits if told to*/
func NewSuccess(err_type, description, position string, quit bool) {
	c_palette := makeColorPalette()
	fmt.Println(string(c_palette.colorGreen), "\r[SUCCESS]"+err_type+"! "+description+" At "+position, string(c_palette.colorReset))
	if quit {
		os.Exit(0)
	}
}
