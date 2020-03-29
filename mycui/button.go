package mycui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// Button button struct
type Button struct {
	*gocui.Gui
	label    string
	name     string
	handlers Handlers
	ctype    ComponentType
	*Position
	*Attributes
}

// NewButton new button
func NewButton(gui *gocui.Gui, label, name string, x, y, width int) *Button {
	if len(label) >= width {
		width = len([]rune(name)) + 1
	}

	b := &Button{
		Gui:   gui,
		label: label,
		name:  name,
		Position: &Position{
			x,
			y,
			x + width + 2,
			y + 2,
		},
		Attributes: &Attributes{
			textColor:      gocui.ColorWhite | gocui.AttrBold,
			textBgColor:    gocui.ColorBlue,
			hilightColor:   gocui.ColorBlue | gocui.AttrBold,
			hilightBgColor: gocui.ColorGreen,
		},
		handlers: make(Handlers),
		ctype:    TypeButton,
	}

	return b
}

// AddHandler add handler
func (b *Button) AddHandler(key Key, handler Handler) *Button {
	b.handlers[key] = handler
	return b
}

// SetTextColor add button fg and bg color
func (b *Button) SetTextColor(fgColor, bgColor gocui.Attribute) *Button {
	b.textColor = fgColor
	b.textBgColor = bgColor
	return b
}

// SetHilightColor add button fg and bg color
func (b *Button) SetHilightColor(fgColor, bgColor gocui.Attribute) *Button {
	b.hilightColor = fgColor
	b.hilightBgColor = bgColor
	return b
}

// GetLabel get button label
func (b *Button) GetLabel() string {
	return b.label
}

// GetPosition get button position
func (b *Button) GetPosition() *Position {
	return b.Position
}

// Focus focus to button
func (b *Button) Focus() {
	b.Gui.Cursor = false
	v, _ := b.Gui.SetCurrentView(b.label)
	v.Highlight = true
}

// UnFocus un focus
func (b *Button) UnFocus() {
	v, _ := b.Gui.View(b.label)
	v.Highlight = false
}

// GetType get component type
func (b *Button) GetType() ComponentType {
	return b.ctype
}

// Draw draw button
func (b *Button) Draw() {
	if v, err := b.Gui.SetView(b.label, b.X, b.Y, b.W, b.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Title = b.name
		v.Frame = false

		v.FgColor = b.textColor
		v.BgColor = b.textBgColor

		v.SelFgColor = b.hilightColor
		v.SelBgColor = b.hilightBgColor

		fmt.Fprint(v, fmt.Sprintf(" %s ", b.name))
	}

	if b.handlers != nil {
		for key, handler := range b.handlers {
			if err := b.Gui.SetKeybinding(b.label, key, gocui.ModNone, handler); err != nil {
				panic(err)
			}
		}
	}

}

// Close close button
func (b *Button) Close() {
	if err := b.DeleteView(b.label); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	if b.handlers != nil {
		b.DeleteKeybindings(b.label)
	}
}

// AddHandlerOnly ...
func (b *Button) AddHandlerOnly(key Key, handler Handler) {
	b.handlers[key] = handler
}
