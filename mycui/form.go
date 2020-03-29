package mycui

import (
	"github.com/jroimartin/gocui"
)

// Form form struct
type Form struct {
	*gocui.Gui
	activeItem  int
	activeRadio int
	id          string
	name        string
	inputs      []*InputField
	checkBoxs   []*CheckBox
	buttons     []*Button
	selects     []*Select
	//radios      []*Radio
	components []Component
	closeFuncs []func() error
	*Position
}

// FormData form data struct
type FormData struct {
	inputs    map[string]string
	checkBoxs map[string]bool
	selects   map[string]string
	radio     map[string]string
}

// NewForm new form
func NewForm(gui *gocui.Gui, id, name string, x, y, w, h int) *Form {
	f := &Form{
		Gui:        gui,
		activeItem: 0,
		id:         id,
		name:       name,
		Position: &Position{
			X: x,
			Y: y,
			W: x + w,
			H: y + h,
		},
	}

	return f
}

// AddInputField add input field to form
func (f *Form) AddInputField(label, name string, labelWidth, fieldWidth int) *InputField {
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		y = p.H
	} else {
		y = f.Y + 1
	}

	input := NewInputField(
		f.Gui,
		label,
		name,
		f.X+1,
		y,
		labelWidth,
		fieldWidth,
	)

	f.inputs = append(f.inputs, input)
	f.components = append(f.components, input)

	return input
}

// AddButton add button to form
func (f *Form) AddButton(id, label string, handler Handler) *Button {
	var x int
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		if f.isButtonLastView() {
			x = p.W
			y = p.Y - 1
		} else {
			x = f.X
			y = p.H
		}
	} else {
		x = f.X
		y = f.Y
	}

	button := NewButton(
		f.Gui,
		id,
		label,
		x+1,
		y+1,
		len([]rune(label))+1,
	)

	button.AddHandler(gocui.KeyEnter, handler)

	f.buttons = append(f.buttons, button)
	f.components = append(f.components, button)

	return button
}

// AddCheckBox add checkbox
func (f *Form) AddCheckBox(label string, width int) *CheckBox {
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		y = p.H
	} else {
		y = f.Y
	}

	checkbox := NewCheckBox(
		f.Gui,
		label,
		f.X+1,
		y,
		width,
	)

	f.checkBoxs = append(f.checkBoxs, checkbox)
	f.components = append(f.components, checkbox)

	return checkbox
}

// AddSelect add select
func (f *Form) AddSelect(label, name string, labelWidth, listWidth int) *Select {
	var y int

	p := f.getLastViewPosition()
	if p != nil {
		y = p.H
	} else {
		y = f.Y
	}

	Select := NewSelect(
		f.Gui,
		label,
		name,
		f.X+1,
		y,
		labelWidth,
		listWidth,
	)

	f.selects = append(f.selects, Select)
	f.components = append(f.components, Select)

	return Select
}

// AddCloseFunc add close function
func (f *Form) AddCloseFunc(function func() error) {
	f.closeFuncs = append(f.closeFuncs, function)
}

// GetFieldTexts form data
func (f *Form) GetFieldTexts() map[string]string {
	data := make(map[string]string)

	if len(f.inputs) == 0 {
		return data
	}

	for _, item := range f.inputs {
		data[item.GetLabel()] = item.GetFieldText()
	}

	return data
}

// GetFieldText get form data with field name
func (f *Form) GetFieldText(target string) string {
	return f.GetFieldTexts()[target]
}

// GetCheckBoxStates get checkbox states
func (f *Form) GetCheckBoxStates() map[string]bool {
	state := make(map[string]bool)

	if len(f.checkBoxs) == 0 {
		return state
	}

	for _, box := range f.checkBoxs {
		state[box.GetLabel()] = box.IsChecked()
	}

	return state
}

// GetCheckBoxState get checkbox states
func (f *Form) GetCheckBoxState(target string) bool {
	return f.GetCheckBoxStates()[target]
}

