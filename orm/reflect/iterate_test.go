package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateArray(t *testing.T) {
	testCases := []struct {
		name     string
		entity   any
		wantVals []any
		wantErr  error
	}{
		{
			name:     "[]int array",
			entity:   [3]int{1, 2, 3},
			wantVals: []any{1, 2, 3},
			wantErr:  nil,
		},
		{
			name:     "[]int slice",
			entity:   []int{1, 2, 3},
			wantVals: []any{1, 2, 3},
			wantErr:  nil,
		},
	}
	for _, tc := range testCases {
		got, err := IterateArrayOrSlice(tc.entity)
		assert.Equal(t, tc.wantErr, err)
		if err != nil {
			return
		}
		assert.Equal(t, tc.wantVals, got)
	}
}

func TestIterateMap(t *testing.T) {
	testCases := []struct {
		name     string
		entity   any
		wantKeys []any
		wantVals []any
		wantErr  error
	}{
		{
			name:     "map[string]int",
			entity:   map[string]int{"a": 1, "b": 2, "c": 3},
			wantKeys: []any{"a", "b", "c"},
			wantVals: []any{1, 2, 3},
			wantErr:  nil,
		},
	}
	for _, tc := range testCases {
		gotKeys, gotVals, err := IterateMap(tc.entity)
		assert.Equal(t, tc.wantErr, err)
		if err != nil {
			return
		}
		// Map 的 key 顺序遍历出来的顺序不一定是每次都相同的
		assert.Equal(t, tc.wantKeys, gotKeys)
		assert.Equal(t, tc.wantVals, gotVals)
		/*
			assert.Equal(t, len(tc.wantKeys), len(gotKeys))
			for i, k := range tc.wantKeys {
				assert.Equal(t, gotKeys[i], k)
				println("...........")
				fmt.Printf("%v: %v", k, gotVals[i])
				assert.Equal(t, gotVals[i], tc.wantVals[i])
			}
		*/
	}
}
