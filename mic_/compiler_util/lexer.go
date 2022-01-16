package compiler_util

import (
	"fmt"
	"strings"
)

// all legal characters for func_names and stuff (numbers also work but that's done later on)
var chars string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
var numbers string = "0123456789"

// All the Tokens that exist
const (
	TT_BITAND     string = "BITAND"   // b&
	TT_BITOR      string = "BITOR"    // b|
	TT_BITXOR     string = "BITXOR"   // ^
	TT_BITNOT     string = "BITNOT"   // b!
	TT_LTHEN      string = "LTHEN"    // <  1
	TT_SHIFTL     string = "SHIFTL"   // << 2
	TT_LEQ        string = "LEQ"      // <=	3
	TT_GTHEN      string = "GTHEN"    // >	4
	TT_SHIFTR     string = "SHIFTR"   // >> 5
	TT_GEQ        string = "GEQ"      // >=	4
	TT_KAND       string = "KAND"     // &	5
	TT_AND        string = "AND"      // &&	6
	TT_OR         string = "OR"       // ||	7
	TT_NOT        string = "NOT"      // !	8
	TT_NEQ        string = "NEQ"      // !=	9
	TT_DOT        string = "DOT"      // .	10
	TT_PLUS       string = "PLS"      // +	11
	TT_PLUSEQ     string = "PLSEQ"    // +=	12
	TT_MINUS      string = "MIN"      // -	13
	TT_ARROW      string = "ARROW"    // ->	14
	TT_MINUSEQ    string = "MINEQ"    // -=	15
	TT_MUL        string = "MUL"      // *	16
	TT_MULEQ      string = "MULEQ"    // *=	17
	TT_DIV        string = "DIV"      // /	18
	TT_DIVEQ      string = "DIVEQ"    // /=	19
	TT_ASSGN      string = "ASSGN"    // =	20
	TT_EQ         string = "EQ"       // ==	21
	TT_REASSGN    string = "REASSGN"  // ?=	22
	TT_RPAREN     string = "RPAREN"   // (	23
	TT_LPAREN     string = "LPAREN"   // )	24
	TT_LBRK       string = "LBRK"     // [	25
	TT_RBRK       string = "RBRK"     // ]	26
	TT_LCURL      string = "LCURL"    // {	27
	TT_RCURL      string = "RCURL"    // }	28
	TT_COMMA      string = "COMMA"    // ,	29
	TT_SEMIC      string = "SEMIC"    // ;	30
	TT_PERCENT    string = "PERCENT"  // %	31
	TT_INT        string = "INT"      // numbers	32
	TT_FLOAT      string = "FLT"      // numbers.numbers	33
	TT_ID         string = "ID"       // a string of chars for example for a var name	34
	TT_STRING     string = "STR"      // "string goes here"	35
	TT_CHAR       string = "CHAR"     // 'single_character'	36
	TT_DEBUG      string = "DEBUG"    // @debug	37
	TT_OVERRIDE   string = "OVERRIDE" // @overwrite	38
	TT_EOF        string = "EOF"      // end of file	39
	TT_NUM_OF_OPS int    = 39
)

// is used as a part of the position and only exists because you can't have a struct in itself
type primitive_position struct {
	section  string
	ln_count int
}

// is used to keep track of the current position in files and to find the previous section once one section is done
type position struct {
	prev_section    []primitive_position
	prev_sec_len    int
	current_section primitive_position
}

// Token is what is given to the parser for each Token
type Token struct {
	section  string
	ln_count int
	type_    string
	value    string
}

var fName string
var text []byte
var text_position int = -1
var current_char byte = ' '
var is_eof bool = false
var section position = position{[]primitive_position{primitive_position{section: fName, ln_count: 0}}, 0, primitive_position{ln_count: 0, section: ""}}
var expect_mikas bool = false
var in_mikas bool = false

func currentPos() string {
	return fmt.Sprintf("%s at ln %d", section.current_section.section, section.current_section.ln_count)
}

// Lexer Functions:
// s forward once
func l_advance() {
	text_position += 1
	if text_position < len(text) {
		current_char = text[text_position]
	} else {
		// assigns the current char to ' ' because that's something that never happens as we ignore white spaces.
		current_char = ' '
		is_eof = true
	}
}

