package env

import "strings"

// Filter returns a subset of secrets whose keys match the given namespace prefix.
// If namespace is empty, all secrets are returned.
func Filter(secrets map[string]string, namespace string) map[string]string {
	if namespace == "" {
		result := make(map[string]string, len(secrets))
		for k, v := range secrets {
			result[k] = v
		}
		return result
	}

	prefix := strings.ToUpper(namespace) + "_"
	result := make(map[string]string)
	for k, v := range secrets {
		if strings.HasPrefix(strings.ToUpper(k), prefix) {
			result[k] = v
		}
	}
	return result
}
