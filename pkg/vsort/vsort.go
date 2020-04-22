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
	"sort"
	"strconv"
	"strings"
)

// Sorter provides comparation and sorting versions
type Sorter interface {
	Compare(v1, v2 string) (int, error)
	Sort(versions []string)
	IsValid(v string) bool
}

type order int

const (
	// Asc should be passed to Sort via `WithOrder(Asc)`
	Asc order = iota
	// Desc should be passed to Sort via `WithOrder(Desc)`
	Desc
)

type sorter struct {
	order  order
	prefix string
	level  int
}

// Option is Functional optional pattern object for Sort
type Option interface {
	apply(*sorter)
}

// WithOrder represent order of Sort
type WithOrder order

func (o WithOrder) apply(s *sorter) {
	s.order = order(o)
}

// String returns "Asc" or "Desc"
func (o WithOrder) String() string {
	switch order(o) {
	case Asc:
		return "order=asc"
	case Desc:
		return "order=Desc"
	default:
		return "order=unknown"
	}
}

// WithPrefix represent expected prefix of version string
type WithPrefix string

func (p WithPrefix) apply(s *sorter) {
	s.prefix = string(p)
}

func (p WithPrefix) String() string {
	return "prefix=" + string(p)
}

// WithLevel represent expected level of version string
type WithLevel int

func (l WithLevel) apply(s *sorter) {
	s.level = int(l)
}

func (l WithLevel) String() string {
	return fmt.Sprintf("level=%d", int(l))
}

// NewSorter returns Sorter initialized by given options
func NewSorter(options ...Option) Sorter {
	defaults := []Option{WithLevel(-1)}
	s := new(sorter)
	for _, o := range append(defaults, options...) {
		o.apply(s)
	}
	return s
}

// Compare returns an integer comparing two version strings.
// The result will be 0 if v1==v2, -1 if v1 < v2, and +1 if v1 > v2.
func (s *sorter) Compare(v1, v2 string) (int, error) {
	nums1 := strings.SplitN(strings.TrimPrefix(v1, s.prefix), ".", s.level)
	nums2 := strings.SplitN(strings.TrimPrefix(v2, s.prefix), ".", s.level)

	for i := 0; i < len(nums1); i++ {
		num1, err := strconv.Atoi(nums1[i])
		if err != nil {
			return 0, err
		}
		num2, err := strconv.Atoi(nums2[i])
		if err != nil {
			return 0, err
		}

		if num1 > num2 {
			return 1, nil
		} else if num1 < num2 {
			return -1, nil
		}
	}

	return 0, nil
}

// Sort sorts given versions
func (s *sorter) Sort(versions []string) {
	sort.Slice(versions, func(i, j int) bool {
		r, _ := s.Compare(versions[i], versions[j])
		if s.order == Asc {
			return r < 0
		}
		return r > 0
	})
}

// IsValid reports whether its argument v is a valid version string.
func (s *sorter) IsValid(v string) bool {
	// check prefix
	if !strings.HasPrefix(v, s.prefix) {
		return false
	}
	v = v[len(s.prefix):]

	// check level
	nums := strings.Split(v, ".")
	if s.level > 0 && len(nums) != s.level {
		return false
	}

	// check each level format
	for _, n := range nums {
		if _, err := strconv.Atoi(n); err != nil || n[0] == '+' || n[0] == '-' {
			return false
		}
	}

	return true
}
