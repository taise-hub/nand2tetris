// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

(LOOP)
    @KBD
    D=M
    @PRESSED
    D;JGT     // KBD > 0
    @flag
    M=0
    @DRAW_INIT
    0;JMP
    
(PRESSED)
    @flag
    M=-1      // -1=11111111なので横一列真っ黒

(DRAW_INIT)
    @i
    M=0

(DRAW)
    @SCREEN
    D=A
    @i
    D=D+M
    @position
    M=D

    @flag
    D=M
    @position
    A=M       //positionには, 画素に対応した番地が入っているので,Aにその値を入れてMでアクセスする
    M=D

    @i
    M=M+1      // i++
    D=M
    @8192      // 8192x4320
    D=D-A
    @LOOP      // if (i>=8192) goto LOOP
    D;JGE

    @DRAW
    0;JMP