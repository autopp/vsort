package vsort

import (
	"strconv"
	"strings"
)

func Compare(v1, v2 string) (int, error) {
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
