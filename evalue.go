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

	"gonum.org/v1/exp/root"
	"gonum.org/v1/gonum/mathext"
	"gonum.org/v1/gonum/stat"
)

type Mom struct {
	G float64
}

func NewMom(deltaMin float64) *Mom {
	return &Mom{G: deltaMin * deltaMin / 2}
}

func (p *Mom) EValue(x, y []float64) float64 {
	t := TStat(x, y, 0)
	s := p.eValue(t.T, t.Nu, t.NEff)
	return s
}

// eValue returns the e-value of a t-statistic.
// See equation B4 in Ly for more details.
func (p *Mom) eValue(t, nu, nEff float64) float64 {
	const k = 1
	g := p.G
	e1 := math.Pow(1+nEff*g, -k-1./2)
	e2 := mathext.Hypergeo((nu+1)/2, k+1./2, 1./2, t*t/(nu+t*t)*nEff*g/(1+nEff*g))
	return e1 * e2
}

func (p *Mom) CI(x, y []float64, alpha float64) [2]float64 {
	t := TStat(x, y, 0)
	nu, nEff := t.Nu, t.NEff
	f := func(t float64) float64 { return p.eValue(t, nu, nEff) - 1./alpha }

	// Construct straddle [a, b] to be fed into Brent's method.
	// Since f(0) < 0 always, a=0.
	const a = 0
	// Since the two-sided 95% t-value for the smallest sample size of 1 is 12.706, start the search from around 12.
	var b float64 = 12
	var maxB float64 = b * math.Pow(2, 15)
	tol := math.Nextafter(1, 2) - 1
	// Solve for tAlpha, where f(tAlpha)=0.
	var tAlpha float64
	var err error
	for ; b < maxB; b *= 2 {
		tAlpha, err = root.Brent(f, a, b, tol)
		if err == nil {
			break
		}
	}
	if err != nil {
		return [2]float64{math.Inf(-1), math.Inf(1)}
	}

	width := t.Sp / math.Sqrt(nEff) * tAlpha
	mean := t.Mean1 - t.Mean2
	return [2]float64{mean - width, mean + width}
}

type NPlan struct {
	NPlan int
	Mean  int
}

func getNPlan(alpha, beta, deltaMin float64) NPlan {
	// g is the parameter of the Mom e-process.
	// g := deltaMin * deltaMin / 2
	// ratio is n1/n2.
	const ratio = 1

	return NPlan{}
}

// TStatistic holds information about a t-statistic.
type TStatistic struct {
	// Nu is the degree of freedom.
	Nu float64
	// NEff is the effective sample size.
	NEff float64
	// Mean1 is the mean of the first group.
	Mean1 float64
	// Mean2 is the mean of the second group.
	Mean2 float64
	// Sp is defined in equation 2 in Ly.
	Sp float64
	// t is the t statistic
	T float64
}

// TStat returns the two sample t-statistic.
// See equation 1 in Ly for more details.
func TStat(x1, x2 []float64, phi0 float64) TStatistic {
	n1, n2 := float64(len(x1)), float64(len(x2))
	nu := n1 + n2 - 2
	nEff := n1 * n2 / (n1 + n2)
	mean1 := stat.Mean(x1, nil)
	mean2 := stat.Mean(x2, nil)

	sp := math.Sqrt(1. / nu * ((n1-1)*stat.Variance(x1, nil) + (n2-1)*stat.Variance(x2, nil)))
	t := math.Sqrt(nEff) * (mean1 - mean2 - phi0) / sp

	ts := TStatistic{
		Nu:    nu,
		NEff:  nEff,
		Mean1: mean1,
		Mean2: mean2,
		Sp:    sp,
		T:     t,
	}
	return ts
}
