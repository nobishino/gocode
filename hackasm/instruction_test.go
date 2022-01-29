package hackasm_test

import (
	"testing"

	"github.com/nobishino/gocode/hackasm"
)

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
