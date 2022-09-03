package utils

import "log"

// Check logs error and exits program
func Check(errs ...error) {
	for _, err := range errs {
		if err != nil {
			log.Fatalf(err.Error())
		}
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
