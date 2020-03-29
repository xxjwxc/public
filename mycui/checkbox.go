package mycui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// CheckBox struct
type CheckBox struct {
	*gocui.Gui
	label     string
	isChecked bool
	box       *box
	ctype     ComponentType
	*Position
	*Attributes
	handlers Handlers
}

type box struct {
	name string
	*Position
	*Attributes
}

// NewCheckBox new checkbox
func NewCheckBox(gui *gocui.Gui, label string, x, y, labelWidth int) *CheckBox {
	if len(label) > labelWidth {
		labelWidth = len(label)
	}
	p := &Position{
		X: x,
		Y: y,
		W: x + labelWidth + 1,
		H: y + 2,
	}

	c := &CheckBox{
		Gui:       gui,
		label:     label,
		isChecked: false,
		Position:  p,
		Attributes: &Attributes{
			textColor:   gocui.ColorYellow | gocui.AttrBold,
			textBgColor: gocui.ColorDefault,
		},
		box: &box{
			name: label + "box",
			Position: &Position{
				X: p.W,
				Y: p.Y,
				W: p.W + 2,
				H: p.H,
			},
			Attributes: &Attributes{
				textColor:   gocui.ColorBlack,
				textBgColor: gocui.ColorCyan,
			},
		},
		handlers: make(Handlers),
		ctype:    TypeCheckBox,
	}

	c.handlers[gocui.KeyEnter] = c.Check
	c.handlers[gocui.KeySpace] = c.Check
	return c
}

// GetLabel get checkbox label
func (c *CheckBox) GetLabel() string {
	return c.label
}

// GetPosition get checkbox position
func (c *CheckBox) GetPosition() *Position {
	return c.box.Position
}

// Check check true or false
func (c *CheckBox) Check(g *gocui.Gui, v *gocui.View) error {
	if v.Buffer() != "" {
		v.Clear()
		c.isChecked = false
	} else {
		fmt.Fprint(v, "X")
		c.isChecked = true
	}

	return nil
}

// AddHandler add handler
func (c *CheckBox) AddHandler(key Key, handler Handler) *CheckBox {
	c.handlers[key] = handler
	return c
}

// AddAttribute add text and bg color
func (c *CheckBox) AddAttribute(textColor, textBgColor gocui.Attribute) *CheckBox {
	c.Attributes = &Attributes{
		textColor:   textColor,
		textBgColor: textBgColor,
	}

	return c
}

// IsChecked return check state
func (c *CheckBox) IsChecked() bool {
	return c.isChecked
}

// Focus focus to checkbox
func (c *CheckBox) Focus() {
	c.Gui.Cursor = true
	c.Gui.SetCurrentView(c.box.name)
}

// UnFocus unfocus
func (c *CheckBox) UnFocus() {
	c.Gui.Cursor = false
}

// GetType get component type
func (c *CheckBox) GetType() ComponentType {
	return c.ctype
}

// Draw draw label and checkbox
func (c *CheckBox) Draw() {
	// draw label
	if v, err := c.Gui.SetView(c.label, c.X, c.Y, c.W, c.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = false
		v.FgColor = c.textColor
		v.BgColor = c.textBgColor
		fmt.Fprint(v, c.label)
	}

	// draw checkbox
	b := c.box
	if v, err := c.Gui.SetView(b.name, b.X, b.Y, b.W, b.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Frame = false
		v.FgColor = b.textColor
		v.BgColor = b.textBgColor

		c.Gui.SetCurrentView(v.Name())

		for key, handler := range c.handlers {
			if err := c.Gui.SetKeybinding(v.Name(), key, gocui.ModNone, handler); err != nil {
				panic(err)
			}
		}
	}
}

// Close close checkbox
func (c *CheckBox) Close() {
	views := []string{
		c.label,
		c.box.name,
	}

	for _, v := range views {
		if err := c.DeleteView(v); err != nil {
			if err != gocui.ErrUnknownView {
				panic(err)
			}
		}
	}

	c.DeleteKeybindings(c.box.name)
}

// AddHandlerOnly add handler not retrun
func (c *CheckBox) AddHandlerOnly(key Key, handler Handler) {
	c.AddHandler(key, handler)
}
