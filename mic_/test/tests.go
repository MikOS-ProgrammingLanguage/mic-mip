package test

import (
	"fmt"
	"math/rand"
	"mik/mic_/compiler_util"
	"time"
)

//
func TestLex(target int) {
	testLexID(target)
	testLexSpecial(target)
	fmt.Println()
	compiler_util.NewSuccess(fmt.Sprintf("Done: Lexing... Used a string of len: %d each", target), "", "", false)
	// put in lex_test.dump
}

//
func testLexID(target int) {
	var test_str_ids []byte = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789\n ")
	var id_test string = ""
	var id_iterations int = 0
	var id_target int = target
	var id_success bool = false
	var id_lex_success bool = false

	go func() {
		fmt.Printf("\r                                ")
		for !id_success {
			fmt.Printf("\r\\  [%d of %d]", id_iterations, id_target)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r|  [%d of %d]", id_iterations, id_target)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r/  [%d of %d]", id_iterations, id_target)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r-  [%d of %d]", id_iterations, id_target)
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Printf("\r✅  [%d of %d]", id_target, id_target)
	}()

	compiler_util.NewInfo("Lexer Test #1", "Testing ids...⭕️", "", false)
	for ; id_iterations < id_target; id_iterations++ {
		rnum := rand.Intn(len(test_str_ids))
		id_test += string(test_str_ids[rnum])
	}
	fmt.Println()
	id_success = true

	go func() {
		fmt.Printf("\r                                ")
		for !id_lex_success {
			fmt.Printf("\r\\  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r|  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r/  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r-  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Printf("\r                                ")
	}()
	start := time.Now()
	_ = compiler_util.Lex(&id_test, "")
	id_lex_success = true
	compiler_util.NewSuccess("Lexer Test #1✅", fmt.Sprintf("It took %s", time.Since(start)), "", false)
}

//
func testLexSpecial(target int) {
	var test_str_ids []byte = []byte("<>=!&|;{}()[].+-*/?%\n ")
	var id_test string = ""
	var id_iterations int = 0
	var id_target int = target
	var id_success bool = false
	var id_lex_success bool = false

	go func() {
		fmt.Printf("\r                                ")
		for !id_success {
			fmt.Printf("\r\\  [%d of %d]", id_iterations, id_target)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r|  [%d of %d]", id_iterations, id_target)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r/  [%d of %d]", id_iterations, id_target)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r-  [%d of %d]", id_iterations, id_target)
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Printf("\r✅  [%d of %d]", id_target, id_target)
	}()

	compiler_util.NewInfo("Lexer Test #2", "Testing special Tokens...⭕️", "", false)
	for ; id_iterations < id_target; id_iterations++ {
		rnum := rand.Intn(len(test_str_ids))
		id_test += string(test_str_ids[rnum])
	}
	id_success = true

	go func() {
		fmt.Printf("\r                                ")
		for !id_lex_success {
			fmt.Printf("\r\\  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r|  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r/  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r-  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Printf("\r                                ")
	}()
	start := time.Now()
	_ = compiler_util.Lex(&id_test, "")
	id_lex_success = true
	compiler_util.NewSuccess("Lexer Test #2✅", fmt.Sprintf("It took %s", time.Since(start)), "", false)
}

//
func TestParse(target int) {
	testParseExpr(target)
	fmt.Println()
	compiler_util.NewSuccess(fmt.Sprintf("Done: Parsing... Used a string of len %d each", target), "", "", false)
	// put in parse_text.dump
}

//
func testParseExpr(target int) {
	var test_str_ids []byte = []byte("+-*/0123456789 ")
	var expr_test string = ""
	var expr_iterations int = 0
	var expr_target int = target
	var expr_success bool = false
	var expr_parse_success bool = false

	go func() {
		fmt.Printf("\r                                ")
		for !expr_success {
			fmt.Printf("\r\\  [%d of %d]", expr_iterations, expr_target)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r|  [%d of %d]", expr_iterations, expr_target)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r/  [%d of %d]", expr_iterations, expr_target)
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r-  [%d of %d]", expr_iterations, expr_target)
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Printf("\r✅  [%d of %d]", expr_target, expr_target)
	}()

	compiler_util.NewInfo("Parser Test #1", "Testing expr...⭕️", "", false)
	for ; expr_iterations < expr_target; expr_iterations++ {
		rnum := rand.Intn(len(test_str_ids))
		expr_test += string(test_str_ids[rnum])
	}
	fmt.Println()
	expr_success = true

	go func() {
		fmt.Printf("\r                                ")
		for !expr_parse_success {
			fmt.Printf("\r\\  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r|  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r/  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
			fmt.Printf("\r-  %s", "Lexing...")
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Printf("\r                                ")
	}()
	var test string = "10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50+10+20*30*50"
	a := compiler_util.Lex(&test, "")
	illegal_name := []string{"", ""}
	illegal_name2 := []string{"", ""}
	start := time.Now()
	_ = compiler_util.Parse(a, illegal_name, illegal_name2, []string{""})
	expr_parse_success = true
	compiler_util.NewSuccess("Parser Test #1✅", fmt.Sprintf("It took %s", time.Since(start)), "", false)
}
