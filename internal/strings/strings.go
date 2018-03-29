package strings

// InSlice checks if the given string exists in the given slice:
func InSlice(str string, list []string) bool {
	for _, s := range list {
		if s == str {
			return true
		}
	}
	return false
}
