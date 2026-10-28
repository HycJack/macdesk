// Command macdesk shows a macOS-style desktop with an app store.
package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func main() {
	d := newDesk()
	mygo.App.WhenReady(func() {
		mygo.NewWindow(mygo.WindowOptions{
			Title:     "macdesk",
			Width:     1280,
			Height:    820,
			MinWidth:  960,
			MinHeight: 600,
			Content:   ui.View(d.view),
		})
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
