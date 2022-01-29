package hackasm_test

import (
	"testing"

	"github.com/nobishino/gocode/hackasm"
)

func TestInstructionCode(t *testing.T) {
	A := func(value uint16) hackasm.Instruction {
		return hackasm.Instruction{
			Kind:  "A",
			Value: value,
		}
	}
	C := func(dest, comp, jump string) hackasm.Instruction {
		return hackasm.Instruction{
			Kind: "C",
			Dest: dest,
			Comp: comp,
			Jump: jump,
		}
	}
	testcases := []struct {
		inst hackasm.Instruction
		want string
	}{
		{A(1), "0000000000000001"},
		{A(32767), "0111111111111111"},
		{C("", "D+A", ""), "1110000010000000"},
		{C("", "D-1", ""), "1110001110000000"},
		{C("", "A-1", ""), "1110110010000000"},
		{C("ADM", "D&M", ""), "1111000000111000"},
		{C("D", "D&M", ""), "1111000000010000"},
		{C("", "D", "JGT"), "1110001100000001"},
		{C("", "D", "JMP"), "1110001100000111"},
		{C("AMD", "1", "JLE"), "1110111111111110"},
		{C("AM", "M-1", ""), "1111110010101000"},
	}
	for _, tt := range testcases {
		if len(tt.want) != 16 {
			t.Fatalf("case %q is incorrect", tt.want)
		}
		t.Run(tt.inst.String(), func(t *testing.T) {
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
