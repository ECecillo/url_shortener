package internal

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratorReturnSpecifiedLengthValue(t *testing.T) {

	testCases := []struct {
		desc      string
		inputSize int
	}{
		{
			desc:      "with size of 10",
			inputSize: 10,
		},
		{
			desc:      "with size of 1000",
			inputSize: 1000,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			value, err := Generate(tC.inputSize)
			expectedLength := tC.inputSize

			require.NoError(t, err)
			require.Len(t, value, expectedLength)

		})
	}

}

func TestGeneratorUniqueness(t *testing.T) {

	testCases := []struct {
		desc      string
		valueSize int
		count     int
	}{
		{
			desc:      "after 100 iterations",
			valueSize: 10,
			count:     100,
		},
		{
			desc:      "after 1000 iterations",
			valueSize: 10,
			count:     1000,
		},
		{
			desc:      "after 10 000 iterations",
			valueSize: 10,
			count:     10000,
		},
	}
	for _, tC := range testCases {
		seenValues := make(map[string]bool)

		for range tC.count {
			t.Run(tC.desc, func(t *testing.T) {
				value, err := Generate(tC.valueSize)

				require.NoError(t, err)
				require.False(t, seenValues[value], "value %s already in generated", value)
				seenValues[value] = true
			})
		}
	}

}

func TestGeneratorDistribution(t *testing.T) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const numSamples = 10000
	const strLen = 100

	freq := make(map[rune]int)
	for _, c := range alphabet {
		freq[c] = 0
	}

	for i := 0; i < numSamples; i++ {
		value, err := Generate(strLen)
		require.NoError(t, err)

		for _, c := range value {

			_, ok := freq[c]
			require.True(t, ok, "unknown character %s given, value : %s", c, value)

			freq[c]++
		}

		chi2 := calculateChi2(t, numSamples, strLen, alphabet, freq)
		assertChi2Distribution(t, chi2, alphabet)
	}

}

// calculateChi2 calculates the chi² statistic for the observed character frequencies
// from the generated string, comparing them to the expected frequency under a uniform distribution.
// It uses the formula: chi² = Σ((observed - expected)² / expected)
//
// For more details, see: https://www.southampton.ac.uk/passs/full_time_education/bivariate_analysis/chi_square.page
func calculateChi2(t testing.TB, numSamples int, strLen int, alphabet string, charFrequencies map[rune]int) float64 {

	t.Helper()

	alphabetSize := len(alphabet)

	totalCount := float64(numSamples * strLen)
	expected := totalCount / float64(alphabetSize)

	chi2 := 0.0
	for _, c := range alphabet {
		observed := float64(charFrequencies[c])
		chi2 += math.Pow(observed-expected, 2) / expected
	}

	return chi2
}

// assertChi2Distribution asserts that the chi² statistic for the generated string's character distribution
// is within the expected bounds for a uniform distribution.
//
// The expected bounds are defined as [ df - 3*sqrt(2*df)           ,           df + 3*sqrt(2*df)],
// where df is the degrees of freedom.
func assertChi2Distribution(t testing.TB, chi2 float64, alphabet string) {
	t.Helper()

	alphabetSize := len(alphabet)

	// degrees
	degree := float64(alphabetSize - 1)

	variance := math.Sqrt(2 * degree)
	stdDeviation := 3 * variance

	lowerBound := degree - stdDeviation
	upperBound := degree + stdDeviation

	require.Less(t, chi2, upperBound, "value greater than upperbound, chi² = %f, expected value between %f and %f", chi2, lowerBound, upperBound)
	require.Greater(t, chi2, lowerBound, "value lower than lowerbound, chi² = %f, expected value between %f and %f", chi2, lowerBound, upperBound)
}
