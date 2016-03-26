package ui

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/nsf/termbox-go"
	"github.com/vhakulinen/girc/utils"
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
	currentWindow *channelOutput

	// Connections contains all connections we have so when ever you
	// create new connection remember to add it here...
	Connections    = utils.NewConnectionPool()
	channelOutputs []*channelOutput

	redrawMutex = &sync.Mutex{}
)

// InputData contains information needed to make actions acording to it
type InputData struct {
	Target  string
	Message string
	Conn    *utils.Connection
}

type channelOutput struct {
	*OutputBox
	channel string
	conn    *utils.Connection
}

func (co *channelOutput) Write(line string) {
	co.OutputBox.Write(line)
}

type writer struct {
	io.Writer
}

func (w *writer) Write(b []byte) (int, error) {
	output <- string(b)
	return len(b), nil
}

func getChannelOutput(channel string, conn *utils.Connection) (*channelOutput, bool) {
	for _, c := range channelOutputs {
		if channel == c.channel && conn == c.conn {
			return c, true
		}
	}
	return nil, false
}

// Write writes msg to channel's output
func Write(channel, msg string, conn *utils.Connection) {
	c, ok := getChannelOutput(channel, conn)
	// If there was no output for this channel, make one
	if !ok {
		w, h := termbox.Size()
		c = &channelOutput{
			OutputBox: NewOutputBox(0, 0, w, h-3, &[]string{}),
			channel:   channel,
			conn:      conn,
		}
		channelOutputs = append(channelOutputs, c)
		setTarget(channel, conn)
	}

	c.Write(msg)
	// If this channel's output is being displayed, refresh the UI
	if c == currentWindow {
		redrawAll()
	}
}

func setTarget(channel string, conn *utils.Connection) {
	c, ok := getChannelOutput(channel, conn)
	if !ok {
		return
	}
	inputBox.Target = channel
	currentWindow = c

	var name = "Du'h"
	if conn != nil {
		name = conn.Name
	}
	statusBar.SetData(fmt.Sprintf("%s - %s", channel, name))
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
					if c == currentWindow {
						if i+1 == len(channelOutputs) {
							setTarget(channelOutputs[0].channel,
								channelOutputs[0].conn)
						} else {
							setTarget(channelOutputs[i+1].channel,
								channelOutputs[i+1].conn)
						}
						break
					}
				}
				break
			case termbox.KeyEnter:
				conn := currentWindow.conn
				if conn == nil {
					pool := Connections.GetPool()
					if len(pool) > 0 {
						conn = Connections.GetPool()[0]
					}
				}
				Input <- &InputData{
					Target:  inputBox.Target,
					Message: inputBox.GetContent(),
					Conn:    conn,
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

	channelOutputs = []*channelOutput{
		&channelOutput{
			OutputBox: statusWindow,
			channel:   "status",
			conn:      nil,
		},
	}
	currentWindow = channelOutputs[0]

	setTarget("status", nil)

	redrawAll()
}

// Close closes termbox and what ever is needed to close on the UI side.
func Close() {
	termbox.Close()
	close(Input)
	close(output)
	close(end)
}
