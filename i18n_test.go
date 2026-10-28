package main

import (
	"testing"

	"github.com/HycJack/macdesk/icon"
)

// TestCustomIcons pins the hand-drawn SVGs into the icon store.
func TestCustomIcons(t *testing.T) {
	for _, name := range []string{"airdrop", "applelogo"} {
		icon.Must(name)
	}
}

func TestPhrasesParity(t *testing.T) {
	if len(zh.finderGroups) != len(en.finderGroups) || len(zh.finderGroups) == 0 {
		t.Fatal("finder groups differ")
	}
	for i := range zh.finderGroups {
		if len(zh.finderGroups[i].Items) != len(en.finderGroups[i].Items) {
			t.Fatalf("finder group %d items differ", i)
		}
	}
	if len(zh.noteGroups) != len(en.noteGroups) || len(zh.noteGroups) == 0 {
		t.Fatal("note groups differ")
	}
	if len(zh.noteBodies) != len(zh.noteGroups[0].Items) ||
		len(en.noteBodies) != len(en.noteGroups[0].Items) {
		t.Fatal("note bodies do not match notes")
	}
	if len(zh.finderFiles) != len(en.finderFiles) || len(zh.finderFiles) == 0 {
		t.Fatal("finder files differ")
	}
	if zh.langToggle == en.langToggle {
		t.Fatal("language toggles must differ")
	}
}

func TestToggleLang(t *testing.T) {
	d := newDesk()
	d.finderIn, d.noteIn = "AirDrop", "会议记录"
	d.toggleLang()
	if d.lang != "en" || d.finderIn != "AirDrop" || d.noteIn != "Meeting notes" {
		t.Fatalf("after toggle: lang %q finder %q note %q", d.lang, d.finderIn, d.noteIn)
	}
	if d.noteText != en.noteBodies[0] {
		t.Errorf("note body did not follow: %q", d.noteText)
	}
	d.toggleLang()
	if d.lang != "zh" || d.finderIn != "AirDrop" || d.noteIn != "会议记录" {
		t.Fatalf("after toggle back: lang %q finder %q note %q", d.lang, d.finderIn, d.noteIn)
	}
	if d.noteText != zh.noteBodies[0] {
		t.Errorf("note body did not follow back: %q", d.noteText)
	}
}
