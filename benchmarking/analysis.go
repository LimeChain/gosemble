package benchmarking

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	benchmarkingtypes "github.com/LimeChain/gosemble/primitives/benchmarking"
)

var (
	errZeroBenchmarkResults = errors.New("provided benchmark results must be more than 0.")
)

// todo other analysis types
// todo enum for analysis types
// todo all enum values to implement Analysis interface
// todo implement benchmark flag for analysis choice
// type analysisType int

// const (
// 	MinSquares analysisType = iota
// 	MedianSlopes
// )

// func (a analysisType) Analysis(benchmarkResults []benchmarkResult) (extrinsicTime, reads, writes analysis, err error) {
// 	if len(benchmarkResults) == 0 {
// 		err = errZeroBenchmarkResults
// 		return
// 	}

// 	if len(benchmarkResults[0].components) == 1 {
// 		return medianValuesAnalysis(benchmarkResults)
// 	}

// 	switch a {
// 	case MinSquares:
// 		return minSquaresAnalysis(benchmarkResults)
// 	case MedianSlopes:
// 		return medianSlopesAnalysis(benchmarkResults)
// 	default:
// 		return
// 	}
// }

type benchmarkResult struct {
	components    []uint32
	extrinsicTime uint64
	reads, writes uint64
}

func newBenchmarkResult(benchmarkRes benchmarkingtypes.BenchmarkResult, componentValues []uint32) benchmarkResult {
	return benchmarkResult{
		extrinsicTime: benchmarkRes.ExtrinsicTime.ToBigInt().Uint64(),
		reads:         uint64(benchmarkRes.Reads),
		writes:        uint64(benchmarkRes.Writes),
		components:    componentValues,
	}
}

type analysis struct {
	base   uint64
	slopes []uint64
	names  []string
	// todo value_dists: Option<Vec<(Vec<u32>, u128, u128)>>
	minimum uint64
}

