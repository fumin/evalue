package evalue

import (
	"bytes"
	"cmp"
	_ "embed"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"slices"
	"strconv"
	"testing"

	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/stat/distuv"
)

// TestSaviTutorial_1_1 tests for the findings in Section 1.1, Savi Tutorial.
func TestSaviTutorial_1_1(t *testing.T) {
	const alpha = 0.05

	if len(grayData) != 61 {
		t.Errorf("wrong data size %d", len(grayData))
	}
	type tTest struct {
		pUsed int
		p     float64
		eUsed int
		e     float64
		stopT int
	}
	tTests := []tTest{}
	for _, study := range grayData {
		var tt tTest
		// p-value based test.
		tt.pUsed = len(study)
		x, y := splitGray(study[:tt.pUsed])
		tStat := TStat(x, y, 0)
		studentT := distuv.StudentsT{Sigma: 1, Nu: tStat.Nu}
		tt.p = 2 * (1 - studentT.CDF(math.Abs(tStat.T)))

		// e-value based test.
		tt.stopT = notStopped
		nPlan := len(study)
		for n := 1; n <= nPlan; n++ {
			tt.eUsed = n
			dataN := study[:tt.eUsed]
			x, y := splitGray(dataN)
			if !(len(x) > 1 && len(y) > 1) {
				continue
			}

			p := NewMom(0.769)
			tt.e = p.EValue(x, y)
			if tt.e > 1./alpha {
				tt.stopT = n
				break
			}
		}
		tTests = append(tTests, tt)
	}

	// Check p-value based results
	pUsed := 0
	pRejects := 0
	for _, tt := range tTests {
		pUsed += tt.pUsed

		if tt.p < alpha {
			pRejects++
		}
	}
	if !(pUsed == 8002 && pRejects == 58) {
		t.Errorf("wrong p-value results %d %d", pUsed, pRejects)
	}

	eUsed := 0
	eRejects := 0
	for _, tt := range tTests {
		eUsed += tt.eUsed
		if tt.stopT != notStopped {
			eRejects++
		}
	}
	if !(eUsed == 2655 && eRejects == 54) {
		t.Errorf("wrong e-value results %d %d", eUsed, eRejects)
	}
}

