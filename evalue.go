// Package evalue provides tools for performing e-value statistical tests.
//
// References:
//   - A. Ly, U. Boehm, G., A. Ramdas, D. van Ravenzwaaij. Safe Anytime-Valid Inference: Practical Maximally Flexible Sampling Designs for Experiments Based on e-Values, doi.org/10.31234/osf.io/h5vae
package evalue

// Future work:
// * Implement the eGauss process using the following references:
//   * Last equation of Chapter 1, The Bayesian two-sample t-test, Mithat Gonen, Wesley O. Johnson, Yonggang Lu, Peter H. Westfall
//   * Equation 41, Anytime-valid t-tests and confidence sequences for Gaussian means with unknown variance, Hongjian Wang, Aaditya Ramdas
// * Implement one-sided tests for the mom process using Theorem A.2,
//   Informed Bayesian T-Tests: Online Appendix, Quentin F. Gronau, Alexander Ly, EJ Wagenmakers

import (
	"math"
	"math/rand/v2"
	"slices"

	"gonum.org/v1/exp/root"
	"gonum.org/v1/gonum/mathext"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"

	edistuv "github.com/fumin/evalue/distuv"
)

// notStopped represents the non-existent stopping time of an experiment
// that has not been stopped by a statistical test.
const notStopped = -1

// A Mom is an e-process based on a non-local moment prior.
type Mom struct {
	// G is the tuning parameter of the mom e-process.
	G float64
}

// NewMom creates a mom e-process.
// deltaMin is a lower bound of the true effect size based on domain knowledge.
// The returned mom e-process is tuned such that it rejects the null hypothesis at the fastest rate, when the true data generating process has effect size deltaMin.
func NewMom(deltaMin float64) *Mom {
	return &Mom{G: deltaMin * deltaMin / 2}
}

// EValue returns the e-value of the two sample data.
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

// CI returns the confidence interval of the two sample data.
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

// GetNPlanOptions are options for GetNPlan.
type GetNPlanOptions struct {
	ratio      float64
	numSamples int
	rsrc       rand.Source
}

// NewGetNPlanOptions returns the default GetNPlan options.
func NewGetNPlanOptions() GetNPlanOptions {
	return GetNPlanOptions{
		ratio:      1,
		numSamples: 1000,
		rsrc:       rand.NewChaCha8([32]byte{0x01, 0x08, 0x02, 0x08, 0x83, 0x15, 0x07, 0x19, 0x64, 0x7a, 0x64, 0x5f, 0x71, 0x7e, 0x07, 0x01, 0xd9, 0x80, 0x61, 0xed, 0xce, 0xaa, 0x4e, 0xf2, 0x2f, 0x36, 0xb5, 0x18, 0x82, 0x85, 0x07, 0x01}),
	}
}

// Ratio sets the ratio n1/n2.
func (opt GetNPlanOptions) Ratio(ratio float64) GetNPlanOptions {
	opt.ratio = ratio
	return opt
}

// NumSamples sets the number of samples in simulation.
func (opt GetNPlanOptions) NumSamples(n int) GetNPlanOptions {
	opt.numSamples = n
	return opt
}

// RandSource sets the random source.
func (opt GetNPlanOptions) RandSource(rsrc rand.Source) GetNPlanOptions {
	opt.rsrc = rsrc
	return opt
}

// NPlan is the planned sample size of an experiment.
type NPlan struct {
	// N is the planned sample size with early stopping.
	N int
	// Mean is the average sample size for rejecting the null hypothesis with early stopping.
	Mean int
	// Batch is the sample size without early stopping.
	Batch int

	// EValue is the e-values during simulation.
	EValue [][]float64
	// StopT is the stopping times during simulation.
	StopT []int
}

