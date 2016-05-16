package main

import (
	"strings"
)

// The Cursor struct stores the location of the cursor in the view
// The complicated part about the cursor is storing its location.
// The cursor must be displayed at an x, y location, but since the buffer
// uses a rope to store text, to insert text we must have an index. It
// is also simpler to use character indicies for other tasks such as
// selection.
type Cursor struct {
	Loc

	Buf *Buffer

	// Last cursor x position
	lastVisualX int

	// The current selection as a range of character numbers (inclusive)
	curSelection [2]Loc
	// The original selection as a range of character numbers
	// This is used for line and word selection where it is necessary
	// to know what the original selection was
	origSelection [2]Loc
}

// ResetSelection resets the user's selection
func (c *Cursor) ResetSelection() {
	c.curSelection[0] = c.Buf.Start()
	c.curSelection[1] = c.Buf.Start()
}

// HasSelection returns whether or not the user has selected anything
func (c *Cursor) HasSelection() bool {
	return c.curSelection[0] != c.curSelection[1]
}

// DeleteSelection deletes the currently selected text
func (c *Cursor) DeleteSelection() {
	if c.GetSelection() == "" {
		return
	}

	if c.curSelection[0].GreaterThan(c.curSelection[1]) {
		c.Buf.Remove(c.curSelection[1], c.curSelection[0])
		c.Loc = c.curSelection[1]
	} else {
		c.Buf.Remove(c.curSelection[0], c.curSelection[1])
		c.Loc = c.curSelection[0]
	}
}

// GetSelection returns the cursor's selection
func (c *Cursor) GetSelection() string {
	if c.curSelection[0].GreaterThan(c.curSelection[1]) {
		return c.Buf.Substr(c.curSelection[1], c.curSelection[0])
	}
	return c.Buf.Substr(c.curSelection[0], c.curSelection[1])
}

// SelectLine selects the current line
func (c *Cursor) SelectLine() {
	c.Start()
	c.curSelection[0] = c.Loc
	c.End()
	if c.Buf.NumLines-1 > c.y {
		c.curSelection[1] = c.Loc.Move(1, c.Buf)
	} else {
		c.curSelection[1] = c.Loc
	}

	c.origSelection = c.curSelection
}

// AddLineToSelection adds the current line to the selection
func (c *Cursor) AddLineToSelection() {
	loc := c.Loc

	if loc.LessThan(c.origSelection[0]) {
		c.Start()
		c.curSelection[0] = c.Loc
		c.curSelection[1] = c.origSelection[1]
	}
	if loc.GreaterThan(c.origSelection[1]) {
		c.End()
		c.curSelection[1] = c.Loc
		c.curSelection[0] = c.origSelection[0]
	}

	if loc.LessThan(c.origSelection[1]) && loc.GreaterThan(c.origSelection[0]) {
		c.curSelection = c.origSelection
	}
}

// SelectWord selects the word the cursor is currently on
func (c *Cursor) SelectWord() {
	// if len(c.Buf.CurLine) == 0 {
	// 	return
	// }
	//
	// if !IsWordChar(string(c.RuneUnder(c.x))) {
	// 	loc := c.Loc
	// 	c.curSelection[0] = loc
	// 	c.curSelection[1] = loc.Move(1, c.Buf)
	// 	c.origSelection = c.curSelection
	// 	return
	// }
	//
	// forward, backward := c.x, c.x
	//
	// for backward > 0 && IsWordChar(string(c.RuneUnder(backward-1))) {
	// 	backward--
	// }
	//
	// c.curSelection[0] = ToCharPos(Loc{backward, c.y}, c.Buf)
	// c.origSelection[0] = c.curSelection[0]
	//
	// for forward < Count(c.Buf.CurLine)-1 && IsWordChar(string(c.RuneUnder(forward+1))) {
	// 	forward++
	// }
	//
	// c.curSelection[1] = ToCharPos(Loc{forward + 1, c.y}, c.Buf)
	// c.origSelection[1] = c.curSelection[1]
	// c.Loc = c.curSelection[1]
}

// AddWordToSelection adds the word the cursor is currently on to the selection
func (c *Cursor) AddWordToSelection() {
	// loc := c.Loc
	//
	// if loc > c.origSelection[0] && loc < c.origSelection[1] {
	// 	c.curSelection = c.origSelection
	// 	return
	// }
	//
	// if loc < c.origSelection[0] {
	// 	backward := c.x
	//
	// 	for backward > 0 && IsWordChar(string(c.RuneUnder(backward-1))) {
	// 		backward--
	// 	}
	//
	// 	c.curSelection[0] = ToCharPos(backward, c.y, c.Buf)
	// 	c.curSelection[1] = c.origSelection[1]
	// }
	//
	// if loc > c.origSelection[1] {
	// 	forward := c.x
	//
	// 	for forward < Count(c.Buf.CurLine)-1 && IsWordChar(string(c.RuneUnder(forward+1))) {
	// 		forward++
	// 	}
	//
	// 	c.curSelection[1] = ToCharPos(forward, c.y, c.Buf) + 1
	// 	c.curSelection[0] = c.origSelection[0]
	// }
	//
	// c.SetLoc(c.curSelection[1])
}