func TestEValue(t *testing.T) {
	t.Parallel()
	data := grayData[slices.IndexFunc(grayData, func(d []grayCase) bool { return d[0].location == "Carleton University, Ottawa, Canada" })]
	tests := []struct {
		n int
		s float64
	}{
		{n: 1, s: 1},
		{n: 2, s: 1},
		{n: 3, s: 1},
		{n: 4, s: 1},
		{n: 5, s: 1},
		{n: 6, s: 1},
		{n: 7, s: 1},
		{n: 8, s: 1},
		{n: 9, s: 1.120292},
		{n: 10, s: 1.226478},
		{n: 11, s: 1.342858},
		{n: 12, s: 1.17818},
		{n: 13, s: 1.270546},
		{n: 14, s: 1.068791},
		{n: 15, s: 1.177223},
		{n: 16, s: 1.043421},
		{n: 17, s: 1.732295},
		{n: 18, s: 1.976986},
		{n: 19, s: 2.251658},
		{n: 20, s: 2.635829},
		{n: 21, s: 3.072126},
		{n: 22, s: 5.186749},
		{n: 23, s: 7.28843},
		{n: 24, s: 8.306582},
		{n: 25, s: 14.75564},
		{n: 26, s: 12.90174},
		{n: 27, s: 12.45841},
		{n: 28, s: 10.91521},
		{n: 29, s: 13.57739},
		{n: 30, s: 21.42713},
		{n: 31, s: 28.88713},
		{n: 32, s: 25.63158},
		{n: 33, s: 28.66248},
		{n: 34, s: 38.5593},
		{n: 35, s: 68.69083},
		{n: 36, s: 60.62451},
		{n: 37, s: 78.08596},
		{n: 38, s: 36.21235},
		{n: 39, s: 41.99709},
		{n: 40, s: 21.89799},
		{n: 41, s: 27.05205},
		{n: 42, s: 37.23797},
		{n: 43, s: 35.06838},
		{n: 44, s: 47.88725},
		{n: 45, s: 60.86173},
		{n: 46, s: 57.41992},
		{n: 47, s: 38.42893},
		{n: 48, s: 21.22604},
		{n: 49, s: 29.29514},
		{n: 50, s: 17.15684},
		{n: 51, s: 13.50019},
		{n: 52, s: 13.18857},
		{n: 53, s: 17.93653},
		{n: 54, s: 18.36552},
		{n: 55, s: 18.03706},
		{n: 56, s: 17.75962},
		{n: 57, s: 17.52732},
		{n: 58, s: 17.33523},
		{n: 59, s: 23.04155},
		{n: 60, s: 22.757},
		{n: 61, s: 26.53848},
		{n: 62, s: 27.07536},
		{n: 63, s: 27.66282},
		{n: 64, s: 33.92066},
		{n: 65, s: 50.16185},
		{n: 66, s: 71.70682},
		{n: 67, s: 107.0897},
		{n: 68, s: 134.3576},
		{n: 69, s: 139.0017},
		{n: 70, s: 143.9816},
		{n: 71, s: 231.0347},
		{n: 72, s: 144.331},
		{n: 73, s: 218.0786},
		{n: 74, s: 207.575},
		{n: 75, s: 264.4869},
		{n: 76, s: 421.306},
		{n: 77, s: 543.1127},
		{n: 78, s: 801.0143},
		{n: 79, s: 851.519},
		{n: 80, s: 711.7536},
		{n: 81, s: 1036.912},
		{n: 82, s: 540.2938},
		{n: 83, s: 775.9812},
		{n: 84, s: 1112.698},
		{n: 85, s: 1593.059},
		{n: 86, s: 1540.007},
		{n: 87, s: 1933.166},
		{n: 88, s: 2354.524},
		{n: 89, s: 2970.068},
		{n: 90, s: 4326.574},
		{n: 91, s: 6310.25},
		{n: 92, s: 6559.974},
		{n: 93, s: 9524.595},
		{n: 94, s: 13809.68},
		{n: 95, s: 19995.5},
		{n: 96, s: 15766.54},
		{n: 97, s: 15326.93},
		{n: 98, s: 19330.53},
		{n: 99, s: 24372.01},
		{n: 100, s: 42520.48},
		{n: 101, s: 65342.26},
		{n: 102, s: 36457.17},
		{n: 103, s: 46522.3},
		{n: 104, s: 80869.53},
		{n: 105, s: 147884.5},
		{n: 106, s: 156659.6},
		{n: 107, s: 232164.9},
		{n: 108, s: 343616.2},
		{n: 109, s: 432289.6},
		{n: 110, s: 408974.1},
		{n: 111, s: 631288.3},
		{n: 112, s: 418056.6},
		{n: 113, s: 439154.8},
		{n: 114, s: 345933.3},
		{n: 115, s: 364562.3},
		{n: 116, s: 344469.9},
		{n: 117, s: 512893.8},
		{n: 118, s: 487417.9},
		{n: 119, s: 268395.9},
		{n: 120, s: 282045.1},
		{n: 121, s: 266929.8},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			data = data[:test.n]
			x, y := splitGray(data)
			if !(len(x) > 1 && len(y) > 1) {
				return
			}

			p := &Mom{G: 0.1339827}
			s := p.EValue(x, y)
			if !scalar.EqualWithinRel(s, test.s, 2e-6) {
				t.Errorf("unexpected result EValue(data[:%d]): got %f want %f", test.n, s, test.s)
			}
		})
	}
}

