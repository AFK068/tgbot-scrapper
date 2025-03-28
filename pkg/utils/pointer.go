package utils

func SliceStringPtr(s []string) *[]string {
	if s == nil {
		empty := []string{}
		return &empty
	}

	return &s
}

func SliceInt64Ptr(s []int64) *[]int64 {
	if s == nil {
		empty := []int64{}
		return &empty
	}

	return &s
}
