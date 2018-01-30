package corelang

import (
	"errors"
	"strconv"
	"strings"
	"text/scanner"
)

type iToken interface{}
type tokenChar rune
type tokenComment string
type tokenFloat float64
type tokenIdent string
type tokenInt int64
type tokenUInt uint64
type tokenOpish string
type tokenSepish string
type tokenStr string

func lex(src string) (tokenStream []iToken, err error) { // accumulating in a return slice defeats the idea of scalable streaming (eg. parallel lex-parse pipeline), but really no matter in this toy
	var lexer scanner.Scanner
	lexer.Init(strings.NewReader(src)).Filename = src
	lexer.Mode = scanner.ScanChars | scanner.ScanComments | scanner.ScanFloats | scanner.ScanIdents | scanner.ScanInts | scanner.ScanRawStrings | scanner.ScanStrings
	lexer.Error = func(_ *scanner.Scanner, msg string) { err = errors.New(msg) }
	for tok := lexer.Scan(); (err == nil) && (tok != scanner.EOF); tok = lexer.Scan() {
		sym := lexer.TokenText()
		switch tok {
		case scanner.Char, scanner.RawString, scanner.String:
			println(sym)
		case scanner.Float:
			var f float64
			if f, err = strconv.ParseFloat(sym, 64); err == nil {
				tokenStream = append(tokenStream, tokenFloat(f))
			}
		case scanner.Int:
			var i int64
			if i, err = strconv.ParseInt(sym, 0, 64); err == nil {
				tokenStream = append(tokenStream, tokenInt(i))
			} else if numerr, _ := err.(*strconv.NumError); numerr != nil && numerr.Err == strconv.ErrRange {
				var u uint64
				if u, err = strconv.ParseUint(sym, 0, 64); err == nil {
					tokenStream = append(tokenStream, tokenUInt(u))
				}
			}
		default:
			switch sym {
			case ":", "+", "-", "*", "/", "!", "?", "^", "§", "°", "$", "%", "&", "|", "<", ">", "·", "×", "÷", "…", "±", "\\", "´", "~", "#", "@":
				tokenStream = append(tokenStream, tokenOpish(sym))
			case "(", ")", "{", "}", "[", "]", ",", ";":
				tokenStream = append(tokenStream, tokenSepish(sym))
			default:
				err = errors.New("unrecognized token: " + sym)
			}
		}
	}
	return
}