func TestEValueT(t *testing.T) {
	t.Parallel()
	tests := []struct {
		t    float64
		n1   int
		n2   int
		want float64
	}{
		{t: 0, n1: 2, n2: 2, want: 0.8281145},
		{t: 0.3429972, n1: 3, n2: 3, want: 0.7874301},
		{t: 0.3429972, n1: 4, n2: 4, want: 0.7307366},
		{t: 1.073751, n1: 5, n2: 5, want: 0.9691936},
		{t: 1.431668, n1: 6, n2: 5, want: 1.230795},
		{t: 1.550761, n1: 7, n2: 6, want: 1.374565},
		{t: 1.94993, n1: 8, n2: 7, want: 2.019262},
		{t: 2.244057, n1: 9, n2: 8, want: 2.8438},
		{t: 2.169469, n1: 10, n2: 9, want: 2.802736},
		{t: 2.251096, n1: 11, n2: 10, want: 3.241674},
		{t: 2.499726, n1: 12, n2: 10, want: 4.47004},
		{t: 2.530404, n1: 13, n2: 11, want: 4.960147},
		{t: 2.94071, n1: 14, n2: 12, want: 9.103413},
		{t: 3.0295, n1: 15, n2: 13, want: 11.1969},
		{t: 3.520671, n1: 16, n2: 14, want: 24.76002},
		{t: 3.332526, n1: 17, n2: 15, want: 20.80801},
		{t: 3.479162, n1: 18, n2: 15, want: 27.2232},
		{t: 3.202694, n1: 19, n2: 16, want: 19.07348},
		{t: 3.292121, n1: 20, n2: 17, want: 23.78914},
		{t: 3.086168, n1: 21, n2: 18, want: 17.69949},
		{t: 3.341982, n1: 22, n2: 19, want: 29.70011},
		{t: 3.398061, n1: 23, n2: 20, want: 35.06838},
		{t: 3.501578, n1: 24, n2: 20, want: 43.9315},
		{t: 3.43804, n1: 25, n2: 21, want: 41.25348},
		{t: 3.542218, n1: 26, n2: 22, want: 53.8048},
		{t: 3.376091, n1: 27, n2: 23, want: 40.28666},
		{t: 3.110812, n1: 28, n2: 24, want: 24.09196},
		{t: 2.991821, n1: 29, n2: 25, want: 19.2767},
		{t: 3.000492, n1: 30, n2: 25, want: 19.83799},
		{t: 3.061378, n1: 31, n2: 26, want: 23.12938},
		{t: 3.17577, n1: 32, n2: 27, want: 30.46774},
		{t: 3.238277, n1: 33, n2: 28, want: 35.95707},
		{t: 3.229884, n1: 34, n2: 29, want: 36.11804},
		{t: 3.222106, n1: 35, n2: 30, want: 36.2442},
		{t: 3.034851, n1: 36, n2: 30, want: 23.6237},
		{t: 3.24621, n1: 37, n2: 31, want: 39.49257},
		{t: 3.305981, n1: 38, n2: 32, want: 46.50529},
		{t: 3.499248, n1: 39, n2: 33, want: 76.64593},
		{t: 3.683814, n1: 40, n2: 34, want: 125.8569},
		{t: 3.780777, n1: 41, n2: 35, want: 166.7825},
		{t: 3.900473, n1: 42, n2: 35, want: 232.4453},
		{t: 4.214193, n1: 43, n2: 36, want: 567.0516},
		{t: 4.497879, n1: 44, n2: 37, want: 1319.02},
		{t: 4.590277, n1: 45, n2: 38, want: 1809.453},
		{t: 4.843779, n1: 46, n2: 39, want: 4034.909},
		{t: 4.667405, n1: 47, n2: 40, want: 2507.086},
		{t: 4.777502, n1: 48, n2: 40, want: 3574.364},
		{t: 4.762545, n1: 49, n2: 41, want: 3589.371},
		{t: 4.931264, n1: 50, n2: 42, want: 6368.81},
		{t: 5.157969, n1: 51, n2: 43, want: 13809.68},
		{t: 5.247082, n1: 52, n2: 44, want: 19532.6},
		{t: 5.318573, n1: 53, n2: 45, want: 26205.62},
		{t: 5.383296, n1: 54, n2: 45, want: 33386.48},
		{t: 5.564101, n1: 55, n2: 46, want: 65342.26},
		{t: 5.533735, n1: 56, n2: 47, want: 62754.71},
		{t: 5.76241, n1: 57, n2: 48, want: 147884.5},
		{t: 5.820789, n1: 58, n2: 49, want: 194347.9},
		{t: 5.897718, n1: 59, n2: 50, want: 273684.7},
		{t: 6.000207, n1: 60, n2: 50, want: 408974.1},
		{t: 6.11447, n1: 61, n2: 51, want: 666879},
		{t: 5.917839, n1: 62, n2: 52, want: 345933.3},
		{t: 5.899927, n1: 63, n2: 53, want: 344469.9},
		{t: 5.976485, n1: 64, n2: 54, want: 487417.9},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			n1, n2 := float64(test.n1), float64(test.n2)
			nu := n1 + n2 - 2
			nEff := n1 * n2 / (n1 + n2)
			p := &Mom{G: 0.1339827}
			s := p.eValue(test.t, nu, nEff)
			if !scalar.EqualWithinRel(s, test.want, 2e-6) {
				t.Errorf("unexpected result eValue(%f, %d, %d): got %f want %f", test.t, test.n1, test.n2, s, test.want)
			}
		})
	}
}

