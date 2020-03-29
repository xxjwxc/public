package mycui

import (
	"fmt"
	"math"
	"strings"

	"github.com/jroimartin/gocui"
)

// Modal struct
type Modal struct {
	*gocui.Gui
	name         string
	textArea     *textArea
	activeButton int
	buttons      []*Button
	*Attributes
	*Position
}

type textArea struct {
	*gocui.Gui
	name string
	text string
	*Attributes
	*Position
}

// NewModal new modal
func NewModal(gui *gocui.Gui, x, y, w int) *Modal {
	p := &Position{
		X: x,
		Y: y,
		W: w,
		H: y + 3,
	}

	return &Modal{
		Gui:          gui,
		name:         "modal",
		activeButton: 0,
		Attributes: &Attributes{
			textColor:   gocui.ColorWhite,
			textBgColor: gocui.ColorBlue,
		},
		Position: p,
		textArea: &textArea{
			Gui:  gui,
			name: "textArea",
			Attributes: &Attributes{
				textColor:   gocui.ColorWhite,
				textBgColor: gocui.ColorBlue,
			},
			Position: &Position{
				X: p.X + 1,
				Y: p.Y + 1,
				W: p.W - 1,
				H: p.H - 1,
			},
		},
	}
}

// SetText set text
func (m *Modal) SetText(text string) *Modal {
	m.textArea.text = text
	h := int(roundUp(float64(len(text)/(m.W-m.X)), 0))

	newLineCount := strings.Count(text, "\n")
	if newLineCount > h {
		h = newLineCount
	}

	m.textArea.H += h
	m.H += h

	return m
}

// SetTextColor set text color
func (m *Modal) SetTextColor(textColor gocui.Attribute) *Modal {
	m.textArea.textColor = textColor
	return m
}

// SetBgColor set bg color
func (m *Modal) SetBgColor(textColor gocui.Attribute) *Modal {
	m.textArea.textBgColor = textColor
	return m
}

// AddButton add button
func (m *Modal) AddButton(id, label string, key Key, handler Handler) *Button {
	var x, y, w int
	if len(m.buttons) == 0 {
		w = m.W - 5
		x = w - len(label)
		y = m.H - 1
		m.H += 2
	} else {
		p := m.buttons[len(m.buttons)-1].GetPosition()
		w = p.W - 10
		x = w - len(label)
		y = p.Y
	}

	button := NewButton(m.Gui, id, label, x, y, len(label)).
		AddHandler(gocui.KeyTab, m.nextButton).
		AddHandler(key, handler).
		SetTextColor(gocui.ColorWhite, gocui.ColorBlack).
		SetHilightColor(gocui.ColorBlack, gocui.ColorWhite)

	m.buttons = append(m.buttons, button)
	return button
}

// GetPosition get modal position
func (m *Modal) GetPosition() *Position {
	return m.Position
}

// Draw draw modal
func (m *Modal) Draw() {
	// modal
	if v, err := m.Gui.SetView(m.name, m.X, m.Y, m.W, m.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = false
		v.FgColor = m.textColor
		v.BgColor = m.textBgColor
	}

	// text area
	area := m.textArea
	if area.text != "" {
		if v, err := area.Gui.SetView(area.name, area.X, area.Y, area.W, area.H); err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}

			v.Wrap = true
			v.Frame = false

			v.FgColor = area.textColor
			v.BgColor = area.textBgColor

			fmt.Fprint(v, area.text)
		}
	}

	// button
	for _, b := range m.buttons {
		b.Draw()
	}

	if len(m.buttons) != 0 {
		m.activeButton = len(m.buttons) - 1
		m.buttons[m.activeButton].Focus()
	}
}

// Close close modal
func (m *Modal) Close() {
	if err := m.DeleteView(m.name); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	if err := m.DeleteView(m.textArea.name); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	for _, b := range m.buttons {
		b.Close()
	}
}

// nextButton focus netxt button
func (m *Modal) nextButton(g *gocui.Gui, v *gocui.View) error {
	m.buttons[m.activeButton].UnFocus()
	m.activeButton = (m.activeButton + 1) % len(m.buttons)
	m.buttons[m.activeButton].Focus()
	return nil
}

func roundUp(num, places float64) float64 {
	shift := math.Pow(10, places)
	return roundUpInt(num*shift) / shift
}

func roundUpInt(num float64) float64 {
	t := math.Trunc(num)
	return t + math.Copysign(1, num)
}
