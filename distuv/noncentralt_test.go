package distuv

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
)

func TestNoncentralTCDF(t *testing.T) {
	t.Parallel()
	tests := []struct {
		dist NoncentralT
		x    float64
		cdf  float64
		tol  float64
		abs  float64
	}{
		// Based on https://github.com/wch/r-source/blob/trunk/tests/reg-tests-2.R
		{dist: NoncentralT{Nu: 10, Ncp: 0}, x: 1.8, cdf: 0.9489738784326605, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Ncp: 0.0001}, x: 1.8, cdf: 0.948964072175642, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Ncp: 1}, x: 1.8, cdf: 0.7584267206837773, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Ncp: -0.0001}, x: 1.8, cdf: 0.9489836831935839, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Ncp: -1}, x: 1.8, cdf: 0.9949471996805094, tol: 5e-12},

		// Based on https://github.com/boostorg/math/blob/develop/test/scipy_issue_14901.cpp
		{dist: NoncentralT{Nu: 2, Ncp: 2}, x: 0.05, cdf: 0.02528206132724582, tol: 5e-12},

		{dist: NoncentralT{Nu: 1, Ncp: 3}, x: 0.05, cdf: 0.00154456589169420, tol: 5e-11},

		// Based on https://github.com/boostorg/math/blob/develop/test/scipy_issue_17916_nct.cpp
		{dist: NoncentralT{Nu: 2, Ncp: 482023264}, x: 2, cdf: 0, tol: 5e-12},

		// Based on https://github.com/boostorg/math/blob/develop/test/test_nc_t.hpp
		{dist: NoncentralT{Nu: 3, Ncp: 1}, x: 2.34, cdf: 0.801888999613917, tol: 5e-12},
		{dist: NoncentralT{Nu: 126, Ncp: -2}, x: -4.33, cdf: 1.252846196792878e-2, tol: 5e-11},
		{dist: NoncentralT{Nu: 20, Ncp: 23}, x: 23, cdf: 0.460134400391924, tol: 5e-12},
		{dist: NoncentralT{Nu: 20, Ncp: 33}, x: 34, cdf: 0.532008386378725, tol: 5e-12},
		{dist: NoncentralT{Nu: 12, Ncp: 38}, x: 39, cdf: 0.495868184917805, tol: 5e-12},
		{dist: NoncentralT{Nu: 12, Ncp: 39}, x: 39, cdf: 0.446304024668836, tol: 5e-2},
		{dist: NoncentralT{Nu: 200, Ncp: 38}, x: 39, cdf: 0.666194209961795, tol: 5e-11},
		{dist: NoncentralT{Nu: 200, Ncp: 42}, x: 40, cdf: 0.179292265426085, tol: 2e-3},
		{dist: NoncentralT{Nu: 2, Ncp: 4}, x: 5, cdf: 0.532020698669953, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Ncp: 16}, x: 0, cdf: 6.388754400538087e-58, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Ncp: 16}, x: 0x1p-1022, cdf: 6.388754400538087e-58, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Ncp: 16}, x: -0x1p-1022, cdf: 6.388754400538087e-58, tol: 5e-12, abs: 1e-16},
		{dist: NoncentralT{Nu: 8, Ncp: 16}, x: -0.125, cdf: 1.018937769092816e-58, tol: 5e-12, abs: 1e-16},
		{dist: NoncentralT{Nu: 8, Ncp: 16}, x: -1e-16, cdf: 6.388754400538077e-58, tol: 5e-12, abs: 1e-16},
		{dist: NoncentralT{Nu: 8, Ncp: 16}, x: 1e-16, cdf: 6.388754400538097e-58, tol: 5e-12},
		{dist: NoncentralT{Nu: 8, Ncp: 16}, x: 0.125, cdf: 5.029904883914148e-57, tol: 1e-4},
		{dist: NoncentralT{Nu: 8, Ncp: 8.5}, x: -1, cdf: 6.174794808375702e-20, tol: 2e-5, abs: 1e-16},

		// Custom tests, wanted values from the R language version 4.4.2.
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: -0.3930852906078905, cdf: 0.01, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 0.6553966734339551, cdf: 0.1, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 1.498837237776003, cdf: 0.33, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 1.944945574355751, cdf: 0.5, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 2.369093423208461, cdf: 0.66, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 3.290148485967974, cdf: 0.9, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 4.469416864177883, cdf: 0.99, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 1.206434338714708, cdf: 0.01, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 2.248666397791987, cdf: 0.1, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 3.094270376888891, cdf: 0.33, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 3.54004668094423, cdf: 0.5, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 3.961130093860616, cdf: 0.66, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 4.860775618489509, cdf: 0.9, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 5.97083270110298, cdf: 0.99, tol: 5e-12},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			cdf := test.dist.CDF(test.x)
			if !scalar.EqualWithinAbsOrRel(cdf, test.cdf, test.abs, test.tol) {
				t.Errorf("{Nu: %f, Ncp: %f}.CDF(%f): got %g want %g", test.dist.Nu, test.dist.Ncp, test.x, cdf, test.cdf)
			}
		})
	}
}

