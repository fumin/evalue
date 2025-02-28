package distuv

import (
	"math"

	"gonum.org/v1/gonum/mathext"
)

type NoncentralT struct {
	Nu float64

	Ncp float64
}

// Russell Lenth, Algorithm AS 243: Cumulative Distribution Function of the Non-Central T Distribution
func (dist NoncentralT) CDF(t float64) float64 {
	df, ncp := dist.Nu, dist.Ncp

	var albeta, a, b, del, lambda, rxb, tt, x float64
	var geven, godd, p, q, s, tnc, xeven, xodd float64
	var negdel bool

	const itrmax = 1000
	const errmax = 1e-12

	if df <= 0 {
		return math.NaN()
	}
	if t >= 0 {
		negdel, tt, del = false, t, ncp
	} else {
		negdel, tt, del = true, -t, -ncp
	}

	// Initialize twin series.
	// Guenther, J. (1978). Statist. Computn. Simuln. vol.6, 199.
	x = t * t
	rxb = df / (x + df)
	x = x / (x + df)
	if x > 0 { // t != 0
		lambda = del * del
		p = 0.5 * math.Exp(-0.5*lambda)
		if p == 0 {
			goto overflow
		}
		// sqrt2dPi is sqrt(2/pi).
		const sqrt2dPi = 0.797884560802865355879892119869
		q = sqrt2dPi * p * del
		s = 0.5 - p
		if s < 1e-7 {
			// s = 0.5 - p = 0.5*(1 - exp(-.5 L)) =  -0.5*expm1(-.5 L))
			s = -0.5 * math.Expm1(-0.5*lambda)
		}
		a = 0.5
		b = 0.5 * df
		rxb = math.Pow(rxb, b) // equivalent to pow(1-x, b)
		// lnSqrtPi is log(sqrt(pi))
		const lnSqrtPi = 0.572364942924700087071713675677
		albeta = lnSqrtPi + lgamma(b) - lgamma(0.5+b)
		xodd = pbeta(x, a, b)
		godd = 2 * rxb * math.Exp(a*math.Log(x)-albeta)
		tnc = b * x
		xeven = 1 - rxb
		geven = tnc * rxb
		tnc = p*xodd + q*xeven

		// repeat until convergence or iteration limit
		for it := 1; it <= itrmax; it++ {
			a += 1
			xodd -= godd
			xeven -= geven
			godd *= x * (a + b - 1) / a
			geven *= x * (a + b - 0.5) / (a + 0.5)
			p *= lambda / (2 * float64(it))
			q *= lambda / (2*float64(it) + 1)
			s -= p
			tnc += p*xodd + q*xeven

			// R 2.4.0 added test for rounding error here.
			if s < -1e-10 { // happens e.g. for (t,df,ncp)=(40,10,38.5), after 799 it.
				goto finis
			}
			if s <= 0 && it > 1 {
				goto finis
			}

			errbd := 2 * s * (xodd - godd)
			if math.Abs(errbd) < errmax {
				goto finis
			}
		}
	} else { // x = t = 0
		tnc = 0
	}

overflow:
	// Approx. from	 Abramowitz & Stegun 26.7.10 (p.949).
	s = 1. / (4 * df)
	return pnorm(tt*(1-s), del, math.Sqrt(1+tt*tt*2*s), !negdel)
finis:
	tnc += pnorm(-del, 0, 1, true)

	p = min(tnc, 1) // Precaution
	if !negdel {
		return p
	} else {
		// Use 0.5 - p + 0.5 to perhaps gain 1 bit of accuracy
		return 0.5 - p + 0.5
	}
}

func (dist NoncentralT) Quantile(p float64) float64 {
	return -1
}

func pbeta(x, a, b float64) float64 {
	return mathext.RegIncBeta(a, b, x)
}

func pnorm(x, mu, sigma float64, lowerTail bool) float64 {
	p := 0.5 * math.Erfc(-(x-mu)/(sigma*math.Sqrt2))
	if lowerTail {
		return p
	}
	return 1 - p
}

func lgamma(x float64) float64 {
	y, _ := math.Lgamma(x)
	return y
}
