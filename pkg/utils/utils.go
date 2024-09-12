package utils

func RemoveZeroValues(slice []string) []string {
	var result []string
	for _, value := range slice {
		if value != "" {
			result = append(result, value)
		}
	}

	return result
}
