package lexer

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

// returns true if char byte is a letter, avoid unicode
func isLetter(char byte) bool {
	// allow _, and ? in identifiers
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_' || char == '?'
}
