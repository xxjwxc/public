package mycui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// Validate validate struct
type Validate struct {
	ErrMsg string
	Do     func(value string) bool
}

// Validator validate struct
type Validator struct {
	*gocui.Gui
	name      string
	errMsg    string
	isValid   bool
	validates []Validate
	*Position
}

// NewValidator new validator
func NewValidator(gui *gocui.Gui, name string, x, y, w, h int) *Validator {
	return &Validator{
		Gui:     gui,
		name:    name,
		isValid: true,
		Position: &Position{
			X: x,
			Y: y,
			W: w,
			H: h,
		},
	}
}

// AddValidate add validate
func (v *Validator) AddValidate(errMsg string, validate func(value string) bool) {
	v.validates = append(v.validates, Validate{
		ErrMsg: errMsg,
		Do:     validate,
	})

	if v.X+len(errMsg) > v.W {
		v.W += len(errMsg)
	}
}

// DispValidateMsg display validate error message
func (v *Validator) DispValidateMsg() {
	if vi, err := v.SetView(v.name, v.X, v.Y, v.W, v.H); err != nil {
		if err != gocui.ErrUnknownView {
			panic(err)
		}

		vi.Frame = false
		vi.BgColor = gocui.ColorDefault
		vi.FgColor = gocui.ColorRed

		fmt.Fprint(vi, v.errMsg)
	}
}

// CloseValidateMsg close validate error message
func (v *Validator) CloseValidateMsg() {
	v.DeleteView(v.name)
}

// IsValid if valid return true
func (v *Validator) IsValid() bool {
	return v.isValid
}

// Validate validate value
func (v *Validator) Validate(value string) {
	for _, validate := range v.validates {
		if !validate.Do(value) {
			v.errMsg = validate.ErrMsg
			v.isValid = false
			v.DispValidateMsg()
			break
		} else {
			v.isValid = true
			v.CloseValidateMsg()
		}
	}
}
