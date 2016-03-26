package ui

import "github.com/nsf/termbox-go"

type StatusBar struct {
	data string
	x, y int
	w    int
}

func NewStatusBar(x, y, w int) *StatusBar {
	return &StatusBar{
		data: "status bar duh'",
		x:    x, y: y, w: w,
	}
}

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
