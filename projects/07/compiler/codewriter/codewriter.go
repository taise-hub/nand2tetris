package codewriter

import (
	"fmt"
	"os"
)

type CodeWriter struct {
	labelCount int
	file       *os.File
}

func New(filename string) (*CodeWriter, error) {
	fp, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &CodeWriter{labelCount: 0, file: fp}, nil
}

// 与えられた算術コマンドをアセンブリコードに変換し, それを書き込む
func (cw *CodeWriter) WriteArithmetic(cmd string) {
	switch cmd {
	case "add", "sub", "and", "or":
		cw.writeBinaryArithmetic(cmd)
	case "neg", "not":
		cw.writeUnaryArithmetic(cmd)
	default:
		panic("invalid command found at WriteArithmetic()")
	}
}

// 1変数算術演算
func (cw *CodeWriter) writeUnaryArithmetic(cmd string) {
	var asm string
	switch cmd {
	case "neg":
		asm = "D=-M"
	case "not":
		asm = "D=!M"
	}

	cw.writePop()
	cw.file.WriteString(asm + "\n")
}

// 2変数算術演算
func (cw *CodeWriter) writeBinaryArithmetic(cmd string) {
	var asm string
	switch cmd {
	case "add":
		asm = "D=D+M"
	case "sub":
		asm = "D=D-M"
	case "and":
		asm = "D=D&M"
	case "or":
		asm = "D=D|M"
	}

	cw.writePop()
	cw.file.WriteString("D=M\n")
	cw.writePop()
	cw.file.WriteString(asm + "\n")
}

// 論理演算
func (cw *CodeWriter) WriteLogic(cmd string) {
	var asm string
	switch cmd {
	case "eq":
		asm = "JEQ"
	case "gt":
		asm = "JGT"
	case "lt":
		asm = "JLT"
	}

	cw.writePop()
	cw.file.WriteString("D=M\n")
	cw.writePop()
	cw.file.WriteString("D=D-M\n")
	cw.file.WriteString(fmt.Sprintf("@TRUE_%d\n", cw.labelCount))
	cw.file.WriteString(fmt.Sprintf("D;%s\n", asm))
	cw.file.WriteString(fmt.Sprintf("@NEXT_%d", cw.labelCount))
	cw.file.WriteString("D=0\nJMP;0\n")
	cw.file.WriteString(fmt.Sprintf("(TRUE_%d)\nD=-1", cw.labelCount))
	cw.file.WriteString(fmt.Sprintf("(NEXT_%d)\n", cw.labelCount))

	cw.labelCount++
}

// C_PUSHまたはC_POPコマンドをアセンブリコードに変換し、それを書き込む
func (cw *CodeWriter) WritePushPop(command string, segment string, index int) {
	panic("implement me")
}

// popした値はMに格納される
func (cw *CodeWriter) writePop() {
	cw.file.WriteString("@SP\nM=M-1\n")
	// cw.file.WriteString("@SP\nM=M-1\nA=M\n")
}

// 出力ファイルを閉じる
func (cw *CodeWriter) Close() {
	cw.file.Close()
}

// func (cw *CodeWriter) newLabel() string {
// 	label := fmt.Sprintf("LABEL_%d", cw.labelCount)
// 	cw.labelCount++
// 	return label
// }
