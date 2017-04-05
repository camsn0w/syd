package main

import (
	"log"
	"os"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"

	"github.com/edsrzf/mmap-go"
	"github.com/mibk/syd/core"
	"github.com/mibk/syd/pkg/undo"
	"github.com/mibk/syd/ui"
	"github.com/mibk/syd/ui/term"
	"github.com/mibk/syd/vi"
)

var (
	UI       = &term.UI{}
	filename = ""
)

func main() {
	log.SetPrefix("syd: ")
	log.SetFlags(0)
	if err := UI.Init(); err != nil {
		log.Fatalln("initializing ui:", err)
	}
	defer UI.Close()

	var b []byte
	if len(os.Args) > 1 {
		filename = os.Args[1]
		m, err := readFile(filename)
		if err != nil {
			panic(err)
		}
		defer m.Unmap()
		b = []byte(m)
	}
	buf := undo.NewBuffer(b)

	win := UI.NewWindow()
	ed := &Editor{
		events:    make(chan ui.Event),
		vi:        vi.NewParser(),
		activeWin: core.NewWindow(win, core.NewUndoBuffer(buf)),
	}
	ed.activeWin.SetFilename(filename)
	ed.Main()
}

func readFile(filename string) (mmap.MMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	m, err := mmap.Map(f, 0, 0)
	if err != nil {
		return nil, err
	}
	return m, nil
}

const (
	ModeNormal = iota
	ModeInsert
)

type Editor struct {
	events     chan ui.Event
	vi         *vi.Parser
	shouldQuit bool

	activeWin *core.Window
	mode      int
}

func (ed *Editor) Main() {
	for !ed.shouldQuit {
		ed.activeWin.Render()
		ev := <-ui.Events
		if ev == ui.Quit {
			return
		}
		switch ev := ev.(type) {
		case key.Event:
			UI.Push_Key_Event(ev)
		case mouse.Event:
			// Temporary reasons...
			UI.Push_Mouse_Event(ev)
		}
	}
}
