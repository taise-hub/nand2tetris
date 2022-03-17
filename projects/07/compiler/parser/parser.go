package parser

import (
	"bufio"
	"os"
	"strings"
)

const (
	C_ARITHMETIC = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

type Parser struct {
	file    *os.File
	scanner *Scanner
}

type Scanner struct {
	*bufio.Scanner
}

func (s *Scanner) Text() string {
	line := strings.TrimSpace(s.Scanner.Text())
	if strings.Contains(line, "//") {
		line = line[:strings.Index(line, "//")]
	}
	return strings.TrimSpace(line)
}

func New(filename string) (*Parser, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := &Scanner{
		Scanner: bufio.NewScanner(fp),
	}
	return &Parser{file: fp, scanner: scanner}, nil
}

// 入力から次のコマンドを読み, それを現在のコマンドにする。
// 次のコマンドがない場合, falseを返す
func (p *Parser) Advance() bool {
	b := p.scanner.Scan()
	if b && p.scanner.Text() == "" {
		b = p.Advance()
	}
	return b
}

func (p *Parser) CommandType() int {
	return p.commandType()
}

// 現コマンドの種類を返す。
func (p *Parser) commandType() int {
	arg0 := strings.Split(p.scanner.Text(), " ")[0]
	switch arg0 {
	case "push":
		return C_PUSH
	case "pop":
		return C_POP
	case "label":
		return C_LABEL
	case "goto":
		return C_GOTO
	case "if-goto":
		return C_IF
	case "function":
		return C_FUNCTION
	case "return":
		return C_RETURN
	default:
		return C_ARITHMETIC
	}
}

// 現コマンドの最初の引数を返す。
// C_ARITHMETICの場合, コマンド自体を返す。
// C_RETURNの場合, 本ルーチンは呼ばない。
func (p *Parser) Arg1() string {
	return strings.Split(p.scanner.Text(), " ")[1]
}

//現コマンドの2番目の引数を返す。
// C_PUSH, C_POP, C_FUNCTION, C_CALLの場合のみ本ルーチンを呼ぶ。
func (p *Parser) Arg2() string {
	return strings.Split(p.scanner.Text(), " ")[2]
}
