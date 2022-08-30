package utils

import (
	"log"
)

// Check logs error and exits program
func Check(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

func StringOr(str1, str2 string) string {
	if str1 != "" {
		return str2
	}
	return str2
}

func IntOr(int1, int2 int) int {
	if int1 > 0 {
		return int1
	}
	return int2
}

func VarOr(var1, var2 any) any {
	if var1 != nil {
		return var1
	}
	return var2
}
