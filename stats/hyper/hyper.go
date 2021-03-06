package hyper

import (
	. "math"
)

func lnG(v int) float64 {
	r, _ := Lgamma(float64(v))
	return r
}

// this is for reference
// as the code below is quite unreadable, but 2x as fast
func ComplementCdfSlow(o, r, O, R int) float64 {
	total := 0.0
	lSOR := lnG(O+1) + lnG(R+1)
	lOR := lnG(O + R + 1)
	for r >= 0 {
		nom := lSOR + lnG(o+r+1) + lnG(O+R-o-r+1)
		denom := lnG(o+1) + lnG(O-o+1) + lnG(r+1) + lnG(R-r+1) + lOR
		add := Exp(nom - denom)
		total += add
		r--
		o++
	}
	return total
}

// returns probability of split of
// o - observed in input , r - observed in validation set
// O - total items in input, R - total items in validation set
// using logarithmic gamma function
// TODO: limits test
func ComplementCdf(chosenA, chosenB, totalA, totalB int) float64 {
	total := 0.0

	o := float64(chosenA)
	r := float64(chosenB)
	O := float64(totalA)
	R := float64(totalB)

	gO, _ := Lgamma(O + 1.0)
	gR, _ := Lgamma(R + 1.0)
	gaOR := gO + gR
	gOR, _ := Lgamma(O + R + 1.0)
	for r >= 0.0 {
		gor, _ := Lgamma(o + r + 1.0)
		gORor, _ := Lgamma(O + R - o - r + 1)
		nom := gaOR + gor + gORor

		ga, _ := Lgamma(o + 1.0)
		gOo, _ := Lgamma(O - o + 1.0)
		gr, _ := Lgamma(r + 1.0)
		gRr, _ := Lgamma(R - r + 1.0)

		denom := ga + gOo + gr + gRr + gOR

		total += Exp(nom - denom)
		r--
		o++
	}
	return total
}

// returns probability of split of
// chosenA - observed in input , chosenB - observed in validation set
// totalA - total items in input, totalB - total items in validation set
// using logarithmic gamma function
// TODO: limits test
const approxEpsilon = 1e-6

func ComplementCdfApprox(chosenA, chosenB, totalA, totalB int) float64 {
	total := 0.0

	o := float64(chosenA)
	r := float64(chosenB)
	O := float64(totalA)
	R := float64(totalB)

	gO, _ := Lgamma(O + 1.0)
	gR, _ := Lgamma(R + 1.0)
	gaOR := gO + gR
	gOR, _ := Lgamma(O + R + 1.0)
	for r >= 0.0 {
		gor, _ := Lgamma(o + r + 1.0)
		gORor, _ := Lgamma(O + R - o - r + 1)
		nom := gaOR + gor + gORor

		ga, _ := Lgamma(o + 1.0)
		gOo, _ := Lgamma(O - o + 1.0)
		gr, _ := Lgamma(r + 1.0)
		gRr, _ := Lgamma(R - r + 1.0)

		denom := ga + gOo + gr + gRr + gOR

		add := Exp(nom - denom)
		total += add

		if add < total*approxEpsilon {
			break
		}

		r--
		o++
	}
	return total
}

// returns probability of split of
// o - observed in input , r - observed in validation set
// O - total items in input, R - total items in validation set
// using logarithmic gamma function
func Cdf(chosenA, chosenB, totalA, totalB int) float64 {
	total := 0.0

	o := float64(chosenA)
	r := float64(chosenB)
	O := float64(totalA)
	R := float64(totalB)

	gO, _ := Lgamma(O + 1.0)
	gR, _ := Lgamma(R + 1.0)
	gaOR := gO + gR
	gOR, _ := Lgamma(O + R + 1.0)
	for o >= 0.0 {
		gor, _ := Lgamma(o + r + 1.0)
		gORor, _ := Lgamma(O + R - o - r + 1)
		nom := gaOR + gor + gORor

		ga, _ := Lgamma(o + 1.0)
		gOo, _ := Lgamma(O - o + 1.0)
		gr, _ := Lgamma(r + 1.0)
		gRr, _ := Lgamma(R - r + 1.0)

		denom := ga + gOo + gr + gRr + gOR
		total += Exp(nom - denom)
		r++
		o--
	}
	return total
}
