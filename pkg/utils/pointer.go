package utils

func SlicePtr(s []string) *[]string {
	if s == nil {
		empty := []string{}
		return &empty
	}

	return &s
}
