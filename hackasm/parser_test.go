package hackasm_test

import (
	"testing"

	"github.com/nobishino/gocode/hackasm"
)

func TestParser(t *testing.T) {
	testcases := []struct {
		name string
		in   string
		want hackasm.Instruction
		err  bool
	}{
		{
			name: "A-inst",
			in:   "@15",
			want: hackasm.Instruction{
				Kind:  "A",
				Value: 15,
			},
		},
		{
			name: "A-inst",
			in:   "@200",
			want: hackasm.Instruction{
				Kind:  "A",
				Value: 200,
			},
		},
		{
			name: "C-inst",
			in:   "D+A",
			want: hackasm.Instruction{
				Kind: "C",
				Comp: "D+A",
			},
		},
		{
			name: "C-inst with destination",
			in:   "AD=M+1",
			want: hackasm.Instruction{
				Kind: "C",
				Dest: "AD",
				Comp: "M+1",
			},
		},
		{
			name: "C-inst with jump",
			in:   "AD=M+1;JEQ",
			want: hackasm.Instruction{
				Kind: "C",
				Dest: "AD",
				Comp: "M+1",
				Jump: "JEQ",
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			p := hackasm.NewParser()
			got := p.Parse(tt.in)
			if got != tt.want {
				t.Errorf("want %v, but got %v", tt.want, got)
			}
		})
	}
}

func TestParseLines(t *testing.T) {
	testcases := []struct {
		name string
		src  string
		want []hackasm.Instruction
	}{
		{
			name: "2 lines",
			src: `
@2
D=A`,
			want: []hackasm.Instruction{
				{Kind: "A", Value: 2},
				{Kind: "C", Dest: "D", Comp: "A"},
			},
		},
		{
			name: "empty line & spaces",
			src:  "\n @2 \nD=A",
			want: []hackasm.Instruction{
				{Kind: "A", Value: 2},
				{Kind: "C", Dest: "D", Comp: "A"},
			},
		},
		{
			name: "handle inline comment",
			src:  "// comment\n @2 // comment\nD=A",
			want: []hackasm.Instruction{
				{Kind: "A", Value: 2},
				{Kind: "C", Dest: "D", Comp: "A"},
			},
		},
		{
			name: "handle virtual register",
			src: `//
@R0
@R1
@R15`,
			want: []hackasm.Instruction{
				{Kind: "A", Value: 0},
				{Kind: "A", Value: 1},
				{Kind: "A", Value: 15},
			},
		},
		{
			name: "handle variable sybol",
			src: `//
@R0
@A
@B`,
			want: []hackasm.Instruction{
				{Kind: "A", Value: 0},
				{Kind: "A", Value: 16},
				{Kind: "A", Value: 17},
			},
		},
		{
			name: "handle variable symbol which starts with R",
			src: `//
@RVariable
@RVariable2`,
			want: []hackasm.Instruction{
				{Kind: "A", Value: 16},
				{Kind: "A", Value: 17},
			},
		},
		{
			name: "handle tab(trimming)",
			src: `//
	@RVariable
@RVariable2
	`,
			want: []hackasm.Instruction{
				{Kind: "A", Value: 16},
				{Kind: "A", Value: 17},
			},
		},
		// 		{
		// 			name: "handle label symbols",
		// 			src: `//
		// @R0
		// @R1
		// (START)
		// M=A
		// @START
		// ;JMP
		// 	`,
		// 			want: []hackasm.Instruction{
		// 				{Kind: "A", Value: 0},
		// 				{Kind: "A", Value: 1},
		// 				{Kind: "C", Comp: "M=A"},
		// 				{Kind: "A", Value: 2},
		// 				{Kind: "C", Jump: "JMP"},
		// 			},
		// 		},
		{
			name: "handle defined symbol",
			src: `//
	@SP
@LCL
@ARG
@THIS
@THAT
@SCREEN
@KBD
	`,
			want: []hackasm.Instruction{
				{Kind: "A", Value: 0},
				{Kind: "A", Value: 1},
				{Kind: "A", Value: 2},
				{Kind: "A", Value: 3},
				{Kind: "A", Value: 4},
				{Kind: "A", Value: 16384},
				{Kind: "A", Value: 24576},
			},
		},
		{
			name: "handle Add.asm",
			src: `//
// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/06/add/Add.asm

// Computes R0 = 2 + 3  (R0 refers to RAM[0])

@2
D=A
@3
D=D+A
@0
M=D
	`,
			want: []hackasm.Instruction{
				{Kind: "A", Value: 2},
				{Kind: "C", Dest: "D", Comp: "A"},
				{Kind: "A", Value: 3},
				{Kind: "C", Dest: "D", Comp: "D+A"},
				{Kind: "A", Value: 0},
				{Kind: "C", Dest: "M", Comp: "D"},
			},
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			p := hackasm.NewParser()
			gotInstructions := p.ParseLines(tt.src)

			for i, want := range tt.want {
				if i >= len(gotInstructions) {
					t.Fatalf("length does not match. want: %d, got %d", len(tt.want), len(gotInstructions))
				}
				got := gotInstructions[i]
				if got != want {
					t.Errorf("want %+v, but got %+v at index %d", want, got, i)
				}
			}
			if len(gotInstructions) > len(tt.want) {
				t.Fatalf("length does not match. want: %d, got %d", len(tt.want), len(gotInstructions))
			}
		})
	}
}
