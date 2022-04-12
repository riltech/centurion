package challenge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultReverseSorter(t *testing.T) {
	assert.True(t, isValidReverseSorterSolution("123456", "654321"))
	assert.False(t, isValidReverseSorterSolution("123456", "65432"))
	assert.False(t, isValidReverseSorterSolution("12345", "12345"))
	assert.False(t, isValidReverseSorterSolution("YoUr_BoY_goT_swAAG", "gaaws_tog_yob_ruoy"))
	assert.True(t, isValidReverseSorterSolution("YoUr_BoY_goT_swAAG", "GAAws_Tog_YoB_rUoY"))
}
