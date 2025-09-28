package symspell

type candidate struct {
	text  string
	runes []rune
}

func (c *candidate) Len() int {
	return len(c.runes)
}
