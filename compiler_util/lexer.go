package compiler_util

import (
	"fmt"
	"strings"
)

// all legal characters for func_names and stuff (numbers also work but that's done later on)
var chars string = "abcdefghijklmnopqrstuvwxyzäöüABCDEFGHIJKLMNOPQRSTUVWXYTÄÖÜ_"
var numbers string = "0123456789"

// All the tokens that exist
const (
	TT_LTHEN    string = "LTHEN"    // <
	TT_LEQ      string = "LEQ"      // <=
	TT_GTHEN    string = "GTHEN"    // >
	TT_GEQ      string = "GEQ"      // >=
	TT_KAND     string = "KAND"     // &
	TT_AND      string = "AND"      // &&
	TT_OR       string = "OR"       // ||
	TT_NOT      string = "NOT"      // !
	TT_NEQ      string = "NEQ"      // !=
	TT_DOT      string = "DOT"      // .
	TT_PLUS     string = "PLS"      // +
	TT_PLUSEQ   string = "PLSEQ"    // +=
	TT_MINUS    string = "MIN"      // -
	TT_ARROW    string = "ARROW"    // ->
	TT_MINUSEQ  string = "MINEQ"    // -=
	TT_MUL      string = "MUL"      // *
	TT_MULEQ    string = "MULEQ"    // *=
	TT_DIV      string = "DIV"      // /
	TT_DIVEQ    string = "DIVEQ"    // /=
	TT_ASSGN    string = "ASSGN"    // =
	TT_EQ       string = "EQ"       // ==
	TT_REASSGN  string = "REASSGN"  // ?=
	TT_RPAREN   string = "RPAREN"   // (
	TT_LPAREN   string = "LPAREN"   // )
	TT_LBRK     string = "LBRK"     // [
	TT_RBRK     string = "RBRK"     // ]
	TT_LCURL    string = "LCURL"    // {
	TT_RCURL    string = "RCURL"    // }
	TT_COMMA    string = "COMMA"    // ,
	TT_SEMIC    string = "SEMIC"    // ;
	TT_PERCENT  string = "PERCENT"  // %
	TT_INT      string = "INT"      // numbers
	TT_FLOAT    string = "FLT"      // numbers.numbers
	TT_ID       string = "ID"       // a string of chars for example for a var name
	TT_STRING   string = "STR"      // "string goes here"
	TT_CHAR     string = "CHAR"     // 'single_character'
	TT_DEBUG    string = "DEBUG"    // @debug
	TT_OVERRIDE string = "OVERRIDE" // @overwrite
	TT_EOF      string = "EOF"      // end of file
)

// is used as a part of the position and only exists because you can't have a struct in itself
type primitive_position struct {
	section  string
	ln_count int
}

// is used to keep track of the current position in files and to find the previous section once one section is done
type position struct {
	prev_section primitive_position
}

// token is what is given to the parser for each token
type token struct {
	section  string
	ln_count int
	type_    string
	value    string
}

var text []byte
var text_position int = -1
var current_char byte = ' '
var is_eof bool = false
var section position = position{primitive_position{ln_count: 1, section: ""}}
var expect_mikas bool = false
var in_mikas bool = false

func currentPos() string {
	return string(section.prev_section.section) + " at ln " + string(section.prev_section.ln_count)
}

// Lexer Functions:
// advances forward once
func advance() {
	text_position += 1
	if text_position < len(text) {
		current_char = text[text_position]
	} else {
		// assigns the current char to ' ' because that's something that never happens as we ignore white spaces.
		current_char = ' '
		is_eof = true
	}
}

// creates a token with the next characters and returns it
func make_str() token {
	var mk_str string

	if in_mikas {
		cnt := 0
		// loops as long as not end of file and aslong as the mikas hasn't ended
		for !is_eof && current_char != '}' {
			if current_char == '\n' && cnt > 0 { // checks for a newline and if it's more than one character appended bcs else the c command would be pretty trashy
				mk_str += ";"
			} else if current_char == '\n' && cnt <= 0 {
				// don't do anything for a newline at the start of the function
			} else if current_char == '\t' {
				// don't do anything for tabs
			} else {
				mk_str += string(current_char) // append every other character
			}
			advance()
			cnt++
		}
	} else {
		advance()
		for !is_eof && current_char != '"' {
			mk_str += string(current_char)
			advance()
		}
	}
	return token{section.prev_section.section, section.prev_section.ln_count, TT_STRING, mk_str}
}

// creates a token for the following characters
func make_id() token {
	var id_str string

	// loops over the next chars and as long as the char is in the alphabet and all numbers it appends
	for !is_eof && strings.Contains(chars+numbers, string(current_char)) {
		id_str += string(current_char)
		advance()
	}
	// sets expect mikas to true if the id is mikas so the compiler can process inline assembly correctly
	if id_str == "mikas" {
		expect_mikas = true
	}
	return token{section.prev_section.section, section.prev_section.ln_count, TT_ID, id_str}
}

