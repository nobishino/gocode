package hackasm

func Assemble(src string) []string {
	var result []string

	p := NewParser()
	var instructions = p.ParseLines(src)
	for _, instruction := range instructions {
		result = append(result, instruction.Code())
	}
	return result
}
