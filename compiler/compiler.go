package compiler

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

type Compiler struct {
	ImportedLines []string
	NasmInstance  Nasm
	srcLines      []string
}

var BitProcessIs64 bool = true

func NewCompiler(entry string) *Compiler {
	return &Compiler{
		ImportedLines: make([]string, 0),
		NasmInstance:  *NewNASM(entry),
		srcLines:      make([]string, 0),
	}
}

func CompilerException(errMsg string) {
	fmt.Printf("Compiler Exception: %s", errMsg)
	os.Exit(1)
}

func (c *Compiler) HandleLines(lines []string) {
	c.srcLines = lines
}

const (
	IMPORT_SCOPE = 1
)

var symbols = [...]rune{'(', ')', '.', ',', ';', '=', '@'}

func isInSymbol(r rune) bool {
	for _, e := range symbols {
		if e == r {
			return true
		}
	}
	return false
}

func isSymbol(s string) bool {
	if len(s) != 1 {
		return false
	}
	return isInSymbol(rune(s[0]))
}
func isString(s string) bool {
	if len(s) < 2 {
		return false
	}
	return s[0] == '"' && s[len(s)-2] == '"'
}
func isIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	if unicode.IsDigit(rune(s[0])) {
		return false
	}
	for _, c := range s {
		if !(unicode.IsDigit(c) || unicode.IsLetter(c)) && c != '_' {
			return false
		}
	}
	return true
}

func tokenize(t string) (tokens []string) {
	built := ""
	insideDoubleQuote := false
	insideSingleQuote := false
	slashPrefix := false
	for _, c := range t {
		if c == '\'' && !(slashPrefix || insideDoubleQuote) {
			insideSingleQuote = !insideSingleQuote
		} else if c == '"' && !(slashPrefix || insideSingleQuote) {
			insideDoubleQuote = !insideDoubleQuote
		} else if c == '\\' {
			if slashPrefix {
				slashPrefix = false
			} else {
				slashPrefix = true
			}
		}

		if (c == ' ') && !(insideSingleQuote || insideDoubleQuote) {
			if built != "" {
				built = strings.Trim(built, " ")
				tokens = append(tokens, built)
				built = ""
			}
		} else if (isInSymbol(c)) && !(insideSingleQuote || insideDoubleQuote) {
			if built != "" {
				built = strings.Trim(built, " ")
				tokens = append(tokens, built)
				built = ""
			}
			tokens = append(tokens, string(c))
		} else {
			built += string(c)
		}
	}

	if built != "" {
		built = strings.Trim(built, " ")
		tokens = append(tokens, built)
		built = ""
	}

	commentIndex := -1
	{
		for i, t := range tokens {
			if t == ";" {
				commentIndex = i
				break
			}
		}
	}
	if commentIndex != -1 {
		tokens = tokens[:commentIndex]
	}

	return tokens
}

const (
	IMPORT    = "import"
	CALL_ADDR = "call"
	PUSH      = "push"
	POP       = "pop"

	PTR_T         = "ptr"
	BYTE_T        = "byte"
	WORD_T        = "word"
	DOUBLE_WORD_T = "dword"
	QUAD_WORD_T   = "qword"
	TEN_BYTES_T   = "tbyte"
	CONST_T       = "const"
)

