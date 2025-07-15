package auth

func containsRole(s []string, role string) bool {
	for _, r := range s {
		if r == role {
			return true
		}
	}
	return false
}
