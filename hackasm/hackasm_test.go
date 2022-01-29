package hackasm_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nobishino/gocode/hackasm"
)

// Parser, Code, SymbolTable, Main
// Parser: each instruction into its underlying fields
// Code: each field into binary
// Staged development
// - handle code without symbols first.
//

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

func TestInstructionCode(t *testing.T) {
	testcases := []struct {
		name string
		inst hackasm.Instruction
		want string
	}{
		{
			name: "A-inst",
			inst: hackasm.Instruction{
				Kind:  "A",
				Value: 1,
			},
			want: "0000000000000001",
		},
		{
			name: "A-inst 2",
			inst: hackasm.Instruction{
				Kind:  "A",
				Value: 32767,
			},
			want: "0111111111111111",
		},
		{
			name: "D+A",
			inst: hackasm.Instruction{
				Kind: "C",
				Comp: "D+A",
			},
			want: "1110000010000000",
		},
		{
			name: "D-1",
			inst: hackasm.Instruction{
				Kind: "C",
				Comp: "D-1",
			},
			want: "1110001110000000",
		},
		{
			name: "A-1",
			inst: hackasm.Instruction{
				Kind: "C",
				Comp: "A-1",
			},
			want: "1110110010000000",
		},
		{
			name: "ADM=D&M",
			inst: hackasm.Instruction{
				Kind: "C",
				Dest: "ADM",
				Comp: "D&M",
			},
			want: "1111000000111000",
		},
		{
			name: "D=D&M",
			inst: hackasm.Instruction{
				Kind: "C",
				Dest: "D",
				Comp: "D&M",
			},
			want: "1111000000010000",
		},
		{
			name: "D;JGT",
			inst: hackasm.Instruction{
				Kind: "C",
				Comp: "D",
				Jump: "JGT",
			},
			want: "1110001100000001",
		},
		{
			name: "D;JMP",
			inst: hackasm.Instruction{
				Kind: "C",
				Comp: "D",
				Jump: "JMP",
			},
			want: "1110001100000111",
		},
		{
			name: "AMD=1;JLE",
			inst: hackasm.Instruction{
				Kind: "C",
				Comp: "1",
				Dest: "AMD",
				Jump: "JLE",
			},
			want: "1110111111111110",
		},
		{
			name: "AM=M-1",
			inst: hackasm.Instruction{
				Kind: "C",
				Comp: "M-1",
				Dest: "AM",
			},
			want: "1111110010101000",
		},
	}
	for _, tt := range testcases {
		if len(tt.want) != 16 {
			t.Fatalf("case %q is incorrect", tt.want)
		}
		t.Run(tt.name, func(t *testing.T) {
			got := tt.inst.Code()
			if len(got) != 16 {
				t.Errorf("want length of 16, but got %d", len(got))
			}
			if got != tt.want {
				t.Errorf("want %q, but got %q", tt.want, got)
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
		{
			name: "handle label symbols",
			src: `//
@R0
@R1
(START)
M=A
@START
;JMP
	`,
			want: []hackasm.Instruction{
				{Kind: "A", Value: 0},
				{Kind: "A", Value: 1},
				{Kind: "C", Comp: "M=A"},
				{Kind: "A", Value: 2},
				{Kind: "C", Jump: "JMP"},
			},
		},
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

func TestAssemble(t *testing.T) {
	testcase := []struct {
		name string
	}{
		{"PongL"},
	}
	for _, tt := range testcase {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			src := open(t, filepath.Join("testdata", tt.name+".asm"))
			dest := new(bytes.Buffer)

			if err := hackasm.Assemble(src, dest); err != nil {
				t.Fatal(err)
			}

			expect := open(t, filepath.Join("testdata", tt.name+".hack"))

			equal(t, expect, dest)
		})
	}
}

func open(t *testing.T, path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		f.Close()
	})
	return f
}

func equal(t *testing.T, want, got io.Reader) {
	wantB := new(strings.Builder)
	gotB := new(strings.Builder)
	if _, err := io.Copy(wantB, want); err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(gotB, got); err != nil {
		t.Fatal(err)
	}
	if wantB.String() != gotB.String() {
		t.Errorf("does not match. want:%s, got:%s", wantB.String(), gotB.String())
	}
}
