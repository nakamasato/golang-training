package iteration


func Repeat(str string, repeatTimes int) string {
	var repeated string
	for i := 0; i < repeatTimes; i++ {
		repeated += str
	}
	return repeated
}
