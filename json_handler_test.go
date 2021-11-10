package json_handler

import (
	"strconv"
	"testing"
)

func Test_main(t *testing.T) {
	example := "{\"store\":{\"book\":[{\"category\":\"小说\",\"author\":\"鲁迅\",\"title\":\"呐喊\",\"price\":12.99},{\"category\":\"\\u5c0f\\u8bf4\",\"author\":\"\\u9c81\\u8fc5\",\"title\":\"\\u5450\\u558a\",\"price\":12.99},{\"category\":\"reference\",\"author\":\"NigelRees\",\"title\":\"SayingsoftheCentury\",\"price\":8.95}],\"bicycle\":{\"color\":\"red\",\"price\":19.95}},\"expensive\":10}"
	main(example)
}

func Test_special(t *testing.T) {
	specials := " {}[]:,\"'\\\n\r\n"
	for i, char := range specials {
		println(i, char)
	}
}

func Test_backSlash(t *testing.T) {
	specials := "\"鲁"
	for i, char := range specials {
		println(i, char)
	}
}

func Test_read_file(t *testing.T) {
	s := "\\u558a"
	//bytes1, _ := ioutil.ReadFile(s)
	//for _, b := range bytes1 {
	//	println(b)
	//}
	for i, char := range s {
		println(i, char)
	}
}

func Test_unicos(t *testing.T) {
	s := "\\u5450"
	var a []rune
	for _, char := range s {
		if char != '\\' && char != 'u' {
			a = append(a, char)
		}
	}
	val, _ := strconv.ParseUint(string(a), 16, 32)
	println(val)
}
