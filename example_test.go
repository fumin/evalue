package evalue_test

import (
	"fmt"

	"github.com/fumin/evalue"
)

func Example() {
	// Prepare data.
	// Data is from the Moral Typecasting study of the Many Labs 2 project (Klein et al., 2018).
	// For simplicity, data is further filtered to contain only those collected from Carleton University, Ottawa, Canada.
	//
	// Klein RA, Vianello M, Hasselman F, et al. Many Labs 2: Investigating Variation in Replicability Across Samples and Settings. Advances in Methods and Practices in Psychological Science. 2018;1(4):443-490. doi:10.1177/2515245918810225
	data := getData()

	// Design the experiment.
	// The null hypothesis is that the difference between the mean of the two groups is zero.
	// The statistical test to be performed is an e-value based two sample t-test.
	// Set the significance level at the usual 5%.
	alpha := 0.05
	// Create an e-process that is tuned to reject the null hypothesis at the fastest rate.
	// delta is the minimal clinically relevant standardized effect size.
	delta := 0.5176537
	eProcess := evalue.NewMom(delta)

	// Perform the e-value based test while the experiment is running.
	// In contrast to p-values, e-values control the Type I error at all times, and thus allow optional stopping.
	stoppingTime := -1
	for n := range data {
		// Prepare group data for the two sample test.
		group1, group2 := splitGroups(data[:n])
		if !(len(group1) > 1 && len(group2) > 1) {
			continue
		}

		// Perform the e-value based test with optional stopping.
		eValue := eProcess.EValue(group1, group2)
		if eValue > 1./alpha {
			stoppingTime = n
			break
		}
	}

	needed := 100 * float64(stoppingTime) / float64(len(data))
	fmt.Printf("Null hypothesis rejected with only %.0f%% (%d/%d) of the data needed.\n", needed, stoppingTime, len(data))
	// Output:
	// Null hypothesis rejected with only 25% (30/121) of the data needed.
}

func ExampleGetNPlan() {
	// alpha is our desired significance level.
	alpha := 0.05
	// power is our desired statistical power.
	power := 0.8
	// deltaMin is a lower bound of the true effect size based on domain knowledge.
	deltaMin := 0.5

	nplan := evalue.GetNPlan(alpha, 1-power, deltaMin)
	fmt.Printf("We should plan for a sample size of %d, which would achieve both our desired statistical power and significance level.\n", nplan.N)
	fmt.Printf("The good news is we can early-stop our experiment with only %d samples on average, which is not allowed with conventional t-tests.\n", nplan.Mean)

	// Output:
	// We should plan for a sample size of 102, which would achieve both our desired statistical power and significance level.
	// The good news is we can early-stop our experiment with only 61 samples on average, which is not allowed with conventional t-tests.
}

type datum struct {
	group int
	value float64
}

func splitGroups(data []datum) ([]float64, []float64) {
	var g1, g2 []float64
	for _, d := range data {
		if d.group == 1 {
			g1 = append(g1, d.value)
		} else {
			g2 = append(g2, d.value)
		}
	}
	return g1, g2
}

func getData() []datum {
	return []datum{
		datum{group: 1, value: 2},
		datum{group: 1, value: 3},
		datum{group: 1, value: 7},
		datum{group: 1, value: 5},
		datum{group: 1, value: 6},
		datum{group: 2, value: 4},
		datum{group: 1, value: 7},
		datum{group: 1, value: 3},
		datum{group: 2, value: 1},
		datum{group: 1, value: 6},
		datum{group: 1, value: 7},
		datum{group: 1, value: 2},
		datum{group: 1, value: 7},
		datum{group: 2, value: 5},
		datum{group: 1, value: 7},
		datum{group: 2, value: 5},
		datum{group: 2, value: 1},
		datum{group: 1, value: 6},
		datum{group: 1, value: 6},
		datum{group: 1, value: 7},
		datum{group: 1, value: 7},
		datum{group: 2, value: 2},
		datum{group: 2, value: 3},
		datum{group: 2, value: 4},
		datum{group: 2, value: 1},
		datum{group: 2, value: 5},
		datum{group: 1, value: 4},
		datum{group: 2, value: 5},
		datum{group: 1, value: 6},
		datum{group: 2, value: 2},
		datum{group: 1, value: 7},
		datum{group: 2, value: 5},
		datum{group: 1, value: 5},
		datum{group: 1, value: 7},
		datum{group: 2, value: 1},
		datum{group: 2, value: 5},
		datum{group: 1, value: 6},
		datum{group: 2, value: 7},
		datum{group: 2, value: 4},
		datum{group: 2, value: 7},
		datum{group: 1, value: 6},
		datum{group: 2, value: 3},
		datum{group: 2, value: 5},
		datum{group: 2, value: 3},
		datum{group: 1, value: 6},
		datum{group: 2, value: 5},
		datum{group: 1, value: 3},
		datum{group: 2, value: 7},
		datum{group: 1, value: 7},
		datum{group: 2, value: 7},
		datum{group: 2, value: 6},
		datum{group: 2, value: 5},
		datum{group: 1, value: 7},
		datum{group: 1, value: 5},
		datum{group: 2, value: 5},
		datum{group: 2, value: 5},
		datum{group: 2, value: 5},
		datum{group: 2, value: 5},
		datum{group: 2, value: 3},
		datum{group: 2, value: 5},
		datum{group: 2, value: 4},
		datum{group: 1, value: 5},
		datum{group: 1, value: 5},
		datum{group: 1, value: 6},
		datum{group: 2, value: 2},
		datum{group: 1, value: 7},
		datum{group: 2, value: 2},
		datum{group: 1, value: 6},
		datum{group: 1, value: 5},
		datum{group: 1, value: 5},
		datum{group: 2, value: 1},
		datum{group: 1, value: 3},
		datum{group: 2, value: 2},
		datum{group: 2, value: 5},
		datum{group: 1, value: 6},
		datum{group: 2, value: 1},
		datum{group: 1, value: 6},
		datum{group: 1, value: 7},
		datum{group: 1, value: 5},
		datum{group: 1, value: 4},
		datum{group: 1, value: 7},
		datum{group: 2, value: 7},
		datum{group: 1, value: 7},
		datum{group: 1, value: 7},
		datum{group: 1, value: 7},
		datum{group: 2, value: 5},
		datum{group: 1, value: 6},
		datum{group: 2, value: 4},
		datum{group: 1, value: 6},
		datum{group: 2, value: 3},
		datum{group: 1, value: 7},
		datum{group: 1, value: 5},
		datum{group: 1, value: 7},
		datum{group: 1, value: 7},
		datum{group: 1, value: 7},
		datum{group: 1, value: 4},
		datum{group: 2, value: 5},
		datum{group: 1, value: 6},
		datum{group: 1, value: 6},
		datum{group: 2, value: 2},
		datum{group: 2, value: 3},
		datum{group: 1, value: 3},
		datum{group: 1, value: 6},
		datum{group: 2, value: 2},
		datum{group: 2, value: 1},
		datum{group: 1, value: 5},
		datum{group: 1, value: 7},
		datum{group: 1, value: 7},
		datum{group: 2, value: 4},
		datum{group: 2, value: 5},
		datum{group: 2, value: 3},
		datum{group: 2, value: 6},
		datum{group: 1, value: 5},
		datum{group: 1, value: 4},
		datum{group: 1, value: 5},
		datum{group: 2, value: 5},
		datum{group: 1, value: 7},
		datum{group: 2, value: 5},
		datum{group: 1, value: 3},
		datum{group: 1, value: 5},
		datum{group: 2, value: 5},
	}
}
