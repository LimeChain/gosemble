package benchmarking

import (
	"errors"
	"math"
	"strings"
)

var (
	errInvalidValues = errors.New("failed to initialize new linear component: linear component min value must be less than or equal to max value")
	errEmptyName     = errors.New("failed to initialize new linear component: name must not be empty")
)

type linear struct {
	name     string
	min, max uint32
	value    uint32
}

// Used by benchmark tests to specify that a benchmarking variable is linear
// over some specified range, i.e. `NewLinear("a", 0, 1000)` means that the corresponding variable
// is allowed to range from `0` to `1000`, inclusive and generated weight file will accept arg named "a"
func NewLinear(name string, min, max uint32) (*linear, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return nil, errEmptyName
	}

	if max < min {
		return nil, errInvalidValues
	}

	// max is the default linear value
	return &linear{
		name:  name,
		min:   min,
		max:   max,
		value: max,
	}, nil
}

func (l *linear) Name() string {
	return l.name
}

// Component value, modified for each step iteration
func (l *linear) Value() uint32 {
	return l.value
}

// Internal function used to set the linear value for each step iteration
func (l *linear) setValue(value uint32) {
	l.value = value
}

// Internal function that calculates linear values for given range and execution steps
func (l *linear) rangeValues(steps int) []uint32 {
	stepSize := math.Max(float64(l.max-l.min)/float64(steps-1), 0)

	values := make([]uint32, steps)
	for step := 0; step < steps; step++ {
		stepValue := uint32(float64(l.min) + stepSize*float64(step))
		values[step] = stepValue
	}

	return values
}

// Internal function used for dereferencing the original components before passing them to analysis
func copyComponents(components []*linear) []linear {
	copiedComponents := make([]linear, len(components))
	for i, c := range components {
		copiedComponents[i] = *c
	}
	return copiedComponents
}

// Internal function that returns the values of component array
func componentValues(components []linear) []uint32 {
	values := make([]uint32, len(components))
	for i, l := range components {
		values[i] = l.Value()
	}
	return values
}
