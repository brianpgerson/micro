package main

// FromCharPos converts from a character position to an x, y position
func FromCharPos(loc int, buf *Buffer) Loc {
	charNum := 0
	x, y := 0, 0

	lineLen := Count(buf.Line(y)) + 1
	for charNum+lineLen <= loc {
		charNum += lineLen
		y++
		lineLen = Count(buf.Line(y)) + 1
	}
	x = loc - charNum

	return Loc{x, y}
}

// ToCharPos converts from an x, y position to a character position
func ToCharPos(start Loc, buf *Buffer) int {
	x, y := start.x, start.y
	loc := 0
	for i := 0; i < y; i++ {
		// + 1 for the newline
		loc += Count(buf.Line(i)) + 1
	}
	loc += x
	return loc
}

// Loc stores a location
type Loc struct {
	x, y int
}

// LessThan returns true if b is smaller
func (l Loc) LessThan(b Loc) bool {
	if l.y < b.y {
		return true
	}
	if l.y == b.y && l.x < b.x {
		return true
	}
	return false
}

// GreaterThan returns true if b is bigger
func (l Loc) GreaterThan(b Loc) bool {
	if l.y > b.y {
		return true
	}
	if l.y == b.y && l.x > b.x {
		return true
	}
	return false
}

func (l Loc) right(n int, buf *Buffer) Loc {
	if l == buf.End() {
		return l
	}
	var res Loc
	if l.x < len(buf.CurLine) {
		res = Loc{l.x + 1, l.y}
	} else {
		res = Loc{0, l.y + 1}
	}
	return res
}
func (l Loc) left(n int, buf *Buffer) Loc {
	if l == buf.Start() {
		return l
	}
	var res Loc
	if l.x > 0 {
		res = Loc{l.x - 1, l.y}
	} else {
		res = Loc{len(buf.Line(l.y - 1)), l.y - 1}
	}
	return res
}

func (l Loc) Move(n int, buf *Buffer) Loc {
	if n > 0 {
		return l.right(n, buf)
	}
	return l.left(Abs(n), buf)
}

// func (l Loc) DistanceTo(b Loc, buf *Buffer) int {
//
// }
