package utils

import "strconv"

// Выполняет проверку номера, заданного в форме строки,
// на корректность в соответствии с алгоритмом Луна.
//
// https://en.wikipedia.org/wiki/Luhn_algorithm
func LuhnCheck(number string) bool {
	sum := 0

	for i := 0; i < len(number); i++ {
		num, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}

		if i%2 == len(number)%2 {
			num *= 2
			if num > 9 {
				num -= 9
			}
		}

		sum += num
	}

	return sum%10 == 0
}
