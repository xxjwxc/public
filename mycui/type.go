package mycui

import "github.com/jroimartin/gocui"

// Key define kye type
type Key interface{}

// Handler define handler type
type Handler func(g *gocui.Gui, v *gocui.View) error

// Handlers handlers
type Handlers map[Key]Handler

// Component form component interface
type Component interface {
	GetLabel() string
	GetPosition() *Position
	GetType() ComponentType
	Focus()
	UnFocus()
	Draw()
	Close()
	AddHandlerOnly(Key, Handler)
}

// Attributes text and hilight color
type Attributes struct {
	textColor      gocui.Attribute
	textBgColor    gocui.Attribute
	hilightColor   gocui.Attribute
	hilightBgColor gocui.Attribute
}

// Position component position
type Position struct {
	X, Y int
	W, H int
}

// ComponentType component type
type ComponentType int

const (
	// TypeInputField type is input component
	TypeInputField ComponentType = iota
	// TypeSelect type is select component
	TypeSelect
	// TypeButton type is button component
	TypeButton
	// TypeCheckBox type is checkbox component
	TypeCheckBox
	// TypeRadio type is radio component
	TypeRadio
	// TypeTable type is table component
	TypeTable
)
