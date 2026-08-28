package parser

func splitLine(line string) []string {
	splitBySpace := true
	words := make([]string, 0)
	currWord := ""
	for _, letter := range line {
		if letter == '"' {
			splitBySpace = !splitBySpace
		}
		if letter == ' ' && splitBySpace {
			words = append(words, currWord)
			currWord = ""
			continue
		}
		currWord += string(letter)
	}
	if currWord != "" {
		words = append(words, currWord)
	}
	return words
}
