package utils

func InInt64(e int64, s []int64) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func InInt(e int, s []int) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func InString(e string, s []string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
