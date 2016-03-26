package ui

import "github.com/nsf/termbox-go"

// StatusBar is simple one line height status line.
type StatusBar struct {
	data string
	x, y int
	w    int
}

// NewStatusBar creates new statusbar with empty data on give position.
func NewStatusBar(x, y, w int) *StatusBar {
	return &StatusBar{
		data: "",
		x:    x, y: y, w: w,
	}
}

// SetData sets the displayed data for this statusbar.
func (sb *StatusBar) SetData(data string) {
	sb.data = data
}

// Draw draws the statusbar using termbox. Does not call termbox.Flush()
func (sb *StatusBar) Draw() {
	l := len(sb.data)
	for i := 0; i < sb.w; i++ {
		var ch = ' '
		if i < l {
			ch = rune(sb.data[i])
		}
		termbox.SetCell(sb.x+i, sb.y, ch, termbox.ColorDefault, termbox.ColorGreen)
	}
}
