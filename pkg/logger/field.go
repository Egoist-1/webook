package logger

type Field struct {
	Key string
	Val any
}

func Error(err error) Field {
	return Field{
		Key: "error",
		Val: err,
	}
}
func Int(key string, val int) Field {
	return Field{
		Key: key,
		Val: val,
	}
}

func Int64(key string, val int64) Field {
	return Field{
		Key: key,
		Val: val,
	}
}
func String(key string, val string) Field {
	return Field{
		Key: key,
		Val: val,
	}
}
func Bool(key string, val bool) Field {
	return Field{
		Key: key,
		Val: val,
	}
}
