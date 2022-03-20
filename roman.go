package printtree

var romanMap = []struct {
	value int
	key   string
}{
	{value: 1000, key: "m"},
	{value: 900, key: "cm"},
	{value: 500, key: "d"},
	{value: 400, key: "cd"},
	{value: 100, key: "c"},
	{value: 90, key: "xc"},
	{value: 50, key: "l"},
	{value: 40, key: "xl"},
	{value: 10, key: "x"},
	{value: 9, key: "ix"},
	{value: 5, key: "v"},
	{value: 4, key: "iv"},
	{value: 1, key: "i"},
}

func RomanNumerals(n int) string {
	s := ""
	for n > 0 {
		for _, rn := range romanMap {
			if n >= rn.value {
				s += rn.key
				n -= rn.value
				break
			}
		}
	}
	return s
}
