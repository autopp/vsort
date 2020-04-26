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
	"errors"
	"fmt"
	"regexp"
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
	suffix *regexp.Regexp
	level  int
}

// Option is Functional optional pattern object for Sort
type Option interface {
	apply(*sorter) error
}

// WithOrder represent order of Sort
type WithOrder order

func (o WithOrder) apply(s *sorter) error {
	ov := order(o)
	if ov != Asc && ov != Desc {
		return errors.New("order should be one of Asc or Desc")
	}
	s.order = ov

	return nil
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

func (p WithPrefix) apply(s *sorter) error {
	s.prefix = string(p)
	return nil
}

func (p WithPrefix) String() string {
	return "prefix=" + string(p)
}

// WithSuffix represents expected suffix pattern of version string
type WithSuffix string

func (suf WithSuffix) apply(s *sorter) error {
	r, err := regexp.Compile(string(suf) + "$")
	if err != nil {
		return err
	}
	s.suffix = r

	return nil
}

// WithLevel represent expected level of version string
type WithLevel int

func (l WithLevel) apply(s *sorter) error {
	level := int(l)
	if level == 0 {
		return errors.New("level shoud not be zero")
	}

	s.level = level

	return nil
}

func (l WithLevel) String() string {
	return fmt.Sprintf("level=%d", int(l))
}

// NewSorter returns Sorter initialized by given options
func NewSorter(options ...Option) (Sorter, error) {
	defaults := []Option{WithLevel(-1)}
	s := new(sorter)
	for _, o := range append(defaults, options...) {
		if err := o.apply(s); err != nil {
			return nil, err
		}
	}
	return s, nil
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

	if s.suffix != nil {
		loc := s.suffix.FindStringIndex(v)
		if loc == nil {
			return false
		}
		v = v[0:loc[0]]
	}

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
