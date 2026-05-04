package uispec

import (
	"errors"

	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
)

type SingleChildBoxSpec struct {
	Parent  scaffui.Constraints
	Wanted  optional.O[scaffui.Constraints]
	Padding scaffui.Padding
}

func (s SingleChildBoxSpec) ChildConstraints() (scaffui.Constraints, error) {
	childConstraints := s.Parent

	if wanted, ok := s.Wanted.Value(); ok {
		if !s.Parent.Fits(wanted) {
			return scaffui.Constraints{}, errors.New("wanted constraints could not be kept")
		}

		if childConstraints.RealMaxWidth() > wanted.RealMaxWidth() {
			childConstraints.MaxWidth = wanted.MaxWidth
			childConstraints.MinWidth = wanted.MinWidth
		}

		if childConstraints.RealMaxHeight() > wanted.RealMaxHeight() {
			childConstraints.MaxHeight = wanted.MaxHeight
			childConstraints.MinHeight = wanted.MinHeight
		}
	}

	return childConstraints.SubtractPadding(s.Padding), nil
}

func (s SingleChildBoxSpec) LayoutWithChild(child scaffui.Node) (scaffui.Size, error) {
	childConstraints, err := s.ChildConstraints()
	if err != nil {
		return scaffui.Size{}, err
	}

	child.SetConstraints(childConstraints)
	childSize, err := child.Layout()
	if err != nil {
		return scaffui.Size{}, err
	}

	return childSize.AddPadding(s.Padding), nil
}

func (s SingleChildBoxSpec) LayoutWithoutChild() (scaffui.Size, error) {
	size := scaffui.Size{}
	if wanted, ok := s.Wanted.Value(); ok {
		var err error
		size, err = wanted.TakeMaxWithin(s.Parent)
		if err != nil {
			return scaffui.Size{}, err
		}
	}

	return size.AddPadding(s.Padding), nil
}
