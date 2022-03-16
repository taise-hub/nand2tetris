package main

import (
	"assembler/code"
	"assembler/parser"
	"assembler/symboltable"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

const fileName = "../rect/rect.asm"
const compiledfileName = "../rect/_Rect.hack"

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
	createSymbolTable()

	for p.Advance() {
		switch p.CommandType() {
		case parser.A_COMMAND:
			symbol := convertSymbol(p.Symbol())
			code := "0" + fmt.Sprintf("%015s", strconv.FormatInt(int64(symbol), 2))
			fp.WriteString(code + "\n")
		case parser.C_COMMAND:
			comp, dest, jump := p.ParseC()
			if len(comp)+len(dest)+len(jump) == 0 {
				continue
			}
			code := "111" + code.Comp(comp) + code.Dest(dest) + code.Jump(jump)
			fp.WriteString(code + "\n")
		default:
			continue
		}
	}
}

func convertSymbol(lavel string) int {
	st := symboltable.GetSymbolTable()

	if isNum(lavel) {
		symbol, _ := strconv.Atoi(lavel)
		return symbol
	}
	return st.GetAddress(lavel)
}

func initSymbolTalbe() {
	symboltable := symboltable.GetSymbolTable()
	symboltable.AddEntry("SP", 0)
	symboltable.AddEntry("LCL", 1)
	symboltable.AddEntry("ARG", 2)
	symboltable.AddEntry("THIS", 3)
	symboltable.AddEntry("THAT", 4)
	symboltable.AddEntry("R0", 0)
	symboltable.AddEntry("R1", 1)
	symboltable.AddEntry("R2", 2)
	symboltable.AddEntry("R3", 3)
	symboltable.AddEntry("R4", 4)
	symboltable.AddEntry("R5", 5)
	symboltable.AddEntry("R6", 6)
	symboltable.AddEntry("R7", 7)
	symboltable.AddEntry("R8", 8)
	symboltable.AddEntry("R9", 9)
	symboltable.AddEntry("R10", 10)
	symboltable.AddEntry("R11", 11)
	symboltable.AddEntry("R12", 12)
	symboltable.AddEntry("R13", 13)
	symboltable.AddEntry("R14", 14)
	symboltable.AddEntry("R15", 15)
	symboltable.AddEntry("SCREEN", 16384)
	symboltable.AddEntry("KBD", 24576)
}

func createSymbolTable() {
	initSymbolTalbe()
	p, err := parser.New(fileName)
	if err != nil {
		log.Fatal(err)
	}
	symboltable := symboltable.GetSymbolTable()
	for p.Advance() {
		switch p.CommandType() {
		case parser.A_COMMAND, parser.C_COMMAND:
			symboltable.Increment()
		case parser.L_COMMAND:
			println(symboltable.GetCount())
			symboltable.AddEntry(p.Symbol(), symboltable.GetCount())
		default:
			continue
		}
	}

	symboltable.InitCount()
	p, err = parser.New(fileName)
	if err != nil {
		log.Fatal(err)
	}
	for p.Advance() {
		switch p.CommandType() {
		case parser.A_COMMAND:
			//ラベルかつsymboltableに登録されていない場合は登録する。
			if !isNum(p.Symbol()) && !symboltable.Contains(p.Symbol()) {
				symboltable.AddEntry(p.Symbol(), symboltable.GetCount()+0x10)
			}
		default:
			continue
		}
	}
}

func isNum(str string) bool {
	for _, c := range str {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
