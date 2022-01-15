package hackasm_test

import (
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
			got := hackasm.Parse(tt.in)
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
