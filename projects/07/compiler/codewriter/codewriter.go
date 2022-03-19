package codewriter

import (
	"fmt"
	"log"
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
		log.Fatal("invalid command found at WriteArithmetic()")
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
	cw.file.WriteString(fmt.Sprintf("@TRUE_%d\nD;%s\n@NEXT_%d\nD=0\nJMP;0\n(TRUE_%d)\nD=-1\n(NEXT_%d)\n",
		cw.labelCount, asm, cw.labelCount, cw.labelCount, cw.labelCount))
	cw.labelCount++
}

func (cw *CodeWriter) WritePushPop(cmd string, segment string, index int) {
	switch cmd {
	case "push":
		cw.WritePush(segment, index)
	case "pop":
		cw.WritePopToMemory(segment, index)
	}
}

// VMのpush命令をHackアセンブラに変換する
func (cw *CodeWriter) WritePush(segment string, index int) {
	switch segment {
	case "constant":
		cw.file.WriteString(fmt.Sprintf("@%d\nD=A\n", index))
	case "local":
		cw.file.WriteString(fmt.Sprintf("@LCL\nD=M\n@%d\nA=D+A\nD=M\n", index))
	case "argument":
		cw.file.WriteString(fmt.Sprintf("@ARG\nD=M\n@%d\nA=D+A\nD=M\n", index))
	case "this":
		cw.file.WriteString(fmt.Sprintf("@THIS\nD=M\n@%d\nA=D+A\nD=M\n", index))
	case "that":
		cw.file.WriteString(fmt.Sprintf("@THAT\nD=M\n@%d\nA=D+A\nD=M\n", index))
	case "pointer":
		cw.file.WriteString(fmt.Sprintf("@R3\nD=A\n@%d\nD=D+A\n", index))
	case "temp":
		cw.file.WriteString(fmt.Sprintf("@R5\nD=A\n@%d\nD=D+A\n", index))
	case "static":
		cw.file.WriteString(fmt.Sprintf("@%s.%d\nD=M\n", cw.file.Name(), index))
	}

	cw.file.WriteString("@SP\nA=M\nM=D\n@SP\nM=M+1")
}

// VMのpop命令をHackアセンブラに変換する
func (cw *CodeWriter) WritePopToMemory(segment string, index int) {
	switch segment {
	case "local":
		cw.file.WriteString(fmt.Sprintf("@LCL\nD=M\n@%d\nD=D+A\n", index))
	case "argument":
		cw.file.WriteString(fmt.Sprintf("@ARG\nD=M\n@%d\nD=D+A\n", index))
	case "this":
		cw.file.WriteString(fmt.Sprintf("@THIS\nD=M\n@%d\nD=D+A\n", index))
	case "that":
		cw.file.WriteString(fmt.Sprintf("@THAT\nD=M\n@%d\nD=D+A\n", index))
	case "pointer":
		cw.file.WriteString(fmt.Sprintf("@R3\nD=A\n@%d\nD=D+A\n", index))
	case "temp":
		cw.file.WriteString(fmt.Sprintf("@R5\nD=A\n@%d\nD=D+A\n", index))
	case "static":
		cw.file.WriteString(fmt.Sprintf("@%s.%d\nD=A\n", cw.file.Name(), index))
	}

	cw.file.WriteString("@R13\nM=D\n")
	cw.writePop()
	cw.file.WriteString("D=M\n@R13\nA=M\nM=D\n")
}

func (cw *CodeWriter) writePop() {
	cw.file.WriteString("@SP\nAM=M-1\n")
}

// 出力用ファイルを閉じる
func (cw *CodeWriter) Close() {
	cw.file.Close()
}
