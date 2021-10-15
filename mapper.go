package xorm

type IMapper interface {
	Obj2Table(string) string
	Table2Obj(string) string
}

type SnakeMapper struct{}

func (mapper SnakeMapper) Table2Obj(name string) string {
	return titleCasedName(name)
}

func (mapper SnakeMapper) Obj2Table(name string) string {
	return snakeCaseName(name)
}

// snakeCase to titleCase
// ab_cd_ef -> AbCdEf
// doesn't work for words start with _
func titleCasedName(name string) string {
	newStr := make([]rune, 0)
	upNextChar := true

	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			chr -= ('a' - 'A')
		case chr == '_':
			upNextChar = true
			continue
		}

		newStr = append(newStr, chr)
	}

	return string(newStr)
}

// titleCase to snakeCase
// AbCdEf -> ab_cd_ef
func snakeCaseName(name string) string {
	newstr := make([]rune, 0)
	firstTime := true

	for _, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if firstTime == true {
				firstTime = false
			} else {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}
