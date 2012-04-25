package stats

import "math"

func lnG(v int) float64 {
	r, _ := math.Lgamma(float64(v))
	return r
}

func gamma(v int) float64 {
	return math.Gamma(float64(v))
}

// returns probability of split of
// o - observed in input , r - observed in validation set
// O - total items in input, R - total items in validation set
// using logarithmic gamma function
func HypergeometricSplit(o int, r int, O int, R int) float64 {
	nom := lnG(O+1) + lnG(R+1) + lnG(o+r+1) + lnG(O+R-o-r+1)
	denom := lnG(o+1) + lnG(O-o+1) + lnG(r+1) + lnG(R-r+1) + lnG(O+R+1)
	if r > 0 {
		return math.Exp(nom-denom) + HypergeometricSplit(o+1, r-1, O, R)
	}
	return math.Exp(nom - denom)
}

// returns probability of split of
// choosing x items from N items
// p - probability of getting one item
func BinomialProb(x int, N int, p float64) float64 {
	nom := lnG(N + 1)
	denom := lnG(x+1) + lnG(N-x+1)
	return math.Exp(nom-denom) * math.Pow(p, float64(x)) * math.Pow(1-p, float64(N-x))
}
