package xorm

import (
//"reflect"
//"strings"
)

// name translation between struct, fields names and table, column names
type IMapper interface {
	Obj2Table(string) string
	Table2Obj(string) string
}

// SameMapper implements IMapper and provides same name between struct and
// database table
type SameMapper struct {
}

func (m SameMapper) Obj2Table(o string) string {
	return o
}

func (m SameMapper) Table2Obj(t string) string {
	return t
}

// SnakeMapper implements IMapper and provides name transaltion between
// struct and database table
type SnakeMapper struct {
}

func snakeCasedName(name string) string {
	newstr := make([]rune, 0)
	for idx, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if idx > 0 {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

/*func pascal2Sql(s string) (d string) {
	d = ""
	lastIdx := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			if lastIdx < i {
				d += s[lastIdx+1 : i]
			}
			if i != 0 {
				d += "_"
			}
			d += string(s[i] + 32)
			lastIdx = i
		}
	}
	d += s[lastIdx+1:]
	return
}*/

func (mapper SnakeMapper) Obj2Table(name string) string {
	return snakeCasedName(name)
}

func titleCasedName(name string) string {
	newstr := make([]rune, 0)
	upNextChar := true

	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			if 'a' <= chr && chr <= 'z' {
				chr -= ('a' - 'A')
			}
		case chr == '_':
			upNextChar = true
			continue
		}

		newstr = append(newstr, chr)
	}

	return string(newstr)
}

func (mapper SnakeMapper) Table2Obj(name string) string {
	return titleCasedName(name)
}
