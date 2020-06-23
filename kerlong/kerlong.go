package kerlong

//IntAbs 绝对值
func IntAbs(num int) int {
	ans, ok := Ternary(num > 0, num, -num).(int)
	if ok {
		return ans
	}
	return 0
}

//Ternary 三目运算符
func Ternary(b bool, t, f interface{}) interface{} {
	if b {
		return t
	}
	return f
}
