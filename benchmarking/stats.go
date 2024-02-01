package benchmarking

import "github.com/montanaflynn/stats"

type StatsResult struct {
	Sum float64
	Min float64
	Max float64

	Mean   float64
	Median float64
	Stddev float64

	P99 float64
	P95 float64
	P75 float64
}

func NewStatsResult(input []float64) (StatsResult, error) {
	sum, err := stats.Sum(input)
	if err != nil {
		return StatsResult{}, err
	}

	min, err := stats.Min(input)
	if err != nil {
		return StatsResult{}, err
	}

	max, err := stats.Max(input)
	if err != nil {
		return StatsResult{}, err
	}

	mean, err := stats.Mean(input)
	if err != nil {
		return StatsResult{}, err
	}

	median, err := stats.Median(input)
	if err != nil {
		return StatsResult{}, err
	}

	stddev, err := stats.StandardDeviation(input)
	if err != nil {
		return StatsResult{}, err
	}

	p99, err := stats.Percentile(input, 99.9)
	if err != nil {
		return StatsResult{}, err
	}

	p95, err := stats.Percentile(input, 95)
	if err != nil {
		return StatsResult{}, err
	}

	p75, err := stats.Percentile(input, 75)
	if err != nil {
		return StatsResult{}, err
	}

	return StatsResult{
		Sum:    sum,
		Min:    min,
		Max:    max,
		Mean:   mean,
		Median: median,
		Stddev: stddev,
		P99:    p99,
		P95:    p95,
		P75:    p75,
	}, nil
}
