package combat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexSlicing(t *testing.T) {
	sliceIndex := func(o []string, i int) []string {
		return append(o[:i], o[i+1:]...)
	}
	d := sliceIndex([]string{"1", "2", "3"}, 1)
	assert.True(t, len(d) == 2)
	assert.Equal(t, "1", d[0])
	assert.Equal(t, "3", d[1])
	d = sliceIndex([]string{"1", "2", "3"}, 0)
	assert.True(t, len(d) == 2)
	assert.Equal(t, "2", d[0])
	assert.Equal(t, "3", d[1])
	d = sliceIndex([]string{"1", "2", "3"}, 2)
	assert.True(t, len(d) == 2)
	assert.Equal(t, "1", d[0])
	assert.Equal(t, "2", d[1])

	orig := []string{"1", "2", "3"}
	newArr := []string{orig[1]}
	orig = sliceIndex(orig, 1)
	assert.True(t, len(newArr) == 1)
	assert.True(t, len(orig) == 2)
	assert.Equal(t, "2", newArr[0])
	assert.Equal(t, "1", orig[0])
	assert.Equal(t, "3", orig[1])

	var arr []string
	assert.True(t, arr == nil)
	assert.True(t, len(append(arr, "1")) == 1)
}
