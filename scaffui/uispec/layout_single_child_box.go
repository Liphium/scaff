package uispec

import (
	"errors"

	"github.com/Liphium/scaff/optional"
	"github.com/Liphium/scaff/scaffui"
	"github.com/Liphium/scaff/scath"
)

type SingleChildBoxSpec struct {
	Parent  scath.Constraints
	Wanted  optional.O[scath.Constraints]
	Padding scath.Padding
}

func (s SingleChildBoxSpec) ChildConstraints() (scath.Constraints, error) {
	childConstraints := s.Parent

	if wanted, ok := s.Wanted.Value(); ok {
		if !s.Parent.Fits(wanted) {
			return scath.Constraints{}, errors.New("wanted constraints could not be kept")
		}

		if childConstraints.RealMaxX() > wanted.RealMaxX() {
			childConstraints.MaxX = wanted.MaxX
			childConstraints.MinX = wanted.MinX
		}

		if childConstraints.RealMaxY() > wanted.RealMaxY() {
			childConstraints.MaxY = wanted.MaxY
			childConstraints.MinY = wanted.MinY
		}
	}

	return childConstraints.SubtractPadding(s.Padding), nil
}

func (s SingleChildBoxSpec) LayoutWithChild(child scaffui.Node) (scath.Vec, error) {
	childConstraints, err := s.ChildConstraints()
	if err != nil {
		return scath.Vec{}, err
	}

	child.SetConstraints(childConstraints)
	childSize, err := child.Layout()
	if err != nil {
		return scath.Vec{}, err
	}

	return childSize.AddPadding(s.Padding), nil
}

func (s SingleChildBoxSpec) LayoutWithoutChild() (scath.Vec, error) {
	size := scath.Vec{}
	if wanted, ok := s.Wanted.Value(); ok {
		var err error
		size, err = wanted.TakeMaxWithin(s.Parent)
		if err != nil {
			return scath.Vec{}, err
		}
	}

	return size.AddPadding(s.Padding), nil
}
