package lexem

import "strings"

type LexemType int

const (
	TypeL LexemType = iota
	ServiceL
	OpenParenL
	CloseParenL
	CommaL
	ReturnArrowL
	OpenBraceL
	CloseBraceL
	MethodL
	IdentefierL
	PackageL
	ObjectL
)

var typeIdentifiers = []string{"byte", "uint32", "uint64", "uint8", "string", "bool", "object"}
var pairTypes = map[string]LexemType{"service": ServiceL, "(": OpenParenL, ")": CloseParenL, ",": CommaL, "->": ReturnArrowL, "{": OpenBraceL, "}": CloseBraceL, "Handler": MethodL, "Stream": MethodL, "package": PackageL, "object": ObjectL}

type Lexem struct {
	Type  LexemType
	Value string
}

func LexemParse(input string) []*Lexem {
	words := Split(input, pairTypes)
	lexems := make([]*Lexem, 0, len(words))

	for _, word := range words {
		var lexemType LexemType
		lexemType, ok := pairTypes[word]
		if ok {
			lexems = append(lexems, &Lexem{
				Type:  lexemType,
				Value: word,
			})
			continue
		}

		var found bool
		for _, typeIdentefier := range typeIdentifiers {
			if strings.Contains(word, typeIdentefier) {
				lexemType = TypeL
				found = true
				break
			}
		}
		if found {
			lexems = append(lexems, &Lexem{
				Type:  lexemType,
				Value: word,
			})
			continue
		}

		lexems = append(lexems, &Lexem{
			Type:  IdentefierL,
			Value: word,
		})

	}

	return lexems
}

func Split(s string, delims map[string]LexemType) []string {
	var result []string
	lastIndex := -1
	runes := []rune(s)
	var i int

	for i < len(runes) {
		if s[i] == ' ' || s[i] == '\n' || s[i] == '\t' {
			if lastIndex+1 < i-1 {
				result = append(result, s[lastIndex+1:i])
			}
			lastIndex = i
			i++
			continue
		}
		if _, ok := delims[string(s[i])]; ok {
			if lastIndex+1 <= i {
				if lastIndex+1 != i {
					result = append(result, s[lastIndex+1:i])
				}
				result = append(result, string(s[i]))
			}
			lastIndex = i
			i++
			continue
		}

		i++
	}

	return result
}
