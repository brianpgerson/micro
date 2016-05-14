package main

import (
	"io/ioutil"
)

// Buffer stores the text for files that are loaded into the text editor
// It uses a rope to efficiently store the string and contains some
// simple functions for saving and wrapper functions for modifying the rope
type Buffer struct {
	LineArray

	// Path to the file on disk
	Path string
	// Name of the buffer on the status line
	Name string

	IsModified bool

	// Provide efficient and easy access to text and lines so the rope String does not
	// need to be constantly recalculated
	// These variables are updated in the update() function
	Lines    []string
	NumLines int
	NumChars int

	// Syntax highlighting rules
	rules []SyntaxRule
	// The buffer's filetype
	FileType string
}

// NewBuffer creates a new buffer from `txt` with path and name `path`
func NewBuffer(txt, path string) *Buffer {
	b := new(Buffer)
	b.Path = path
	b.Name = path

	b.Update()
	b.UpdateRules()

	return b
}

// UpdateRules updates the syntax rules and filetype for this buffer
// This is called when the colorscheme changes
func (b *Buffer) UpdateRules() {
	b.rules, b.FileType = GetRules(b)
}

// Update fetches the string from the rope and updates the `text` and `lines` in the buffer
func (b *Buffer) Update() {
	b.NumLines = len(b.Lines)
}

// Save saves the buffer to its default path
func (b *Buffer) Save() error {
	return b.SaveAs(b.Path)
}

// SaveAs saves the buffer to a specified path (filename), creating the file if it does not exist
func (b *Buffer) SaveAs(filename string) error {
	b.UpdateRules()
	data := []byte(b.String())
	err := ioutil.WriteFile(filename, data, 0644)
	if err == nil {
		b.IsModified = false
	}
	return err
}

// Insert a string into the rope
func (b *Buffer) Insert(idx Loc, value string) {
	b.IsModified = true
	b.NumChars += len(value)
	b.LineArray.Insert(idx, []byte(value))
	b.Update()
}

// Remove a slice of the rope from start to end (exclusive)
// Returns the string that was removed
func (b *Buffer) Remove(start, end Loc) string {
	b.IsModified = true
	if start.LessThan(b.Start()) {
		start = b.Start()
	}
	if end.GreaterThan(b.End()) {
		end = b.End()
	}
	removed := b.LineArray.Delete(start, end)
	b.NumChars -= len(removed)
	b.Update()
	return removed
}

// Start returns the location of the first character in the buffer
func (b *Buffer) Start() Loc {
	return Loc{0, 0}
}

// End returns the location of the last character in the buffer
func (b *Buffer) End() Loc {
	return Loc{len(b.Lines[len(b.Lines)-1]) - 1, b.NumLines}
}

// Len gives the length of the buffer
func (b *Buffer) Len() int {
	var sum int
	for _, l := range b.Lines {
		sum += len(l)
	}
	return sum
}
