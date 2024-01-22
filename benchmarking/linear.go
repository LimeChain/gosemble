package benchmarking

import (
	"errors"
	"fmt"
	"math"
)

var (
	errInvalidValues = errors.New("failed to initialize new linear component: linear component min value must be less than or equal to max value")
)

type linear struct {
	min, max uint32
	value    uint32
}

// Used by benchmark tests to specify that a benchmarking variable is linear
// over some specified range, i.e. `NewLinear(0, 1000)` means that the corresponding variable
// is allowed to range from `0` to `1000`, inclusive.
func NewLinear(min, max uint32) (*linear, error) {
	if max < min {
		return nil, errInvalidValues
	}

	// max is the default linear value
	// learn more about linears and benchmark exxecution in substrate:
	// https://paritytech.github.io/polkadot-sdk/master/frame_benchmarking/v2/index.html
	// https://docs.substrate.io/test/benchmark/
	return &linear{min: min, max: max, value: max}, nil
}

// Component value, modified for each step iteration
func (l *linear) Value() uint32 {
	return l.value
}

// Internal function used to set the linear value for each step iteration
func (l *linear) setValue(value uint32) {
	l.value = value
}

// Internal function that calculates linear values for given steps
func (l *linear) values(steps int) ([]uint32, error) {
	stepSize := math.Max(float64(l.max-l.min)/float64(steps-1), 0)

	values := make([]uint32, steps)
	for step := 0; step < steps; step++ {
		stepValue := uint32(float64(l.min) + stepSize*float64(step))
		if stepValue < l.min || stepValue > l.max {
			return []uint32{}, fmt.Errorf("failed to generate values for linear(min: %d, max: %d) at step %d/%d: linear value must not be lower than linear.min or higher than linear.max", l.min, l.max, step, steps)
		}
		values[step] = stepValue
	}

	return values, nil
}

// Internal function that returns all component values
func componentValues(linearComponents []*linear) []uint32 {
	values := make([]uint32, len(linearComponents))
	for i, l := range linearComponents {
		values[i] = l.Value()
	}

	return values
}
