package fermat

func pow(x, y uint) uint {
	var result uint = 1
	for ; y > 0; x, y = x*x, y>>1 {
		if y%2 == 1 {
			result *= x
		}
	}
	return result
}
