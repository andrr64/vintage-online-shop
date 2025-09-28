package helper

// Helper: cek apakah slice string mengandung target
func Contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