func (c *Compiler) Compile() (reqPackages []string) {
	reqPackages = make([]string, 0)

	for i, l := range c.srcLines {
		lineTokens := tokenize(l)
		if len(lineTokens) == 0 {
			continue
		}
		switch lineTokens[0] {
		case IMPORT:
			if len(lineTokens) != 2 {
				CompilerException(fmt.Sprintf("LINE %d - Expected import [PKG]", i+1))
			}
			if !isString(lineTokens[1]) {
				CompilerException(fmt.Sprintf("LINE %d - Expected string as the package name", i+1))
			}
			reqPackages = append(reqPackages, lineTokens[1][1:len(lineTokens[1])-2])
		case CALL_ADDR:
			if len(lineTokens) < 2 {
				CompilerException(fmt.Sprintf("LINE %d - Expected call [ADDR]", i+1))
			}
			if lineTokens[1] == "@" {
				if len(lineTokens) != 2 {
					CompilerException(fmt.Sprintf("LINE %d - Cannot use macro accumilator with expr", i+1))
				}
				c.NasmInstance.TextSection.CallAddr(macroTotalExpr + macroTotalExprEnd)
			} else {
				c.NasmInstance.TextSection.CallAddr(strings.Join(lineTokens[1:], " "))
			}
		case PUSH:
			if len(lineTokens) < 2 {
				CompilerException(fmt.Sprintf("LINE %d - Expected push [VALUE]", i+1))
			}
			if lineTokens[1] == "@" {
				if len(lineTokens) != 2 {
					CompilerException(fmt.Sprintf("LINE %d - Cannot use macro accumilator with expr", i+1))
				}
				c.NasmInstance.TextSection.PushData(macroTotalExpr + macroTotalExprEnd)
			} else {
				c.NasmInstance.TextSection.PushData(strings.Join(lineTokens[1:], " "))
			}
		case POP:
			if len(lineTokens) < 2 {
				CompilerException(fmt.Sprintf("LINE %d - Expected pop [DST]", i+1))
			}
			c.NasmInstance.TextSection.PopInto(strings.Join(lineTokens[1:], " "))

		case "@":
			if len(lineTokens) < 2 {
				CompilerException(fmt.Sprintf("LINE %d - Expected @[MACRO NAME]", i+1))
			}
			if !isMacroName(lineTokens[1]) {
				CompilerException(fmt.Sprintf("LINE %d - Invalid macro name", i+1))
			}
			if macroNeedsNoArgs(lineTokens[1]) && len(lineTokens[1]) != 2 {
				CompilerException(fmt.Sprintf("LINE %d - This macro requires no arguments", i+1))
			}
			handleMacro(i+1, lineTokens[1], lineTokens[2:]...)

		case PTR_T:
			fallthrough
		case CONST_T:
			fallthrough
		case BYTE_T:
			fallthrough
		case WORD_T:
			fallthrough
		case DOUBLE_WORD_T:
			fallthrough
		case QUAD_WORD_T:
			fallthrough
		case TEN_BYTES_T:
			if len(lineTokens) < 3 {
				CompilerException(fmt.Sprintf("LINE %d - Expected [TYPE] [ALIAS] ...", i+1))
			}
			if !isIdentifier(lineTokens[1]) {
				CompilerException(fmt.Sprintf("LINE %d - Expected [TYPE] [ALIAS(IDENTIFIER)] ...", i+1))
			}
			if lineTokens[2] != "=" {
				CompilerException(fmt.Sprintf("LINE %d - Expected [TYPE] [ALIAS] = ...", i+1))
			}
			{
				alias := lineTokens[1]
				expr := strings.Join(lineTokens[3:], " ")
				switch lineTokens[0] {
				case PTR_T:
					if BitProcessIs64 {
						c.NasmInstance.DataSection.DefineQuadWord(alias, expr)
					} else {
						c.NasmInstance.DataSection.DefineDoubleWord(alias, expr)
					}
				case CONST_T:
					c.NasmInstance.DataSection.DefineConstant(alias, expr)
				case BYTE_T:
					c.NasmInstance.DataSection.DefineByte(alias, expr)
				case WORD_T:
					c.NasmInstance.DataSection.DefineWord(alias, expr)
				case DOUBLE_WORD_T:
					c.NasmInstance.DataSection.DefineDoubleWord(alias, expr)
				case QUAD_WORD_T:
					c.NasmInstance.DataSection.DefineQuadWord(alias, expr)
				case TEN_BYTES_T:
					c.NasmInstance.DataSection.DefineTenBytes(alias, expr)
				}
			}

		}
	}

	return reqPackages
}
