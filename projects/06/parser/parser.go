package parser

import (
	"bufio"
	"os"
	"strings"
)

type CommandType string

const (
	A_COMMAND CommandType = "A" // @Xxx
	C_COMMAND CommandType = "C" // dest=comp;jump
	L_COMMAND CommandType = "L" // (Xxx)
)

type Parser struct {
	file    *os.File
	scanner *Scanner
}

type Scanner struct {
	*bufio.Scanner
}

func (s *Scanner) Text() string {
	return strings.TrimSpace(s.Scanner.Text())
}

func New(fileName string) (*Parser, error) {
	fp, err := os.Open(fileName)
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
func (p *Parser) advance() bool {
	scan := p.scanner.Scan()
	if p.scanner.Text() == "" {
		scan = p.scanner.Scan()
	}
	return scan
}

// 現コマンドの種類を返す。
func (p *Parser) commandType() CommandType {
	switch {
	case strings.HasPrefix(p.scanner.Text(), "@"):
		return A_COMMAND
	case strings.HasPrefix(p.scanner.Text(), "(") && strings.HasPrefix(p.scanner.Text(), ")"):
		return L_COMMAND
	default:
		return C_COMMAND
	}
}

// 現コマンド@Xxx, (Xxx)のXxxの部分を返す。
// commandType()がA_COMMANDまたは, L_COMMANDの場合に呼ぶ。
func (p *Parser) symbol() string {
	command := p.scanner.Text()
	switch p.commandType() {
	case A_COMMAND:
		return command[1:]
	case L_COMMAND:
		return command[1 : len(command)-1]
	default:
		panic("command is not A_COMMAND OR L_COMMAND")
	}
}

// 現C命令のdestニーモニックを返す。
// commandType()がC_COMMANDの場合に呼ぶ。
func (p *Parser) dest() string {
	command := p.scanner.Text()
	return command[:strings.Index(command, "=")]
}

// 現C命令のcompニーモニックを返す。
// commandType()がC_COMMANDの場合に呼ぶ。
func (p *Parser) comp() string {
	command := p.scanner.Text()
	if strings.Contains(command, "=") {
		command = command[strings.Index(command, "="):]
	}
	if strings.Contains(command, ";") {
		command = command[:strings.Index(command, ";")]
	}
	return command
}

// 現C命令のjumpニーモニックを返す。
// commandType()がC_COMMANDの場合に呼ぶ。
func (p *Parser) jump() string {
	command := p.scanner.Text()
	return command[strings.Index(command, ";"):]
}