// creates a Token with the next characters and returns it
func make_str() Token {
	var mk_str string

	if in_mikas {
		cnt := 0
		// loops as long as not end of file and aslong as the mikas hasn't ended
		for !is_eof && current_char != '}' {
			if current_char == '\n' && cnt > 0 { // checks for a newline and if it's more than one character appended bcs else the c command would be pretty trashy
				mk_str += ";"
				section.current_section.ln_count++
			} else if current_char == '\n' && cnt <= 0 {
				section.current_section.ln_count++
				// don't do anything for a newline at the start of the function except incrementing
			} else if current_char == '\t' {
				// don't do anything for tabs
			} else {
				mk_str += string(current_char) // append every other character
			}
			l_advance()
			cnt++
		}
	} else {
		l_advance()
		for !is_eof && current_char != '"' {
			mk_str += string(current_char)
			l_advance()
		}
		l_advance()
	}
	if len(mk_str) > 1 {
		mk_str = fmt.Sprintf("\"%s\"", mk_str)
	}
	return Token{section.current_section.section, section.current_section.ln_count, TT_STRING, mk_str}
}

// creates a Token for the following characters
func make_id() Token {
	var id_str string

	// loops over the next chars and as long as the char is in the alphabet and all numbers it appends
	for !is_eof && strings.Contains(chars+numbers, string(current_char)) {
		id_str += string(current_char)
		l_advance()
	}
	// sets expect mikas to true if the id is mikas so the compiler can process inline assembly correctly
	if id_str == "mikas" {
		expect_mikas = true
	}
	return Token{section.current_section.section, section.current_section.ln_count, TT_ID, id_str}
}

// creates a Token for either a int or float based on if it has a decimal point
func make_number() Token {
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
		l_advance()
	}

	// return a float if we have a decimal point else return a int
	if dot_cnt > 0 {
		return Token{section.current_section.section, section.current_section.ln_count, TT_FLOAT, num_str}
	} else {
		return Token{section.current_section.section, section.current_section.ln_count, TT_INT, num_str}
	}
}

// creates a Token for the next character
func make_char() Token {
	l_advance()
	char_str := "'"

	if current_char != '\'' {
		char_str += string(current_char) + "'"
		l_advance()
		if current_char == '\'' {
			return Token{section.current_section.section, section.current_section.ln_count, TT_CHAR, char_str}
		} else {
			NewError("CharStatementNotEnded", "You started a char in a single qoute but never closed it", currentPos(), true)
		}
	}
	return Token{section.current_section.section, section.current_section.ln_count, TT_CHAR, char_str}
}