// GetNPlan returns the planned sample size of an experiment.
// alpha is the significance level, and beta is one minus statistical power.
// deltaMin is a lower bound of the true effect size based on domain knowledge.
func GetNPlan(alpha, beta, deltaMin float64, options ...GetNPlanOptions) NPlan {
	opt := NewGetNPlanOptions()
	if len(options) > 0 {
		opt = options[0]
	}

	// Bound the length of a simulation by the sample size in batch mode.
	// Experiments with early stopping always need smaller sample sizes than those in batch mode which are done without early stopping.
	p := NewMom(deltaMin)
	nPlanBatch1, nPlanBatch2 := getNPlanBatch(alpha, beta, deltaMin, opt.ratio, p)
	nPlan := NPlan{Batch: nPlanBatch1}

	// Interpolate n1 and n2.
	var n1Vector, n2Vector []int
	for i := 1; i <= nPlanBatch1; i++ {
		n1Vector = append(n1Vector, i)
		n2Vector = append(n2Vector, int(math.Ceil(opt.ratio*float64(i))))
	}

	// Simulation experiments.
	rnd := rand.New(opt.rsrc)
	sampleLen := max(nPlanBatch1, nPlanBatch2)
	sample1, sample2 := make([]float64, sampleLen), make([]float64, sampleLen)
	interpolate1 := newInterpolator(len(n1Vector), len(sample1))
	interpolate2 := newInterpolator(len(n2Vector), len(sample2))
	for range opt.numSamples {
		// Generate simulation data.
		for i := range sampleLen {
			sample1[i] = deltaMin/2 + rnd.NormFloat64()
			sample2[i] = -deltaMin/2 + rnd.NormFloat64()
		}

		// Interpolate between n1 and n2, so that the resulting slices are of the same length.
		x1Bar, x1Square := interpolate1.do(n1Vector, sample1)
		x2Bar, x2Square := interpolate2.do(n2Vector, sample2)

		// Simulate an experiment with early stopping.
		var eValues []float64
		stopT := notStopped
		for i := range n1Vector {
			n1, n2 := float64(n1Vector[i]), float64(n2Vector[i])
			nu, nEff := n1+n2-2, n1*n2/(n1+n2)
			x1, x2 := x1Bar[i], x2Bar[i]
			x1Sq, x2Sq := x1Square[i], x2Square[i]

			// Compute e-value.
			var eVal float64 = 1
			if nu > 0 {
				sp := math.Sqrt(1. / nu * (x1Sq - n1*x1*x1 + x2Sq - n2*x2*x2))
				t := math.Sqrt(nEff) * (x1 - x2) / sp
				eVal = p.eValue(t, nu, nEff)
			}
			eValues = append(eValues, eVal)

			// Perform test with optional stopping.
			if eVal > 1./alpha {
				stopT = int(n1)
				break
			}
		}

		nPlan.EValue = append(nPlan.EValue, eValues)
		nPlan.StopT = append(nPlan.StopT, stopT)
	}

	// Compute sample size for the desired statistical power.
	stopT := make([]float64, len(nPlan.StopT))
	for i, t := range nPlan.StopT {
		if t == notStopped {
			stopT[i] = math.Inf(1)
		} else {
			stopT[i] = float64(t)
		}
	}
	slices.Sort(stopT)
	nPlan.N = int(math.Ceil(stat.Quantile(1-beta, stat.LinInterp, stopT, nil)))

	// Calculate the average stopping time, assuming we go according to plan.
	for i := range stopT {
		stopT[i] = min(float64(nPlan.N), stopT[i])
	}
	nPlan.Mean = int(math.Ceil(stat.Mean(stopT, nil)))

	return nPlan
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

func getNPlanBatch(alpha, beta, delta, ratio float64, p *Mom) (int, int) {
	// Define the function f that returns eValue - 1/alpha, given nEff.
	delta = math.Abs(delta)
	f := func(nEff float64) float64 {
		nu := math.Pow(1+ratio, 2)/ratio*nEff - 2
		t := edistuv.NoncentralT{Nu: nu, Ncp: math.Sqrt(nEff) * delta}.Quantile(beta)
		s := p.eValue(t, nu, nEff)
		return s - 1./alpha
	}

	// Solve for the root of f.
	//
	// Find the bracket that wraps the root.
	qB := distuv.Normal{Sigma: 1}.Quantile(beta)
	guess := 2 / (delta * delta) * (qB*qB - qB*math.Sqrt(qB*qB+2*math.Log(1./alpha)) + math.Log(1./alpha))
	a, b := edistuv.FindBracketMono(f, guess)
	// Find the root inside the bracket.
	eps := math.Nextafter(1, 2) - 1
	tol := math.Pow(eps, 0.25)
	nEff, err := root.Brent(f, a, b, tol)
	if err != nil {
		return -1, -1
	}

	n1 := int(math.Ceil(nEff * (1 + ratio) / ratio))
	n2 := int(math.Ceil(nEff * (1 + ratio)))
	return n1, n2
}

// interpolator holds buffers for interpolating between two groups of data of different sizes.
type interpolator struct {
	x    []float64
	x2   []float64
	cum  []float64
	cum2 []float64
}

func newInterpolator(n, sampleLen int) interpolator {
	var b interpolator
	b.x = make([]float64, n)
	b.x2 = make([]float64, n)
	b.cum = make([]float64, sampleLen)
	b.cum2 = make([]float64, sampleLen)
	return b
}

func (b interpolator) do(ns []int, sample []float64) ([]float64, []float64) {
	// Compute the cumulative sums.
	s0 := sample[0]
	b.cum[0], b.cum2[0] = s0, s0*s0
	for i := 1; i < len(b.cum); i++ {
		x := sample[i]
		b.cum[i] = b.cum[i-1] + x
		b.cum2[i] = b.cum2[i-1] + x*x
	}

	// Interpolate according to ns.
	for i, n := range ns {
		b.x[i] = b.cum[n-1] / float64(n)
		b.x2[i] = b.cum2[n-1]
	}
	return b.x, b.x2
}
