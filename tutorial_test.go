package evalue

import (
	"bytes"
	"encoding/csv"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"testing"

	"gonum.org/v1/gonum/stat/distuv"
)

// TestOptionalContinuation tests that e-values support optional continuation.
// The output data can be further visualized by plot.py.
func TestOptionalContinuation(t *testing.T) {
	t.Parallel()

	rsrc := rand.NewChaCha8([32]byte{0xb2, 0x11, 0x8a, 0x08, 0x83, 0x15, 0x07, 0x19, 0x64, 0x7a, 0x64, 0x5f, 0x71, 0x7e, 0x07, 0x01, 0xd9, 0x80, 0x61, 0xed, 0xce, 0xaa, 0x4e, 0xf2, 0x2f, 0x36, 0xb5, 0x18, 0x82, 0x85, 0x1c, 0x24})

	const alpha = 0.05
	const delta = 0
	const numSamples = 1e3
	const numBatches = 5
	const batchSize = 40
	rawData := normData(rsrc, delta, numSamples, numBatches*batchSize)

	// Compute the p-values and e-values at every timestep.
	type sampleStat struct {
		data   [2][]float64
		pValue []float64
		eValue []float64
	}
	getPValue := func(x, y []float64) float64 {
		ts := TStat(x, y, 0)
		dist := distuv.StudentsT{Sigma: 1, Nu: ts.Nu, Src: rsrc}
		pValue := 2 * (1 - dist.CDF(math.Abs(ts.T)))
		return pValue
	}
	getEValue := func(x, y []float64) float64 {
		p := NewMom(0.51765)
		eValue := p.EValue(x, y)
		return eValue
	}
	var data []sampleStat
	for _, sampleData := range rawData {
		xFull, yFull := sampleData[0], sampleData[1]
		sample := sampleStat{
			data:   sampleData,
			pValue: make([]float64, len(xFull)),
			eValue: make([]float64, len(xFull)),
		}
		sample.pValue[0] = 1
		sample.eValue[0] = 1
		for i := 1; i < len(xFull); i++ {
			x, y := xFull[:i+1], yFull[:i+1]
			sample.pValue[i] = getPValue(x, y)
			sample.eValue[i] = getEValue(x, y)
		}
		data = append(data, sample)
	}

	// Perform the standard statistical tests at the end of the experiment.
	type statTest struct {
		name  string
		stopT []int
	}
	pValueStd := statTest{name: "p-value", stopT: newStoppingTimes(len(data))}
	for i, sample := range data {
		n := len(sample.pValue) - 1
		if sample.pValue[n] < alpha {
			pValueStd.stopT[i] = n
		}
	}
	eValueStd := statTest{name: "e-value", stopT: newStoppingTimes(len(data))}
	for i, sample := range data {
		n := len(sample.pValue) - 1
		if sample.eValue[n] > 1./alpha {
			eValueStd.stopT[i] = n
		}
	}

	// Perform statistical tests with optional continuation.
	pValueOC := statTest{name: "p-value OC", stopT: newStoppingTimes(len(data))}
	for i, sample := range data {
		for batch := range numBatches {
			if pValueOC.stopT[i] != notStopped {
				continue
			}
			n := (1+batch)*batchSize - 1
			if sample.pValue[n] < alpha {
				pValueOC.stopT[i] = n
			}
		}
	}
	eValueOC := statTest{name: "e-value OC", stopT: newStoppingTimes(len(data))}
	for i, sample := range data {
		for batch := range numBatches {
			if eValueOC.stopT[i] != notStopped {
				continue
			}
			n := (1+batch)*batchSize - 1
			if sample.eValue[n] > 1./alpha {
				eValueOC.stopT[i] = n
			}
		}
	}
	eValueOS := statTest{name: "e-value OS", stopT: newStoppingTimes(len(data))}
	for i, sample := range data {
		for n := range len(sample.eValue) {
			if eValueOS.stopT[i] != notStopped {
				continue
			}
			if sample.eValue[n] > 1./alpha {
				eValueOS.stopT[i] = n
			}
		}
	}

	// Check the type I errors.
	tests := []statTest{pValueStd, pValueOC, eValueStd, eValueOC, eValueOS}
	typeIs := []float64{0.048, 0.147, 0.001, 0.012, 0.041}
	for i, test := range tests {
		var stopped float64
		for _, st := range test.stopT {
			if st != notStopped {
				stopped++
			}
		}
		typeI := stopped / float64(len(test.stopT))
		if typeI != typeIs[i] {
			t.Errorf("%-10s type I error: got %f want %f", test.name, typeI, typeIs[i])
		}
	}

	// Dump data for analysis.
	buf := bytes.NewBuffer(nil)
	w := csv.NewWriter(buf)
	row := []string{"s", "t", "x", "y", "p", "e", "stopTP", "stopTPOC", "stopTE", "stopTEOC"}
	if err := w.Write(row); err != nil {
		t.Fatalf("%+v", err)
	}
	for i, sample := range data {
		x, y := sample.data[0], sample.data[1]
		for j := range len(x) {
			row[0] = strconv.Itoa(i)
			row[1] = strconv.Itoa(j)
			row[2] = strconv.FormatFloat(x[j], 'f', -1, 64)
			row[3] = strconv.FormatFloat(y[j], 'f', -1, 64)
			row[4] = strconv.FormatFloat(sample.pValue[j], 'f', -1, 64)
			row[5] = strconv.FormatFloat(sample.eValue[j], 'f', -1, 64)
			row[6] = strconv.Itoa(pValueStd.stopT[i])
			row[7] = strconv.Itoa(pValueOC.stopT[i])
			row[8] = strconv.Itoa(eValueStd.stopT[i])
			row[9] = strconv.Itoa(eValueOC.stopT[i])
			if err := w.Write(row); err != nil {
				t.Fatalf("%+v", err)
			}
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		t.Fatalf("%+v", err)
	}
	fpath := "/dev/null"
	if err := os.WriteFile(fpath, buf.Bytes(), 0755); err != nil {
		t.Fatalf("%+v", err)
	}
}

func newStoppingTimes(n int) []int {
	stopT := make([]int, n)
	for i := range n {
		stopT[i] = notStopped
	}
	return stopT
}

func normData(rsrc rand.Source, delta float64, numSamples, sampleLen int) [][2][]float64 {
	rnd := rand.New(rsrc)
	var data [][2][]float64
	for range numSamples {
		var sample [2][]float64
		for range sampleLen {
			sample[0] = append(sample[0], rnd.NormFloat64())
			sample[1] = append(sample[1], delta+rnd.NormFloat64())
		}
		data = append(data, sample)
	}
	return data
}