func TestCI(t *testing.T) {
	t.Parallel()
	data := grayData[slices.IndexFunc(grayData, func(d []grayCase) bool { return d[0].location == "Carleton University, Ottawa, Canada" })]
	tests := []struct {
		n  int
		c0 float64
		c1 float64
	}{
		{n: 1, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 2, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 3, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 4, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 5, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 6, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 7, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 8, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 9, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 10, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 11, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 12, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 13, c0: math.Inf(-1), c1: math.Inf(1)},
		{n: 14, c0: -98.62792, c1: 101.9612},
		{n: 15, c0: -15.42631, c1: 19.09297},
		{n: 16, c0: -6.51332, c1: 9.346654},
		{n: 17, c0: -3.997455, c1: 7.930789},
		{n: 18, c0: -3.390325, c1: 7.451863},
		{n: 19, c0: -2.914129, c1: 7.085558},
		{n: 20, c0: -2.54023, c1: 6.94023},
		{n: 21, c0: -2.218689, c1: 6.818689},
		{n: 22, c0: -1.32434, c1: 6.32434},
		{n: 23, c0: -0.8378349, c1: 5.837835},
		{n: 24, c0: -0.6311292, c1: 5.381129},
		{n: 25, c0: -0.2024012, c1: 5.424623},
		{n: 26, c0: -0.2640066, c1: 5.064007},
		{n: 27, c0: -0.2703812, c1: 4.893911},
		{n: 28, c0: -0.3181225, c1: 4.596197},
		{n: 29, c0: -0.1938408, c1: 4.537275},
		{n: 30, c0: 0.03226012, c1: 4.523295},
		{n: 31, c0: 0.1665839, c1: 4.552714},
		{n: 32, c0: 0.1053324, c1: 4.331915},
		{n: 33, c0: 0.146038, c1: 4.238577},
		{n: 34, c0: 0.258684, c1: 4.268789},
		{n: 35, c0: 0.4691052, c1: 4.388038},
		{n: 36, c0: 0.4001898, c1: 4.209334},
		{n: 37, c0: 0.4728427, c1: 4.175642},
		{n: 38, c0: 0.204692, c1: 3.977126},
		{n: 39, c0: 0.242893, c1: 3.880102},
		{n: 40, c0: 0.02959331, c1: 3.70778},
		{n: 41, c0: 0.09503225, c1: 3.677915},
		{n: 42, c0: 0.187117, c1: 3.661853},
		{n: 43, c0: 0.1630632, c1: 3.554328},
		{n: 44, c0: 0.2437182, c1: 3.545102},
		{n: 45, c0: 0.2999012, c1: 3.521527},
		{n: 46, c0: 0.2753359, c1: 3.42921},
		{n: 47, c0: 0.1695833, c1: 3.324962},
		{n: 48, c0: 0.01554447, c1: 3.198369},
		{n: 49, c0: 0.09703757, c1: 3.230722},
		{n: 50, c0: -0.03914359, c1: 3.109656},
		{n: 51, c0: -0.09881216, c1: 3.012658},
		{n: 52, c0: -0.101948, c1: 2.948102},
		{n: 53, c0: -0.02592925, c1: 2.977496},
		{n: 54, c0: -0.0197774, c1: 2.926371},
		{n: 55, c0: -0.02336174, c1: 2.867277},
		{n: 56, c0: -0.02621921, c1: 2.811933},
		{n: 57, c0: -0.02844918, c1: 2.759976},
		{n: 58, c0: -0.03013572, c1: 2.711088},
		{n: 59, c0: 0.02910432, c1: 2.733568},
		{n: 60, c0: 0.02600273, c1: 2.688283},
		{n: 61, c0: 0.05555108, c1: 2.673886},
		{n: 62, c0: 0.05803339, c1: 2.629532},
		{n: 63, c0: 0.0606844, c1: 2.5878},
		{n: 64, c0: 0.09619272, c1: 2.580249},
		{n: 65, c0: 0.164917, c1: 2.643432},
		{n: 66, c0: 0.2235697, c1: 2.673489},
		{n: 67, c0: 0.2893588, c1: 2.73207},
		{n: 68, c0: 0.3201729, c1: 2.723983},
		{n: 69, c0: 0.3191499, c1: 2.687573},
		{n: 70, c0: 0.3184285, c1: 2.653},
		{n: 71, c0: 0.3934321, c1: 2.751012},
		{n: 72, c0: 0.3192517, c1: 2.680748},
		{n: 73, c0: 0.3795185, c1: 2.730091},
		{n: 74, c0: 0.3659826, c1: 2.689573},
		{n: 75, c0: 0.3947169, c1: 2.686364},
		{n: 76, c0: 0.4634515, c1: 2.771476},
		{n: 77, c0: 0.4909541, c1: 2.768155},
		{n: 78, c0: 0.5380819, c1: 2.795251},
		{n: 79, c0: 0.5372391, c1: 2.766607},
		{n: 80, c0: 0.5069519, c1: 2.720065},
		{n: 81, c0: 0.5503801, c1: 2.746323},
		{n: 82, c0: 0.4655226, c1: 2.677335},
		{n: 83, c0: 0.5074294, c1: 2.701873},
		{n: 84, c0: 0.547685, c1: 2.725042},
		{n: 85, c0: 0.5863884, c1: 2.746945},
		{n: 86, c0: 0.5733911, c1: 2.711162},
		{n: 87, c0: 0.592393, c1: 2.706653},
		{n: 88, c0: 0.6058952, c1: 2.694312},
		{n: 89, c0: 0.6240386, c1: 2.690044},
		{n: 90, c0: 0.6579574, c1: 2.703745},
		{n: 91, c0: 0.6927843, c1: 2.723882},
		{n: 92, c0: 0.6879246, c1: 2.69983},
		{n: 93, c0: 0.7209153, c1: 2.719085},
		{n: 94, c0: 0.7527797, c1: 2.737416},
		{n: 95, c0: 0.7835767, c1: 2.754885},
		{n: 96, c0: 0.7536425, c1: 2.718056},
		{n: 97, c0: 0.7403016, c1: 2.685942},
		{n: 98, c0: 0.7542303, c1: 2.681797},
		{n: 99, c0: 0.7677302, c1: 2.677724},
		{n: 100, c0: 0.8158708, c1: 2.719483},
		{n: 101, c0: 0.8458001, c1: 2.732066},
		{n: 102, c0: 0.7919001, c1: 2.687914},
		{n: 103, c0: 0.8054778, c1: 2.684987},
		{n: 104, c0: 0.8505302, c1: 2.723191},
		{n: 105, c0: 0.9067115, c1: 2.788464},
		{n: 106, c0: 0.9019647, c1: 2.769012},
		{n: 107, c0: 0.9293165, c1: 2.786079},
		{n: 108, c0: 0.9558729, c1: 2.80246},
		{n: 109, c0: 0.9622272, c1: 2.790154},
		{n: 110, c0: 0.9463786, c1: 2.760288},
		{n: 111, c0: 0.9712439, c1: 2.769933},
		{n: 112, c0: 0.9316632, c1: 2.727311},
		{n: 113, c0: 0.9261686, c1: 2.708762},
		{n: 114, c0: 0.9010986, c1: 2.678306},
		{n: 115, c0: 0.896407, c1: 2.66098},
		{n: 116, c0: 0.8821625, c1: 2.63386},
		{n: 117, c0: 0.9070171, c1: 2.650176},
		{n: 118, c0: 0.8933751, c1: 2.623986},
		{n: 119, c0: 0.8479675, c1: 2.58622},
		{n: 120, c0: 0.8438691, c1: 2.570272},
		{n: 121, c0: 0.8308287, c1: 2.544929},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			data = data[:test.n]
			x, y := splitGray(data)
			if !(len(x) > 1 && len(y) > 1) {
				return
			}

			const alpha = 0.05
			p := &Mom{G: 0.1339827}
			ci := p.CI(x, y, alpha)
			if !scalar.EqualWithinAbsOrRel(ci[0], test.c0, 5e-5, 5e-5) {
				t.Errorf("unexpected result: got %f want %f", ci[0], test.c0)
			}
			if !scalar.EqualWithinAbsOrRel(ci[1], test.c1, 5e-5, 5e-5) {
				t.Errorf("unexpected result: got %f want %f", ci[1], test.c1)
			}
		})
	}
}

