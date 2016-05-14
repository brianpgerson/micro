package main

// Loc stores a location
type Loc struct {
	x, y int
}

// LessThan returns true if b is smaller
func (a Loc) LessThan(b Loc) bool {
	if a.y < b.y {
		return true
	}
	if a.y == b.y && a.x < b.x {
		return true
	}
	return false
}

// GreaterThan returns true if b is bigger
func (a Loc) GreaterThan(b Loc) bool {
	if a.y > b.y {
		return true
	}
	if a.y == b.y && a.x > b.x {
		return true
	}
	return false
}
