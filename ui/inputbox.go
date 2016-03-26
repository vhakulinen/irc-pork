package ui

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type InputBox struct {
	Target string
	Prompt string
	data   string

	x int
	y int
}

func NewInputBox(x, y int) *InputBox {
	return &InputBox{
		Target: "status",
		Prompt: ">> ",
		data:   "",

		x: x, y: y,
	}
}

func (ib *InputBox) GetContent() string {
	return ib.data
}

func (ib *InputBox) InsertRune(ch rune) {
	ib.data += string(ch)
}

func (ib *InputBox) RemoveRuneBackwards() {
	len := len(ib.data)
	if len == 0 {
		return
	}
	ib.data = ib.data[:len-1]
}

func (ib *InputBox) Clear() {
	for i, _ := range ib.getDisplayContent() {
		termbox.SetCell(ib.x+i, ib.y, ' ', termbox.ColorDefault, termbox.ColorDefault)
	}
	ib.data = ""
}

func (ib *InputBox) getDisplayContent() string {
	return fmt.Sprintf("[%s]%s %s", ib.Target, ib.Prompt, ib.data)
}

func (ib *InputBox) Draw() {
	data := ib.getDisplayContent()
	for i, ch := range data {
		termbox.SetCell(ib.x+i, ib.y, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.SetCursor(ib.x+len(data), ib.y)
}