func TestGetNPlan(t *testing.T) {
	t.Parallel()
	tests := []struct {
		alpha    float64
		beta     float64
		deltaMin float64
		nPlan    NPlan
	}{
		{alpha: 0.05, beta: 0.2, deltaMin: 0.51765, nPlan: NPlan{N: 94, Mean: 56, Batch: 113}},
		{alpha: 0.01, beta: 0.1, deltaMin: 0.7688172, nPlan: NPlan{N: 75, Mean: 41, Batch: 83}},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			nPlan := GetNPlan(test.alpha, test.beta, test.deltaMin)
			if nPlan.N != test.nPlan.N {
				t.Errorf("GetNPlan(%f, %f, %f).N: got %d want %d", test.alpha, test.beta, test.deltaMin, nPlan.N, test.nPlan.N)
			}
			if nPlan.Mean != test.nPlan.Mean {
				t.Errorf("GetNPlan(%f, %f, %f).Mean: got %d want %d", test.alpha, test.beta, test.deltaMin, nPlan.Mean, test.nPlan.Mean)
			}
			if nPlan.Batch != test.nPlan.Batch {
				t.Errorf("GetNPlan(%f, %f, %f).Batch: got %d want %d", test.alpha, test.beta, test.deltaMin, nPlan.Batch, test.nPlan.Batch)
			}
		})
	}
}