func medianSlopesAnalysis(benchmarkResults []benchmarkResult) (extrinsicTime, reads, writes analysis, err error) {
	if len(benchmarkResults) == 0 {
		err = errZeroBenchmarkResults
		return
	}

	if len(benchmarkResults[0].components) == 1 {
		return medianValuesAnalysis(benchmarkResults)
	}

	results := make([]struct {
		others []float64
		values []struct {
			componentValue               float64
			extrinsicTime, reads, writes float64
		}
	}, len(benchmarkResults[0].components))

	for i, _ := range benchmarkResults[0].components {
		// count each component value combination
		counted := map[string]int{}
		for _, br := range benchmarkResults {
			componentValues := make([]uint32, len(br.components))
			copy(componentValues, br.components)
			componentValues[i] = 0
			counted[fmt.Sprintf("%d", componentValues)]++
		}

		// get the component values with highest count (as string)
		highestCountKey := ""
		highestCount := 0
		for key, count := range counted {
			if count > highestCount {
				highestCount = count
				highestCountKey = key
			}
		}

		// convert component values from string to []uint64
		others := make([]float64, len(benchmarkResults[0].components))
		for y, v := range strings.Split(highestCountKey[1:len(highestCountKey)-1], " ") {
			num, errParse := strconv.ParseUint(v, 10, 64)
			if errParse != nil {
				err = errParse
				return
			}
			others[y] = float64(num)
		}

		results[i].others = others

		for _, br := range benchmarkResults {
			isValid := true
			for y, v := range br.components {
				if y != i && float64(v) != others[y] {
					isValid = false
					continue
				}
			}

			if !isValid {
				continue
			}

			results[i].values = append(
				results[i].values, struct {
					componentValue               float64
					extrinsicTime, reads, writes float64
				}{float64(br.components[i]), float64(br.extrinsicTime), float64(br.reads), float64(br.writes)},
			)
		}
	}

	models := make([]struct {
		offsetExtrinsicTime, offsetReads, offsetWrites float64
		slopeExtrinsicTime, slopeReads, slopeWrites    float64
	}, len(results))

	for i, r := range results {
		slopes := []struct{ slopeExtrinsicTime, slopeReads, slopeWrites float64 }{}
		for y, v1 := range r.values {
			for _, v2 := range r.values[y+1:] {
				if v1.componentValue != v2.componentValue {
					slopes = append(slopes, struct{ slopeExtrinsicTime, slopeReads, slopeWrites float64 }{
						(v1.extrinsicTime - v2.extrinsicTime) / (v1.componentValue - v2.componentValue),
						(v1.reads - v2.reads) / (v1.componentValue - v2.componentValue),
						(v1.writes - v2.writes) / (v1.componentValue - v2.componentValue),
					})
				}
			}
		}

		midIndex := len(slopes) / 2

		// slope extrinsic time
		sort.Slice(slopes, func(i, j int) bool {
			return uint64(slopes[i].slopeExtrinsicTime) < uint64(slopes[j].slopeExtrinsicTime)
		})
		models[i].slopeExtrinsicTime = slopes[midIndex].slopeExtrinsicTime

		// slope reads
		sort.Slice(slopes, func(i, j int) bool {
			return uint64(slopes[i].slopeReads) < uint64(slopes[j].slopeReads)
		})
		models[i].slopeReads = slopes[midIndex].slopeReads

		// slope writes
		sort.Slice(slopes, func(i, j int) bool {
			return uint64(slopes[i].slopeWrites) < uint64(slopes[j].slopeWrites)
		})
		models[i].slopeWrites = slopes[midIndex].slopeWrites

		offsets := []struct{ offsetExtrinsicTime, offsetReads, offsetWrites float64 }{}
		for _, v := range r.values {
			offsets = append(offsets, struct{ offsetExtrinsicTime, offsetReads, offsetWrites float64 }{
				float64(v.extrinsicTime) - models[i].slopeExtrinsicTime*float64(v.componentValue),
				float64(v.reads) - models[i].slopeReads*float64(v.componentValue),
				float64(v.writes) - models[i].slopeWrites*float64(v.componentValue),
			})
		}

		midIndex = len(offsets) / 2

		// offset extrinsic time
		sort.Slice(offsets, func(i, j int) bool {
			return uint64(offsets[i].offsetExtrinsicTime) < uint64(offsets[j].offsetExtrinsicTime)
		})
		models[i].offsetExtrinsicTime = offsets[midIndex].offsetExtrinsicTime

		// offset reads
		sort.Slice(offsets, func(i, j int) bool {
			return uint64(offsets[i].offsetReads) < uint64(offsets[j].offsetReads)
		})
		models[i].offsetReads = offsets[midIndex].offsetReads

		// offset writes
		sort.Slice(offsets, func(i, j int) bool {
			return uint64(offsets[i].offsetWrites) < uint64(offsets[j].offsetWrites)
		})
		models[i].offsetWrites = offsets[midIndex].offsetWrites
	}

	for i, _ := range models {
		over := struct{ overExtrinsicTime, overReads, overWrites float64 }{}

		for y, o := range results[i].others {
			if y != i {
				over.overExtrinsicTime += models[y].slopeExtrinsicTime * o
				over.overReads += models[y].slopeReads * o
				over.overWrites += models[y].slopeWrites * o
			}
		}

		models[i].offsetExtrinsicTime -= over.overExtrinsicTime
		models[i].offsetReads -= over.overReads
		models[i].offsetWrites -= over.overWrites
	}

	// extrinsic time
	extrinsicTime.base = uint64((math.Max(models[0].offsetExtrinsicTime, 0) + 0.000_000_005) * 1000)
	for _, m := range models {
		extrinsicTime.slopes = append(extrinsicTime.slopes, uint64((math.Max(m.slopeExtrinsicTime, 0)+0.000_000_005)*1000))
	}

	sort.Slice(benchmarkResults, func(i, j int) bool {
		return benchmarkResults[i].extrinsicTime < benchmarkResults[j].extrinsicTime
	})
	extrinsicTime.minimum = benchmarkResults[0].extrinsicTime

	// reads
	reads.base = uint64(math.Max(models[0].offsetReads, 0) + 0.000_000_005)
	for _, m := range models {
		reads.slopes = append(reads.slopes, uint64(math.Max(m.slopeReads, 0)+0.000_000_005))
	}

	sort.Slice(benchmarkResults, func(i, j int) bool {
		return benchmarkResults[i].reads < benchmarkResults[j].reads
	})
	reads.minimum = benchmarkResults[0].reads

	//  writes
	writes.base = uint64(math.Max(models[0].offsetWrites, 0) + 0.000_000_005)
	for _, m := range models {
		writes.slopes = append(writes.slopes, uint64(math.Max(m.slopeWrites, 0)+0.000_000_005))
	}

	sort.Slice(benchmarkResults, func(i, j int) bool {
		return benchmarkResults[i].writes < benchmarkResults[j].writes
	})
	writes.minimum = benchmarkResults[0].writes

	return
}

func minSquaresAnalysis(benchmarkResults []benchmarkResult) (extrinsicTime, reads, writes analysis, err error) {
	// todo
	return
}

func medianValuesAnalysis(benchmarkResults []benchmarkResult) (extrinsicTime, reads, writes analysis, err error) {
	midIndex := len(benchmarkResults) / 2

	// extrinsic time
	sort.Slice(benchmarkResults, func(i, j int) bool {
		return benchmarkResults[i].extrinsicTime < benchmarkResults[j].extrinsicTime
	})

	extrinsicTime.base = benchmarkResults[midIndex].extrinsicTime * 1000
	extrinsicTime.minimum = benchmarkResults[0].extrinsicTime

	// reads
	sort.Slice(benchmarkResults, func(i, j int) bool {
		return benchmarkResults[i].reads < benchmarkResults[j].reads
	})

	reads.base = benchmarkResults[midIndex].reads
	reads.minimum = benchmarkResults[0].reads

	// writes
	sort.Slice(benchmarkResults, func(i, j int) bool {
		return benchmarkResults[i].writes < benchmarkResults[j].writes
	})

	writes.base = benchmarkResults[midIndex].writes
	writes.minimum = benchmarkResults[0].writes

	return
}
