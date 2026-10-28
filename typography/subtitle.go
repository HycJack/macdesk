package typography

import (
	"github.com/HycJack/macdesk/kit"
	"github.com/egoist/mygo/ui"
)

// Subtitle is muted text whose changes cross-fade.
func Subtitle(c *ui.Context, text string) ui.Element {
	return kit.Muted(c, text).ContentTransition(ui.ContentTransition{})
}
