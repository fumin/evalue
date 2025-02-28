package evalue

import (
	"math/rand"
)

type m0Gauss struct {
	delta float64
}

func (m m0Gauss) getData(numSamples, sampleLen int) [][2][]float64 {
	var data [][2][]float64
	for _ = range numSamples {
		var sample [2][]float64
		for _ = range sampleLen {
			sample[0] = append(sample[0], rand.NormFloat64())
			sample[1] = append(sample[1], m.delta+rand.NormFloat64())
		}
		data = append(data, sample)
	}
	return data
}

func stoppingProb(data [][2][]float64, eProcess *Mom, alpha float64, n int, phi0 float64) float64 {
	var numStopped float64
	for _, sample := range data {
		x, y := sample[0], sample[1]

		stoppingTime := -1
		for i := 2; i <= n; i++ {
			// for i := n; i <= n; i++ {
			xi, yi := x[:i], y[:i]
			s := eProcess.EValue(xi, yi, phi0)
			if s >= 1./alpha {
				stoppingTime = n
				break
			}
		}

		if stoppingTime != -1 {
			numStopped++
		}
	}
	return numStopped / float64(len(data))
}
