// Package evalue provides tools for perform e-value statistical tests.
//
// References:
//   - A. Ly, U. Boehm, G., A. Ramdas, D. van Ravenzwaaij. Safe Anytime-Valid Inference: Practical Maximally Flexible Sampling Designs for Experiments Based on e-Values, doi.org/10.31234/osf.io/h5vae
package evalue

// eGauss:
// last equation of Chapter 1, The Bayesian two-sample t-test, Mithat Gonen, Wesley O. Johnson, Yonggang Lu, Peter H. Westfall
// Equation 41, Anytime-valid t-tests and confidence sequences for Gaussian means with unknown variance, Hongjian Wang, Aaditya Ramdas

import (
	"math"

	"gonum.org/v1/gonum/mathext"
	"gonum.org/v1/gonum/stat"
)

type Mom struct {
	G float64
}

func NewMom(deltaMin float64) *Mom {
	return &Mom{G: deltaMin * deltaMin / 2}
}

func (p *Mom) EValue(x, y []float64, phi0 float64) float64 {
	t := tStatistic(x, y, phi0)
	s := p.eValue(t, len(x), len(y))
	// log.Printf("eVal %f tVal %f", s, t)
	return s
}

// eValue returns the e-value of a t-statistic.
// See equation B4 in Ly for more details.
func (p *Mom) eValue(t float64, n1i, n2i int) float64 {
	const k float64 = 1
	n1, n2 := float64(n1i), float64(n2i)
	nu := n1 + n2 - 2
	nEff := n1 * n2 / (n1 + n2)
	g := p.G
	e1 := math.Pow(1+nEff*g, -k-1./2)
	e2 := mathext.Hypergeo((nu+1)/2, k+1./2, 1./2, t*t/(nu+t*t)*nEff*g/(1+nEff*g))
	return e1 * e2
}

// tStatistic returns the two sample t-statistic.
// See equation 1 in Ly for more details.
func tStatistic(x1, x2 []float64, phi0 float64) float64 {
	n1, n2 := float64(len(x1)), float64(len(x2))
	nu := n1 + n2 - 2
	nEff := n1 * n2 / (n1 + n2)
	x1Mean := stat.Mean(x1, nil)
	x2Mean := stat.Mean(x2, nil)

	sp := math.Sqrt(1. / nu * ((n1-1)*stat.Variance(x1, nil) + (n2-1)*stat.Variance(x2, nil)))
	tStat := math.Sqrt(nEff) * (x1Mean - x2Mean - phi0) / sp
	// log.Printf("tStat %f nEff %f x1-x2 %f sp %f", tStat, nEff, x1Mean-x2Mean, sp)
	return tStat
}
