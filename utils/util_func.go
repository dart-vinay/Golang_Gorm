package utils

func Unique(e []int) []int {
	r := []int{}
	for _, s := range e {
		if !Contains(r[:], s) {
			r = append(r, s)
		}
	}
	return r
}

func Contains(e []int, c int) bool {
	if c == 0 {
		return false
	}
	for _, s := range e {
		if s == c {
			return true
		}
	}
	return false
}

func ContainsString(e []string, c string) bool {
	if c == "" {
		return false
	}
	for _, s := range e {
		if s == c {
			return true
		}
	}
	return false
}
