package paypal

import "fmt"

func (s *Service) convertToFormData(param map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range param {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
