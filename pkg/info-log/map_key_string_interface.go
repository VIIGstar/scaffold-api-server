package info_log

func ErrorToLogFields(key string, err error) map[string]interface{} {
	return map[string]interface{}{
		key: err,
	}
}