// creates a token for either a int or float based on if it has a decimal point
func make_number() token {
	var num_str string
	dot_cnt := 0

	for !is_eof && strings.Contains(numbers+".", string(current_char)) {
		if current_char == '.' {
			if dot_cnt == 1 {
				break
			}
			dot_cnt++
			num_str += "."
		} else {
			num_str += string(current_char)
		}
		advance()
	}

	// return a float if we have a decimal point else return a int
	if dot_cnt > 0 {
		return token{section.prev_section.section, section.prev_section.ln_count, TT_FLOAT, num_str}
	} else {
		return token{section.prev_section.section, section.prev_section.ln_count, TT_INT, num_str}
	}
}

// creates a token for the next character
func make_char() token {
	advance()
	char_str := ""

	if current_char != '\'' {
		char_str += string(current_char)
		advance()
		if current_char == '\'' {
			return token{section.prev_section.section, section.prev_section.ln_count, TT_CHAR, char_str}
		} else {
			NewError("CharStatementNotEnded", "You started a char in a single qoute but never closed it", currentPos(), true)
		}
	}
	return token{section.prev_section.section, section.prev_section.ln_count, TT_CHAR, char_str}
}

// lexes and returns a pointer to a token array
func Lex(text_ptr *string) *[]token {
	// sets the variable that holds the text to a byte array of the string that is given in the function
	text = []byte(*text_ptr)
	advance() // advance once to set the first character

	var tokens []token
	var sections []string

	for !is_eof {
		if in_mikas {
			tokens = append(tokens, make_str()) // in case we are in mikas we just make everything a big str and append it
			advance()
			in_mikas = false
		} else {
			// multiline comments, @..., comments and strings, chars, ints, floats, mikas
			switch string(current_char) {
			case " ":
				advance()
				break
			case "\t":
				advance()
				break
			case "\n":
				section.prev_section.ln_count += 1
				advance()
				break
			case "<":
				advance()
				if current_char == '=' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_LEQ, "<="})
					advance()
				} else {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_LTHEN, "<"})
				}
				break
			case ">":
				advance()
				if current_char == '=' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_GEQ, ">="})
					advance()
				} else {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_GTHEN, ">"})
				}
				break
			case "&":
				advance()
				if current_char == '&' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_AND, "&&"})
					advance()
				} else {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_KAND, "&"})
				}
				break
			case "|":
				advance()
				if current_char == '|' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_OR, "||"})
					advance()
				} else {
					NewError("IllegalTokenError", "A '|' was expected after a '|' but not found", currentPos(), true)
				}
				break
			case "!":
				advance()
				if current_char == '=' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_NEQ, "!="})
					advance()
				} else {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_NOT, "!"})
				}
				break
			case ".":
				tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_DOT, "."})
				advance()
				break
			case "+":
				advance()
				if current_char == '=' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_PLUSEQ, "+="})
					advance()
				} else {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_PLUS, "+"})
				}
				break
			case "-":
				advance()
				if current_char == '=' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_MINUSEQ, "-="})
					advance()
				} else if current_char == '>' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_ARROW, "->"})
					advance()
				} else {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_MINUS, "-"})
				}
				break
			case "*":
				advance()
				if current_char == '=' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_MULEQ, "*="})
					advance()
				} else {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_MUL, "*"})
				}
				break
			case "/":
				advance()
				if current_char == '=' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_DIVEQ, "/="})
					advance()
				} else if current_char == '/' {
					for current_char != '\n' && !is_eof {
						advance()
					}
				} else if current_char == '*' {
					advance()
					for !is_eof {
						if current_char == '*' {
							advance()
							if current_char == '/' {
								advance()
								break
							}
						}
						advance()
					}
					if is_eof {
						NewError("MultilineCommentNeverClosed", "A multiline comment was started but never closed", currentPos(), true)
					}
				} else {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_DIV, "/"})
				}
				break
			case "=":
				advance()
				if current_char == '=' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_EQ, "=="})
					advance()
				} else {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_ASSGN, "="})
				}
				break
			case "?":
				advance()
				if current_char == '=' {
					tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_REASSGN, "?="})
					advance()
				} else {
					NewError("IllegalTokenError", "A '=' was expected after a '?' but nor found", currentPos(), true)
				}
			case "(":
				tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_LPAREN, "("})
				advance()
				break
			case ")":
				tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_RPAREN, ")"})
				advance()
				break
			case "[":
				tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_LBRK, "["})
				advance()
				break
			case "]":
				tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_RBRK, "]"})
				advance()
				break
			case "{":
				if expect_mikas {
					in_mikas = true
					expect_mikas = false
				}
				tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_LCURL, "{"})
				advance()
				break
			case "}":
				tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_RCURL, "}"})
				advance()
				break
			case ",":
				tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_COMMA, ","})
				advance()
				break
			case ";":
				tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_SEMIC, ";"})
				advance()
				break
			case "%":
				tokens = append(tokens, token{section.prev_section.section, section.prev_section.ln_count, TT_PERCENT, "%"})
				advance()
				break
			case "\"":
				tokens = append(tokens, make_str())
				advance()
			case "'":
				tokens = append(tokens, make_char())
				advance()
			default:
				if strings.Contains(chars, string(current_char)) {
					tokens = append(tokens, make_id())
				} else if strings.Contains(numbers, string(current_char)) {
					tokens = append(tokens, make_number())
				} else {
					NewError("IllegalTokenError", "A token was not expected.", currentPos(), true)
				}
			}
		}
	}

	fmt.Println(tokens, sections)
	return &tokens
}
