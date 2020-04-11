// Copyright (C) 2020 Akira Tanimura (@autopp)
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vsort

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparatorCompare(t *testing.T) {
	type Case struct {
		v1       string
		v2       string
		expected int
	}
	cases := []Case{
		{v1: "0.1.0", v2: "0.1.0", expected: 0},
		{v1: "0.1.1", v2: "0.1.0", expected: 1},
		{v1: "0.1.0", v2: "0.1.1", expected: -1},
		{v1: "0.1.0", v2: "0.0.1", expected: 1},
		{v1: "0.2.0", v2: "0.10.1", expected: -1},
	}

	genSubtestName := func(c Case) string {
		var op string
		switch c.expected {
		case 0:
			op = "="
		case 1:
			op = ">"
		case -1:
			op = "<"
		}
		return fmt.Sprintf("%q%s%q", c.v1, op, c.v2)
	}

	for _, tt := range cases {
		t.Run(genSubtestName(tt), func(t *testing.T) {
			var c Comparator
			actual, err := c.Compare(tt.v1, tt.v2)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestSort(t *testing.T) {
	type Case struct {
		versions []string
		order    SortOrder
		expected []string
	}
	cases := []Case{
		{
			versions: []string{"0.2.0", "0.0.1", "0.10.0", "0.0.2"},
			order:    Asc,
			expected: []string{"0.0.1", "0.0.2", "0.2.0", "0.10.0"},
		},
		{
			versions: []string{"0.2.0", "0.0.1", "0.10.0", "0.0.2"},
			order:    Desc,
			expected: []string{"0.10.0", "0.2.0", "0.0.2", "0.0.1"},
		},
	}

	genSubtestName := func(c Case) string {
		return fmt.Sprintf("%q(%s)", c.versions, c.order)
	}

	for _, tt := range cases {
		t.Run(genSubtestName(tt), func(t *testing.T) {
			copied := make([]string, len(tt.versions))
			copy(copied, tt.versions)
			Sort(copied, tt.order)
			assert.Equal(t, tt.expected, copied)
		})
	}
}