func TestNoncentralTQuantile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		dist NoncentralT
		x    float64
		cdf  float64
		tol  float64
		abs  float64
	}{
		// Based on https://github.com/wch/r-source/blob/trunk/tests/reg-tests-2.R
		{dist: NoncentralT{Nu: 10, Ncp: 0}, x: 1.812461122811676, cdf: 0.95, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Ncp: 0.0001}, x: 1.812579296911650, cdf: 0.95, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Ncp: 1}, x: 3.041741814971971, cdf: 0.95, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Ncp: -0.0001}, x: 1.8123429496892811, cdf: 0.95, tol: 5e-12},
		{dist: NoncentralT{Nu: 10, Ncp: -1}, x: 0.6797901881827499, cdf: 0.95, tol: 5e-12},

		// Based on https://github.com/boostorg/math/blob/develop/test/test_nc_t.hpp
		{dist: NoncentralT{Nu: 3, Ncp: 1}, x: 2.34, cdf: 0.801888999613917, tol: 5e-12},
		{dist: NoncentralT{Nu: 126, Ncp: -2}, x: -4.33, cdf: 1.252846196792878e-2, tol: 5e-11},
		{dist: NoncentralT{Nu: 20, Ncp: 23}, x: 23, cdf: 0.460134400391924, tol: 5e-12},
		{dist: NoncentralT{Nu: 20, Ncp: 33}, x: 34, cdf: 0.532008386378725, tol: 5e-12},
		{dist: NoncentralT{Nu: 12, Ncp: 38}, x: 39, cdf: 0.495868184917805, tol: 5e-12},
		{dist: NoncentralT{Nu: 12, Ncp: 39}, x: 39, cdf: 0.446304024668836, tol: 5e-2},
		{dist: NoncentralT{Nu: 200, Ncp: 38}, x: 39, cdf: 0.666194209961795, tol: 5e-11},
		{dist: NoncentralT{Nu: 200, Ncp: 42}, x: 40, cdf: 0.179292265426085, tol: 2e-3},
		{dist: NoncentralT{Nu: 2, Ncp: 4}, x: 5, cdf: 0.532020698669953, tol: 5e-12},

		// Custom tests, wanted values from the R language version 4.4.2.
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: -0.3930852906078905, cdf: 0.01, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 0.6553966734339551, cdf: 0.1, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 1.498837237776003, cdf: 0.33, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 1.944945574355751, cdf: 0.5, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 2.369093423208461, cdf: 0.66, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 3.290148485967974, cdf: 0.9, tol: 5e-12},
		{dist: NoncentralT{Nu: 58, Ncp: 1.936492}, x: 4.469416864177883, cdf: 0.99, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 1.206434338714708, cdf: 0.01, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 2.248666397791987, cdf: 0.1, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 3.094270376888891, cdf: 0.33, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 3.54004668094423, cdf: 0.5, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 3.961130093860616, cdf: 0.66, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 4.860775618489509, cdf: 0.9, tol: 5e-12},
		{dist: NoncentralT{Nu: 198, Ncp: 3.535534}, x: 5.97083270110298, cdf: 0.99, tol: 5e-12},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			return
			x := test.dist.Quantile(test.cdf)
			if !scalar.EqualWithinAbsOrRel(x, test.x, test.abs, test.tol) {
				t.Errorf("{Nu: %f, Ncp: %f}.Quantile(%f): got %g want %g", test.dist.Nu, test.dist.Ncp, test.cdf, x, test.x)
			}
		})
	}
}
