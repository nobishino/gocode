package hackasm

func Assemble(src string) []string {
	var result []string
	var instructions = ParseLines(src)
	for _, instruction := range instructions {
		result = append(result, instruction.Code())
	}
	return result
}
