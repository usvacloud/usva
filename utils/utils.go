package utils

import "log"

// Must will iterate through it's arguments,
// log the first error and exit the program with log.Fatalf()
func Must(errs ...error) {
	for _, err := range errs {
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

// StringOr will return the first argument which is not an
// empty string
func StringOr(str1, str2 string) string {
	if str1 != "" {
		return str1
	}
	return str2
}

func IntOr[T ~int | ~uint](int1, int2 T) T {
	if int1 != 0 {
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

func IsBetween[T ~int | ~uint](base, min, max T) bool {
	return min <= base && base <= max
}
