package utils

func ToString(v interface{}) string {
	if val, ok := v.(string); ok {
		return val
	}
	return ""
}

func ToInt(v interface{}) int {
	if val, ok := v.(float64); ok {
		return int(val)
	}
	return 0
}

func ToFloat(v interface{}) float64 {
	if val, ok := v.(float64); ok {
		return val
	}
	return 0
}

func ToBool(v interface{}) bool {
	if val, ok := v.(bool); ok {
		return val
	}
	return false
}
