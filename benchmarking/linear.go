package benchmarking

import (
	"fmt"
	"math"
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
		return nil, fmt.Errorf("failed to initialize new linear component: linear component min value must be less than or equal to max value")
	}

	return &linear{min: min, max: max}, nil
}

// Current value for linear component. Modified for each step before testFn execution
func (l *linear) Value() uint32 {
	return l.value
}

// Internal function that sets the linear value before executing testFn
func (l *linear) setValue(value uint32) {
	l.value = value
}

// Internal function that calculates linear values for given steps parameter
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

// Internal function that calculates step values for provided linear components and steps and executes executeFn for each step
func forEachStep(steps int, components *[]*linear, executeFn func(currentStep int, currentComponentIndex int)) error {
	if len(*components) == 0 {
		executeFn(0, 0)
		return nil
	}

	// set all linear values to the max possible value
	// learn more about linears and benchmark exxecution in substrate:
	// https://paritytech.github.io/polkadot-sdk/master/frame_benchmarking/v2/index.html
	// https://docs.substrate.io/test/benchmark/
	for i, linear := range *components {
		(*components)[i].setValue(linear.max)
	}

	// iterate steps for each linear component
	for currentComponentIndex, linear := range *components {
		values, err := linear.values(steps)
		if err != nil {
			return err
		}

		for currentStep, v := range values {
			linear.setValue(v)
			executeFn(currentStep+1, currentComponentIndex)
		}
	}

	return nil
}

// Internal function that returns current values for linear components
func componentValues(linearComponents []*linear) []uint32 {
	values := make([]uint32, len(linearComponents))
	for i, l := range linearComponents {
		values[i] = l.Value()
	}

	return values
}
