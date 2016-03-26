package ui

import (
	"io"
	"log"
	"sync"

	"github.com/nsf/termbox-go"
)

var (
	// Writer writes data to status window. Is thread safe because
	// uses channel to pass the data along to output loop
	Writer = &writer{}

	// Input passes input from UI forward. Listen this on main.
	Input  chan *InputData
	output chan string
	end    chan bool

	inputBox  *InputBox
	statusBar *StatusBar

	statusWindow  *OutputBox
	currentWindow *OutputBox

	channelOutputs []*ChannelOutput

	redrawMutex = &sync.Mutex{}
)

// InputData contains information needed to make actions acording to it
type InputData struct {
	Target  string
	Message string
}

// ChannelOutput .
type ChannelOutput struct {
	*OutputBox
	channel string
}

func (co *ChannelOutput) Write(line string) {
	co.OutputBox.Write(line)
}

type writer struct {
	io.Writer
}

func (w *writer) Write(b []byte) (int, error) {
	output <- string(b)
	return len(b), nil
}

func getChannelOutput(channel string) (*ChannelOutput, bool) {
	for _, c := range channelOutputs {
		if channel == c.channel {
			return c, true
		}
	}
	return nil, false
}

// Write writes msg to channel's output
func Write(channel, msg string) {
	c, ok := getChannelOutput(channel)
	// If there was no output for this channel, make one
	if !ok {
		w, h := termbox.Size()
		c = &ChannelOutput{
			OutputBox: NewOutputBox(0, 0, w, h-3, &[]string{}),
			channel:   channel,
		}
		channelOutputs = append(channelOutputs, c)
		setTarget(channel)
	}

	c.Write(msg)
	// If this channel's output is being displayed, refresh the UI
	if c.OutputBox == currentWindow {
		redrawAll()
	}
}

func setTarget(channel string) {
	c, ok := getChannelOutput(channel)
	if !ok {
		return
	}
	inputBox.Target = channel
	currentWindow = c.OutputBox
}

func outputLoop(wg *sync.WaitGroup) {
loop:
	for {
		select {
		case msg, ok := <-output:
			if !ok {
				log.Println("Failed to read data")
				break
			}
			// write to status window
			statusWindow.Write(msg)
			break
		case <-end:
			break loop
		}
		redrawAll()
	}
	log.Println("Output loop exited")
	wg.Done()
}

func inputLoop(wg *sync.WaitGroup) {
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				end <- true
				break loop
			case termbox.KeyCtrlN:
				for i, c := range channelOutputs {
					if c.OutputBox == currentWindow {
						if i+1 == len(channelOutputs) {
							setTarget(channelOutputs[0].channel)
						} else {
							setTarget(channelOutputs[i+1].channel)
						}
						break
					}
				}
				break
			case termbox.KeyEnter:
				Input <- &InputData{
					Target:  inputBox.Target,
					Message: inputBox.GetContent(),
				}
				inputBox.Clear()
				break
			case termbox.KeyTab:
				inputBox.InsertRune('\t')
				break
			case termbox.KeySpace:
				inputBox.InsertRune(' ')
				break
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				inputBox.RemoveRuneBackwards()
				break
			default:
				inputBox.InsertRune(ev.Ch)
				break
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redrawAll()
	}
	log.Println("Input loop exited")
	wg.Done()
}

// Loop runs the main loop and blocks untill it exists.
func Loop() {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go outputLoop(wg)
	go inputLoop(wg)
	wg.Wait()
}

func redrawAll() {
	redrawMutex.Lock()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	currentWindow.Draw()
	inputBox.Draw()
	statusBar.Draw()

	termbox.Flush()
	redrawMutex.Unlock()
}

// Init initializes termbox and the whole UI (du'h)
func Init() {
	err := termbox.Init()
	if err != nil {
		log.Fatalf("Failed to initialized termbox: %v", err)
	}
	termbox.SetInputMode(termbox.InputEsc)
	termbox.SetOutputMode(termbox.OutputNormal)

	Input = make(chan *InputData)
	output = make(chan string)
	end = make(chan bool)

	w, h := termbox.Size()
	statusWindow = NewOutputBox(0, 0, w, h-3, &[]string{})
	inputBox = NewInputBox(0, h-1)
	statusBar = NewStatusBar(0, h-2, w)

	channelOutputs = []*ChannelOutput{
		&ChannelOutput{
			OutputBox: statusWindow,
			channel:   "",
		},
	}
	currentWindow = channelOutputs[0].OutputBox

	setTarget("status")

	redrawAll()
}

// Close closes termbox and what ever is needed to close on the UI side.
func Close() {
	termbox.Close()
	close(Input)
	close(output)
	close(end)
}
