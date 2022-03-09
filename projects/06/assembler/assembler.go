package main

import (
	"assembler/code"
	"assembler/parser"
	"fmt"
	"log"
	"os"
	"strconv"
)

const fileName = "../pong/PongL.asm"
const compiledfileName = "../pong/_PongL.hack"

func main() {
	fp, err := os.Create(compiledfileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	p, err := parser.New(fileName)
	if err != nil {
		log.Fatal(err)
	}

	for p.Advance() {
		switch p.CommandType() {
		case parser.A_COMMAND, parser.L_COMMAND:
			symbol, err := strconv.Atoi(p.Symbol())
			if err != nil {
				log.Fatal(err)
			}
			code := "0" + fmt.Sprintf("%015s", strconv.FormatInt(int64(symbol), 2))
			fp.WriteString(code + "\n")
		default:
			comp, dest, jump := p.ParseC()
			code := "111" + code.Comp(comp) + code.Dest(dest) + code.Jump(jump)
			fp.WriteString(code + "\n")
		}
	}
}
