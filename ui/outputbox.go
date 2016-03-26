package ui

import (
	"io"

	"github.com/nsf/termbox-go"
)

// OutputBox is used as output window for channels/status window
type OutputBox struct {
	io.Writer
	data       *[]string
	x, y, w, h int
}

// NewOutputBox creates new outputbox with initial data "data"
func NewOutputBox(x, y, w, h int, data *[]string) *OutputBox {
	return &OutputBox{
		data: data,

		x: x, y: y, w: w, h: h,
	}
}

func (ob *OutputBox) Write(line []byte) (int, error) {
	*ob.data = append(*ob.data, string(line))
	redrawAll()
	return len(line), nil
}

// Draw draws OutputBox using termbox. Doesn't call termbox.Flush()
func (ob *OutputBox) Draw() {
	l := len(*ob.data)
	var display = make([]string, l)
	if l > ob.h {
		for i, line := range (*ob.data)[l-ob.h:] {
			display[i] = line
		}
	} else {
		display = *ob.data
	}
	for y, line := range display {
		for x, ch := range line {
			termbox.SetCell(ob.x+x, ob.y+y, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}
