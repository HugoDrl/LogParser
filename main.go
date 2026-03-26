package main

import "git.hugoderlyn.com/Hugo/goLogParser.git/parser"

func main() {
	content, _ := parser.ReadFile("./test.txt")
	parser.ParseLog(string(content))
}
