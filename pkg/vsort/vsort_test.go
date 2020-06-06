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

func TestSorterCompare(t *testing.T) {
	type Case struct {
		options  []Option
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
		{
			options:  []Option{WithPrefix("v")},
			v1:       "v0.1.1",
			v2:       "v0.1.0",
			expected: 1,
		},
		{
			options:  []Option{WithPrefix("[[:alpha:]]+-")},
			v1:       "rc-0.1.1",
			v2:       "version-0.1.0",
			expected: 1,
		},
		{
			options:  []Option{WithSuffix(`-\d+`)},
			v1:       "0.1.1-1",
			v2:       "0.1.0-2",
			expected: 1,
		},
		{
			options:  []Option{WithLevel(2)},
			v1:       "0.10",
			v2:       "0.1",
			expected: 1,
		},
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
			s, err := NewSorter(tt.options...)
			if assert.NoError(t, err) {
				actual, err := s.Compare(tt.v1, tt.v2)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.expected, actual)
				}
			}
		})
	}
}

func TestSorterSort(t *testing.T) {
	type Case struct {
		versions []string
		options  []Option
		expected []string
	}
	cases := []Case{
		{
			versions: []string{"0.2.0", "0.0.1", "0.10.0", "0.0.2"},
			options:  []Option{},
			expected: []string{"0.0.1", "0.0.2", "0.2.0", "0.10.0"},
		},
		{
			versions: []string{"0.2.0", "0.0.1", "0.10.0", "0.0.2"},
			options:  []Option{WithOrder(Desc)},
			expected: []string{"0.10.0", "0.2.0", "0.0.2", "0.0.1"},
		},
		{
			versions: []string{"v0.2.0", "v0.0.1", "v0.10.0", "v0.0.2"},
			options:  []Option{WithPrefix("v")},
			expected: []string{"v0.0.1", "v0.0.2", "v0.2.0", "v0.10.0"},
		},
		{
			versions: []string{"2.0", "0.1", "10.0", "0.2"},
			options:  []Option{WithLevel(2)},
			expected: []string{"0.1", "0.2", "2.0", "10.0"},
		},
	}

	genSubtestName := func(c Case) string {
		return fmt.Sprintf("%q(%s)", c.versions, c.options)
	}

	for _, tt := range cases {
		t.Run(genSubtestName(tt), func(t *testing.T) {
			copied := make([]string, len(tt.versions))
			copy(copied, tt.versions)

			s, err := NewSorter(tt.options...)
			if assert.NoError(t, err) {
				s.Sort(copied)
				assert.Equal(t, tt.expected, copied)
			}
		})
	}
}

func TestSorterIsValid(t *testing.T) {
	type versionCase struct {
		version  string
		expected bool
	}
	cases := []struct {
		options  []Option
		versions []versionCase
	}{
		{
			options: []Option{},
			versions: []versionCase{
				{
					version:  "0.1.0",
					expected: true,
				},
				{
					version:  "1.0",
					expected: true,
				},
				{
					version:  "v0.1.0",
					expected: false,
				},
				{
					version:  "0.1.0-rc1",
					expected: false,
				},
			},
		},
		{
			options: []Option{WithPrefix("v")},
			versions: []versionCase{
				{
					version:  "v0.1.0",
					expected: true,
				},
				{
					version:  "v1.0",
					expected: true,
				},
				{
					version:  "0.1.0",
					expected: false,
				},
			},
		},
		{
			options: []Option{WithPrefix("[[:alpha:]]+-")},
			versions: []versionCase{
				{
					version:  "version-0.1.0",
					expected: true,
				},
				{
					version:  "v1.0",
					expected: false,
				},
				{
					version:  "0.1.0",
					expected: false,
				},
			},
		}, {
			options: []Option{WithSuffix(`-\d+`)},
			versions: []versionCase{
				{
					version:  "0.1.0-1",
					expected: true,
				},
				{
					version:  "1.0-10",
					expected: true,
				},
				{
					version:  "0.1.0-alpha",
					expected: false,
				},
				{
					version:  "0.1.0",
					expected: false,
				},
			},
		},
		{
			options: []Option{WithLevel(3)},
			versions: []versionCase{
				{
					version:  "0.1.0",
					expected: true,
				},
				{
					version:  "1.0",
					expected: false,
				},
				{
					version:  "0.1.0.1",
					expected: false,
				},
			},
		},
	}

	genSubtestName := func(options []Option, c versionCase) string {
		return fmt.Sprintf("%q(%s)", c.version, options)
	}

	for _, c := range cases {
		options := c.options
		for _, tt := range c.versions {
			t.Run(genSubtestName(options, tt), func(t *testing.T) {
				s, err := NewSorter(options...)
				if assert.NoError(t, err) {
					assert.Equal(t, tt.expected, s.IsValid(tt.version))
				}
			})
		}
	}
}
