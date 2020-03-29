package mycui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

// Select struct
type Select struct {
	*InputField
	options      []string
	currentOpt   int
	isExpanded   bool
	ctype        ComponentType
	listColor    *Attributes
	listHandlers Handlers
}

// NewSelect new select
func NewSelect(gui *gocui.Gui, label, name string, x, y, labelWidth, fieldWidth int) *Select {

	s := &Select{
		InputField:   NewInputField(gui, label, name, x, y, labelWidth, fieldWidth),
		listHandlers: make(Handlers),
		ctype:        TypeSelect,
	}

	s.AddHandler(gocui.KeyEnter, s.expandOpt)
	s.AddAttribute(gocui.ColorBlack, gocui.ColorWhite, gocui.ColorBlack, gocui.ColorGreen).
		AddListHandler('j', s.nextOpt).
		AddListHandler('k', s.preOpt).
		AddListHandler(gocui.KeyArrowDown, s.nextOpt).
		AddListHandler(gocui.KeyArrowUp, s.preOpt).
		AddListHandler(gocui.KeyEnter, s.selectOpt).
		SetEditable(false)

	return s
}

// AddOptions add select options
func (s *Select) AddOptions(opts ...string) *Select {
	for _, opt := range opts {
		s.options = append(s.options, opt)
	}
	return s
}

// AddOption add select option
func (s *Select) AddOption(opt string) *Select {
	s.options = append(s.options, opt)
	return s
}

// AddAttribute add select attribute
func (s *Select) AddAttribute(textColor, textBgColor, fgColor, bgColor gocui.Attribute) *Select {
	s.listColor = &Attributes{
		textColor:      textColor,
		textBgColor:    textBgColor,
		hilightColor:   fgColor,
		hilightBgColor: bgColor,
	}

	return s
}

// AddListHandler add list handler
func (s *Select) AddListHandler(key Key, handler Handler) *Select {
	s.listHandlers[key] = handler
	return s
}

// GetSelected get selected option
func (s *Select) GetSelected() string {
	return s.options[s.currentOpt]
}

// SetSelected set the default selected
func (s *Select) SetSelected(str string) {
	for i := 0; i < len(s.options); i++ {
		if strings.EqualFold(s.options[i], str) {
			s.currentOpt = i
			break
		}
	}
}

// Focus set focus to select
func (s *Select) Focus() {
	s.Gui.Cursor = true
	s.Gui.SetCurrentView(s.GetLabel())
}

// UnFocus un focus
func (s *Select) UnFocus() {
	s.Gui.Cursor = false
}

// GetType get component type
func (s *Select) GetType() ComponentType {
	return s.ctype
}

// Close close select
func (s *Select) Close() {
	s.InputField.Close()
	if s.isExpanded {
		for _, opt := range s.options {
			s.DeleteView(opt)
			s.DeleteKeybindings(opt)
		}
	}
}

// Draw draw select
func (s *Select) Draw() {
	if len(s.options) > 0 {
		s.InputField.SetText(s.options[s.currentOpt])
	}
	s.InputField.Draw()
}

func (s *Select) nextOpt(g *gocui.Gui, v *gocui.View) error {
	maxOpt := len(s.options)
	if maxOpt == 0 {
		return nil
	}

	v.Highlight = false

	next := s.currentOpt + 1
	if next >= maxOpt {
		next = s.currentOpt
	}

	s.currentOpt = next
	v, _ = g.SetCurrentView(s.options[next])

	v.Highlight = true

	return nil
}

func (s *Select) preOpt(g *gocui.Gui, v *gocui.View) error {
	maxOpt := len(s.options)
	if maxOpt == 0 {
		return nil
	}

	v.Highlight = false

	next := s.currentOpt - 1
	if next < 0 {
		next = 0
	}

	s.currentOpt = next
	v, _ = g.SetCurrentView(s.options[next])

	v.Highlight = true

	return nil
}

func (s *Select) selectOpt(g *gocui.Gui, v *gocui.View) error {
	if !s.isExpanded {
		s.expandOpt(g, v)
	} else {
		s.closeOpt(g, v)
	}

	return nil
}

func (s *Select) expandOpt(g *gocui.Gui, vi *gocui.View) error {
	if s.hasOpts() {
		s.isExpanded = true
		g.Cursor = false

		x := s.field.X
		w := s.field.W

		y := s.field.Y
		h := y + 2

		for _, opt := range s.options {
			y++
			h++
			if v, err := g.SetView(opt, x, y, w, h); err != nil {
				if err != gocui.ErrUnknownView {
					panic(err)
				}

				v.Frame = false
				v.SelFgColor = s.listColor.textColor
				v.SelBgColor = s.listColor.textBgColor
				v.FgColor = s.listColor.hilightColor
				v.BgColor = s.listColor.hilightBgColor

				for key, handler := range s.listHandlers {
					if err := g.SetKeybinding(v.Name(), key, gocui.ModNone, handler); err != nil {
						panic(err)
					}
				}

				fmt.Fprint(v, opt)
			}

		}

		v, _ := g.SetCurrentView(s.options[s.currentOpt])
		v.Highlight = true
	}

	return nil
}

func (s *Select) closeOpt(g *gocui.Gui, v *gocui.View) error {
	s.isExpanded = false
	g.Cursor = true

	for _, opt := range s.options {
		g.DeleteView(opt)
		g.DeleteKeybindings(opt)
	}

	v, _ = g.SetCurrentView(s.GetLabel())

	v.Clear()

	fmt.Fprint(v, s.GetSelected())

	return nil
}

func (s *Select) hasOpts() bool {
	if len(s.options) > 0 {
		return true
	}
	return false
}
