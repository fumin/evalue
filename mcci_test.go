package evalue

import (
	"math"
	"testing"
)

func TestMCCI(t *testing.T) {
	t.Parallel()
	m0 := m0Gauss{delta: 0}
	const numSamples = 1e6
	const sampleLen = 150
	data := m0.getData(numSamples, sampleLen)

	const alpha = 0.05
	eProcess := NewMom(0.51765)
	t.Logf("g %f", eProcess.G)

	n := 40
	nEff := float64(n * n / (n + n))

	phi0 := 0.65724
	// phi0 = 0.185
	t.Logf("stoppingProb %f", stoppingProb(data, eProcess, alpha, n, phi0))

	var sdObs float64 = 1
	tVal := -math.Sqrt(nEff) / sdObs * 0.65724
	eVal := eProcess.eValue(tVal, n, n)
	t.Logf("eVal %f tVal %f", eVal, tVal)
}
