package kerlong

import "fmt"

//CheckErr 有错误就显示
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
