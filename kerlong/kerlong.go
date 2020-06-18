package kerlong

//ternary 三目运算符
func ternary(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}