// GetSelectedOpts get selected options
func (f *Form) GetSelectedOpts() map[string]string {
	opts := make(map[string]string)

	if len(f.selects) == 0 {
		return opts
	}

	for _, Select := range f.selects {
		opts[Select.GetLabel()] = Select.GetSelected()
	}

	return opts
}

// GetSelectedOpt get selected options
func (f *Form) GetSelectedOpt(target string) string {
	return f.GetSelectedOpts()[target]
}

// GetSelectedRadio get selected radio
func (f *Form) GetSelectedRadio(target string) string {
	return f.GetSelectedOpts()[target]
}

// GetFormData get form data
func (f *Form) GetFormData() *FormData {
	fd := &FormData{
		inputs:    f.GetFieldTexts(),
		checkBoxs: f.GetCheckBoxStates(),
		selects:   f.GetSelectedOpts(),
	}

	return fd
}

// GetInputs get inputs
func (f *Form) GetInputs() []*InputField {
	return f.inputs
}

// GetCheckBoxs get checkboxs
func (f *Form) GetCheckBoxs() []*CheckBox {
	return f.checkBoxs
}

// GetButtons get buttons
func (f *Form) GetButtons() []*Button {
	return f.buttons
}

// GetSelects get selects
func (f *Form) GetSelects() []*Select {
	return f.selects
}

// GetItems get items
func (f *Form) GetItems() []Component {
	return f.components
}

// SetCurrentItem set current item index
func (f *Form) SetCurrentItem(index int) *Form {
	f.activeItem = index
	f.components[index].Focus()
	return f
}

// GetCurrentItem get current item index
func (f *Form) GetCurrentItem() int {
	return f.activeItem
}

// Validate validate form items
func (f *Form) Validate() bool {
	isValid := true
	for _, item := range f.inputs {
		if !item.Validate() {
			isValid = false
		}
	}

	return isValid
}

// NextItem to next item
func (f *Form) NextItem(g *gocui.Gui, v *gocui.View) error {
	f.components[f.activeItem].UnFocus()
	f.activeItem = (f.activeItem + 1) % len(f.components)
	f.components[f.activeItem].Focus()
	return nil
}

// PreItem to pre item
func (f *Form) PreItem(g *gocui.Gui, v *gocui.View) error {
	f.components[f.activeItem].UnFocus()

	if f.activeItem-1 < 0 {
		f.activeItem = len(f.components) - 1
	} else {
		f.activeItem = (f.activeItem - 1) % len(f.components)
	}

	f.components[f.activeItem].Focus()

	return nil
}

// Draw form
func (f *Form) Draw() {
	if v, err := f.Gui.SetView(f.id, f.X, f.Y, f.W+1, f.H+1); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		v.Title = f.name
	}

	for _, cp := range f.components {
		p := cp.GetPosition()
		if p.W > f.W {
			f.W = p.W
		}
		if p.H > f.H {
			f.H = p.H
		}
		cp.AddHandlerOnly(gocui.KeyTab, f.NextItem)
		cp.AddHandlerOnly(gocui.KeyArrowDown, f.NextItem)
		cp.AddHandlerOnly(gocui.KeyArrowUp, f.PreItem)
		cp.Draw()
	}

	f.SetView(f.id, f.X, f.Y, f.W+1, f.H+1)

	if len(f.components) != 0 {
		f.components[0].Focus()
	}
}

// Close close form
func (f *Form) Close(g *gocui.Gui, v *gocui.View) error {
	if err := f.Gui.DeleteView(f.id); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}
	}

	for _, c := range f.components {
		c.Close()
	}

	if len(f.closeFuncs) != 0 {
		for _, f := range f.closeFuncs {
			f()
		}
	}

	return nil
}

func (f *Form) getLastViewPosition() *Position {
	cpl := len(f.components)
	if cpl == 0 {
		return nil
	}

	return f.components[cpl-1].GetPosition()
}

func (f *Form) isButtonLastView() bool {
	cpl := len(f.components)
	if cpl == 0 {
		return false
	}

	return f.components[cpl-1].GetType() == TypeButton
}
