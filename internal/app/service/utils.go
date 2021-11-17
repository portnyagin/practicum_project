package service

func CheckOrderNum(orderNum string) bool {
	var (
		number int
		check  int
	)
	chars := []rune(orderNum)
	for i := 0; i < len(chars); i++ {
		if !(chars[i] >= '0' && chars[i] <= '9') {
			return false
		}
		number = int(chars[i] - '0')
		if (len(chars)-i)%2 == 0 {
			number *= 2
			if number > 9 {
				number -= 9
			}
		}
		check += number
	}

	return (check*9)%10 == 0
}
