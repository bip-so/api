package utils

// StringExistsInList: Will check if string exists in a slice of string
func SliceContainsInt(list []uint64, a uint64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Keys(data map[uint64]string) []uint64 {
	keys := make([]uint64, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}

func Values(data map[uint64]string) []string {
	values := make([]string, 0, len(data))
	for _, value := range data {
		values = append(values, value)
	}
	return values
}

func Contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func KeysForNestedMap(data map[uint64]map[uint64]string) []uint64 {
	keys := make([]uint64, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}
