package vsort

import (
	"sort"
	"strconv"
	"strings"
)

// Comparator contains comparation settings and provides Compare(v1, v2)
type Comparator struct {
}

// Compare returns an integer comparing two version strings.
// The result will be 0 if v1==v2, -1 if v1 < v2, and +1 if v1 > v2.
func (*Comparator) Compare(v1, v2 string) (int, error) {
	nums1 := strings.Split(v1, ".")
	nums2 := strings.Split(v2, ".")

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

// SortOrder represent order of Sort
type SortOrder int

const (
	// Asc should be passed to Sort
	Asc = iota
	// Desc should be passed to Sort
	Desc
)

// String returns "Asc" or "Desc"
func (o SortOrder) String() string {
	if o == Asc {
		return "Asc"
	}
	return "Desc"
}

// Sort sorts given versions
func Sort(versions []string, order SortOrder) {
	c := new(Comparator)
	sort.Slice(versions, func(i, j int) bool {
		r, _ := c.Compare(versions[i], versions[j])
		if order == Asc {
			return r < 0
		}
		return r > 0
	})
}