// lexes and returns a pointer to a Token array
func Lex(text_ptr *string, f_name string) *[]Token {
	// sets the variable that holds the text to a byte array of the string that is given in the function
	text = []byte(*text_ptr)
	fName = f_name
	l_advance() //  once to set the first character

	var Tokens []Token
	var sections []string

	section.current_section.ln_count = -1
	for !is_eof {
		if in_mikas {
			Tokens = append(Tokens, make_str()) // in case we are in mikas we just make everything a big str and append it
			in_mikas = false
		} else {
			// multiline comments, @..., comments and strings, chars, ints, floats, mikas
			switch string(current_char) {
			case " ":
				l_advance()
				break
			case "\t":
				l_advance()
				break
			case "\n":
				section.current_section.ln_count++
				l_advance()
				break
			case "<":
				l_advance()
				if current_char == '=' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_LEQ, "<="})
					l_advance()
				} else if current_char == '<' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_SHIFTL, "<<"})
					l_advance()
				} else {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_LTHEN, "<"})
				}
				break
			case ">":
				l_advance()
				if current_char == '=' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_GEQ, ">="})
					l_advance()
				} else if current_char == '>' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_SHIFTR, ">>"})
					l_advance()
				} else {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_GTHEN, ">"})
				}
				break
			case "&":
				l_advance()
				if current_char == '&' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_AND, "&&"})
					l_advance()
				} else {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_KAND, "&"})
				}
				break
			case "|":
				l_advance()
				if current_char == '|' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_OR, "||"})
					l_advance()
				} else {
					NewError("IllegalTokenError", "A '|' was expected after a '|' but not found", currentPos(), true)
				}
				break
			case "!":
				l_advance()
				if current_char == '=' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_NEQ, "!="})
					l_advance()
				} else {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_NOT, "!"})
				}
				break
			case ".":
				Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_DOT, "."})
				l_advance()
				break
			case "+":
				l_advance()
				if current_char == '=' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_PLUSEQ, "+="})
					l_advance()
				} else {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_PLUS, "+"})
				}
				break
			case "-":
				l_advance()
				if current_char == '=' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_MINUSEQ, "-="})
					l_advance()
				} else if current_char == '>' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_ARROW, "->"})
					l_advance()
				} else {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_MINUS, "-"})
				}
				break
			case "*":
				l_advance()
				if current_char == '=' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_MULEQ, "*="})
					l_advance()
				} else {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_MUL, "*"})
				}
				break
			case "/":
				l_advance()
				if current_char == '=' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_DIVEQ, "/="})
					l_advance()
				} else if current_char == '/' {
					for current_char != '\n' && !is_eof {
						l_advance()
					}
					section.current_section.ln_count++
					l_advance()
				} else if current_char == '*' {
					l_advance()
					for !is_eof {
						if current_char == '*' {
							l_advance()
							if current_char == '/' {
								l_advance()
								break
							}
						} else if current_char == '\n' {
							section.current_section.ln_count++
						}
						l_advance()
					}
					if is_eof {
						NewError("MultilineCommentNeverClosed", "A multiline comment was started but never closed", currentPos(), true)
					}
				} else {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_DIV, "/"})
				}
				break
			case "=":
				l_advance()
				if current_char == '=' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_EQ, "=="})
					l_advance()
				} else {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_ASSGN, "="})
				}
				break
			case "?":
				l_advance()
				if current_char == '=' {
					Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_REASSGN, "?="})
					l_advance()
				} else {
					NewError("IllegalTokenError", "A '=' was expected after a '?' but nor found", currentPos(), true)
				}
			case "(":
				Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_LPAREN, "("})
				l_advance()
				break
			case ")":
				Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_RPAREN, ")"})
				l_advance()
				break
			case "[":
				Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_LBRK, "["})
				l_advance()
				break
			case "]":
				Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_RBRK, "]"})
				l_advance()
				break
			case "{":
				if expect_mikas {
					in_mikas = true
					expect_mikas = false
				}
				Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_LCURL, "{"})
				l_advance()
				break
			case "}":
				Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_RCURL, "}"})
				l_advance()
				break
			case ",":
				Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_COMMA, ","})
				l_advance()
				break
			case ";":
				Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_SEMIC, ";"})
				l_advance()
				break
			case "%":
				Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_PERCENT, "%"})
				l_advance()
				break
			case "\"":
				Tokens = append(Tokens, make_str())
			case "'":
				Tokens = append(Tokens, make_char())
				l_advance()
			case "@":
				l_advance()
				if strings.Contains(chars, string(current_char)) {
					var id_str string

					// make a id for what ever comes after @ so for example: @section -> section
					for !is_eof && strings.Contains(chars+numbers, string(current_char)) {
						id_str += string(current_char)
						l_advance()
					}
					switch id_str {
					case "section":
						if current_char == '(' {
							l_advance()
							if current_char == '"' {
								sec_name := make_str().value

								if current_char == ')' {
									section.prev_section = append(section.prev_section, section.current_section)
									section.prev_sec_len++
									section = position{prev_section: section.prev_section, prev_sec_len: section.prev_sec_len, current_section: primitive_position{sec_name, 1}}
									sections = append(sections, sec_name)
									l_advance()
									section.current_section.ln_count = 1
								} else {
									NewError("ClosingParentheseExpectedError", "A closing parenthese expected after '@section(\"name\"' but not found at ", currentPos(), true)
								}
							} else {
								NewError("StringExpectedError", "A srting (\"\") was expected after '@section(' but wasn't found At ", currentPos(), true)
							}
						} else {
							NewError("OpeningParentheseExpected", "A opening paranthese was expected but not found after '@section' at ", currentPos(), true)
						}
					case "secend":
						if section.prev_sec_len > 0 {
							section.current_section = section.prev_section[section.prev_sec_len-1]
							section.prev_sec_len--
						}
					case "debug":
						Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_DEBUG, ""})
						l_advance()
					case "override":
						Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_OVERRIDE, ""})
						l_advance()
					default:
						NewError("InvalidCompilerFlagError", "A invalid compiler flag was found at ", currentPos(), true)
					}
				}
			case "^":
				Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_BITXOR, "^"})
				l_advance()
			default:
				if strings.Contains(chars, string(current_char)) {
					if current_char == 'b' {
						l_advance()
						if current_char == '&' {
							Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_BITAND, "b&"})
							l_advance()
						} else if current_char == '|' {
							Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_BITOR, "b|"})
							l_advance()
						} else if current_char == '!' {
							Tokens = append(Tokens, Token{section.current_section.section, section.current_section.ln_count, TT_BITNOT, "b!"})
							l_advance()
						} else {
							current_char = text[text_position-1]
							text_position--
							Tokens = append(Tokens, make_id())
						}
					} else {
						Tokens = append(Tokens, make_id())
					}
				} else if strings.Contains(numbers, string(current_char)) {
					Tokens = append(Tokens, make_number())
				} else {
					NewError("IllegalTokenError", "A Token was not expected.", currentPos(), true)
				}
			}
		}
	}
	Tokens = append(Tokens, Token{"", 0, TT_EOF, ""})
	return &Tokens
}
