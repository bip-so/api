package utils

// StringExistsInList: Will check if string exists in a slice of string
func SliceContainsItem(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func RemoveItemAtIndexFromSlice(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func RemoveItemFromSlice(s []string, item string) []string {
	result := []string{}
	for _, element := range s {
		if item != element {
			result = append(result, element)
		}
	}
	return result
}

func GetIntSliceIndex(slice []uint64, item uint64) *uint64 {
	var index *uint64
	for i := range slice {
		if slice[i] == item {
			tempIndex := uint64(i)
			index = &tempIndex
			return index
		}
	}
	return index
}

func MergeMaps(ms ...map[uint64]map[uint64]string) map[uint64]map[uint64]string {
	res := map[uint64]map[uint64]string{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}
