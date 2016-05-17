package main

import (
	"bytes"
)

type LineArray struct {
	Lines [][]byte
}

func NewLineArray(text []byte) *LineArray {
	la := new(LineArray)
	split := bytes.Split(text, []byte("\n"))
	la.Lines = make([][]byte, len(split))
	for i := range split {
		la.Lines[i] = make([]byte, len(split[i]))
		copy(la.Lines[i], split[i])
	}

	return la
}

func (la *LineArray) String() string {
	return string(bytes.Join(la.Lines, []byte("\n")))
}

func (la *LineArray) NewlineBelow(y int) {
	la.Lines = append(la.Lines, []byte(" "))
	copy(la.Lines[y+2:], la.Lines[y+1:])
	la.Lines[y+1] = []byte("")
}

func (la *LineArray) Insert(pos Loc, value []byte) {
	x, y := pos.x, pos.y
	for i := 0; i < len(value); i++ {
		if value[i] == '\n' {
			la.Split(Loc{x, y})
			x = 0
			y++
			continue
		}
		la.insertByte(Loc{x, y}, value[i])
		x++
	}
}

func (la *LineArray) insertByte(pos Loc, value byte) {
	la.Lines[pos.y] = append(la.Lines[pos.y], 0)
	copy(la.Lines[pos.y][pos.x+1:], la.Lines[pos.y][pos.x:])
	la.Lines[pos.y][pos.x] = value
}

func (la *LineArray) JoinLines(a, b int) {
	la.Insert(Loc{len(la.Lines[a]), a}, la.Lines[b])
	la.DeleteLine(b)
}

func (la *LineArray) Split(pos Loc) {
	la.NewlineBelow(pos.y)
	la.Insert(Loc{0, pos.y + 1}, la.Lines[pos.y][pos.x:])
	la.DeleteToEnd(pos)
}

func (la *LineArray) Delete(start, end Loc) string {
	sub := la.Substr(start, end)
	if start.y == end.y {
		la.Lines[start.y] = append(la.Lines[start.y][:start.x], la.Lines[start.y][end.x:]...)
	} else {
		for i := start.y + 1; i <= end.y-1; i++ {
			la.DeleteLine(i)
		}
		la.DeleteToEnd(start)
		la.DeleteFromStart(Loc{end.x - 1, start.y + 1})
		la.JoinLines(start.y, start.y+1)
	}
	return sub
}

func (la *LineArray) DeleteToEnd(pos Loc) {
	la.Lines[pos.y] = la.Lines[pos.y][:pos.x]
}

func (la *LineArray) DeleteFromStart(pos Loc) {
	la.Lines[pos.y] = la.Lines[pos.y][pos.x+1:]
}

func (la *LineArray) DeleteLine(y int) {
	la.Lines = la.Lines[:y+copy(la.Lines[y:], la.Lines[y+1:])]
}

func (la *LineArray) DeleteByte(pos Loc) {
	la.Lines[pos.y] = la.Lines[pos.y][:pos.x+copy(la.Lines[pos.y][pos.x:], la.Lines[pos.y][pos.x+1:])]
}

func (la *LineArray) Substr(start, end Loc) string {
	if start.y == end.y {
		return string(la.Lines[start.y][start.x:end.x])
	}
	var str string
	str += string(la.Lines[start.y][start.x:]) + "\n"
	for i := start.y + 1; i <= end.y-1; i++ {
		str += string(la.Lines[i]) + "\n"
	}
	str += string(la.Lines[end.y][:end.x])
	return str
}