// SelectTo selects from the current cursor location to the given location
func (c *Cursor) SelectTo(loc Loc) {
	if loc.GreaterThan(c.origSelection[0]) {
		c.curSelection[0] = c.origSelection[0]
		c.curSelection[1] = loc
	} else {
		c.curSelection[0] = loc
		c.curSelection[1] = c.origSelection[0]
	}
}

// WordRight moves the cursor one word to the right
func (c *Cursor) WordRight() {
	c.Right()
	for IsWhitespace(c.RuneUnder(c.x)) {
		if c.x == Count(c.Buf.CurLine) {
			return
		}
		c.Right()
	}
	for !IsWhitespace(c.RuneUnder(c.x)) {
		if c.x == Count(c.Buf.CurLine) {
			return
		}
		c.Right()
	}
}

// WordLeft moves the cursor one word to the left
func (c *Cursor) WordLeft() {
	c.Left()
	for IsWhitespace(c.RuneUnder(c.x)) {
		if c.x == 0 {
			return
		}
		c.Left()
	}
	for !IsWhitespace(c.RuneUnder(c.x)) {
		if c.x == 0 {
			return
		}
		c.Left()
	}
	c.Right()
}

// RuneUnder returns the rune under the given x position
func (c *Cursor) RuneUnder(x int) rune {
	line := []rune(c.Buf.CurLine)
	if len(line) == 0 {
		return '\n'
	}
	if x >= len(line) {
		return '\n'
	} else if x < 0 {
		x = 0
	}
	return line[x]
}

// Up moves the cursor up one line (if possible)
func (c *Cursor) Up() {
	if c.y > 0 {
		c.y--

		runes := []rune(c.Buf.CurLine)
		c.x = c.GetCharPosInLine(c.y, c.lastVisualX)
		if c.x > len(runes) {
			c.x = len(runes)
		}
	}
	c.UpdateCurLine()
}

// Down moves the cursor down one line (if possible)
func (c *Cursor) Down() {
	if c.y < c.Buf.NumLines-1 {
		c.y++

		runes := []rune(c.Buf.CurLine)
		c.x = c.GetCharPosInLine(c.y, c.lastVisualX)
		if c.x > len(runes) {
			c.x = len(runes)
		}
	}
	c.UpdateCurLine()
}

// Left moves the cursor left one cell (if possible) or to the last line if it is at the beginning
func (c *Cursor) Left() {
	c.Loc = c.Loc.Move(-1, c.Buf)
	c.UpdateCurLine()
	c.lastVisualX = c.GetVisualX()
}

// Right moves the cursor right one cell (if possible) or to the next line if it is at the end
func (c *Cursor) Right() {
	c.Loc = c.Loc.Move(1, c.Buf)
	c.UpdateCurLine()
	c.lastVisualX = c.GetVisualX()
}

// End moves the cursor to the end of the line it is on
func (c *Cursor) End() {
	c.x = Count(c.Buf.CurLine)
	c.UpdateCurLine()
	c.lastVisualX = c.GetVisualX()
}

// Start moves the cursor to the start of the line it is on
func (c *Cursor) Start() {
	c.x = 0
	c.UpdateCurLine()
	c.lastVisualX = c.GetVisualX()
}

func (c *Cursor) UpdateCurLine() {
	c.Buf.CurLine = string(c.Buf.Lines[c.Buf.Cursor.y])
}

// GetCharPosInLine gets the char position of a visual x y coordinate (this is necessary because tabs are 1 char but 4 visual spaces)
func (c *Cursor) GetCharPosInLine(lineNum, visualPos int) int {
	// Get the tab size
	tabSize := int(settings["tabsize"].(float64))
	// This is the visual line -- every \t replaced with the correct number of spaces
	visualLine := strings.Replace(c.Buf.Line(lineNum), "\t", "\t"+Spaces(tabSize-1), -1)
	if visualPos > Count(visualLine) {
		visualPos = Count(visualLine)
	}
	numTabs := NumOccurences(visualLine[:visualPos], '\t')
	if visualPos >= (tabSize-1)*numTabs {
		return visualPos - (tabSize-1)*numTabs
	}
	return visualPos / tabSize
}

// GetVisualX returns the x value of the cursor in visual spaces
func (c *Cursor) GetVisualX() int {
	runes := []rune(c.Buf.CurLine)
	tabSize := int(settings["tabsize"].(float64))
	return c.x + NumOccurences(string(runes[:c.x]), '\t')*(tabSize-1)
}

// Relocate makes sure that the cursor is inside the bounds of the buffer
// If it isn't, it moves it to be within the buffer's lines
func (c *Cursor) Relocate() {
	if c.y < 0 {
		c.y = 0
	} else if c.y >= c.Buf.NumLines {
		c.y = c.Buf.NumLines - 1
	}

	if c.x < 0 {
		c.x = 0
	} else if c.x > Count(c.Buf.CurLine) {
		c.x = Count(c.Buf.CurLine)
	}
}

// Display draws the cursor to the screen at the correct position
func (c *Cursor) Display(v *View) {
	// Don't draw the cursor if it is out of the viewport or if it has a selection
	if (c.y-v.Topline < 0 || c.y-v.Topline > v.height-1) || c.HasSelection() {
		screen.HideCursor()
	} else {
		screen.ShowCursor(c.GetVisualX()+v.lineNumOffset-v.leftCol, c.y-v.Topline)
	}
}