// Downloaded from https://github.com/ManyLabsOpenScience/ManyLabs2/blob/master/OSFdata/Moral%20Typecasting%20(Gray%20%26%20Wegner%2C%202009)/Gray.1/Global/Data/Gray_1_study_global_include_all_CLEAN_CASE.csv
//
//go:embed testdata/Gray_1_study_global_include_all_CLEAN_CASE.csv
var Gray_1_study_global_include_all_CLEAN_CASE []byte

var grayData = getGray()

const adultHarmsBaby = "Adult harms Baby"

type grayCase struct {
	uID      int
	variable int
	factor   string
	source   string
	location string
}

func getGray() [][]grayCase {
	rows, err := csv.NewReader(bytes.NewBuffer(Gray_1_study_global_include_all_CLEAN_CASE)).ReadAll()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	// Skip header.
	rows = rows[1:]

	dataM := make(map[string][]grayCase)
	for _, row := range rows {
		var c grayCase
		c.uID, err = strconv.Atoi(row[0])
		if err != nil {
			log.Fatalf("%+v", err)
		}
		c.variable, err = strconv.Atoi(row[1])
		if err != nil {
			log.Fatalf("%+v", err)
		}
		c.factor = row[2]
		c.source = row[8] + row[9] + row[10]
		c.location = row[12]

		dataM[c.source] = append(dataM[c.source], c)
	}

	// Sort data by collection time for each source.
	data := make([][]grayCase, 0)
	for _, v := range dataM {
		slices.SortFunc(v, func(a, b grayCase) int { return cmp.Compare(a.uID, b.uID) })
		data = append(data, v)
	}

	// Sort by source name.
	slices.SortFunc(data, func(a, b []grayCase) int { return cmp.Compare(a[0].source, b[0].source) })

	return data
}

func splitGray(data []grayCase) ([]float64, []float64) {
	var x, y []float64
	for _, d := range data {
		if d.factor == adultHarmsBaby {
			x = append(x, float64(d.variable))
		} else {
			y = append(y, float64(d.variable))
		}
	}
	return x, y
}

func TestMain(m *testing.M) {
	flag.Parse()
	log.SetFlags(log.Lmicroseconds | log.Llongfile | log.LstdFlags)

	m.Run()
}
