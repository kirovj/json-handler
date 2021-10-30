package json_handler

const (
	SPACE         = 32  // 空格
	COLON         = 58  // 冒号
	COMMA         = 44  // 逗号
	DOUBLE_QUOTE  = 34  // 双引号
	SINGLE_QUOTE  = 39  // 单引号
	OPEN_BRACKET  = 91  // 左中括号
	CLOSE_BRACKET = 93  // 右中括号
	OPEN_CURLY    = 123 // 左花括号
	CLOSE_CURLY   = 125 // 右花括号
)

var example = "{\"store\":{\"book\":[{\"category\":\"小说\",\"author\":\"鲁迅\",\"title\":\"呐喊\",\"price\":12.99},{\"category\":\"\\u5c0f\\u8bf4\",\"author\":\"\\u9c81\\u8fc5\",\"title\":\"\\u5450\\u558a\",\"price\":12.99},{\"category\":\"reference\",\"author\":\"NigelRees\",\"title\":\"SayingsoftheCentury\",\"price\":8.95}],\"bicycle\":{\"color\":\"red\",\"price\":19.95}},\"expensive\":10}"

func handle(s string) {
	//for idx, char := range s {
	//	println(idx, char)
	//}
	for i := 0; i < len(s); i++ {
		println(i, s[i])
	}
}

func main() {
	handle(example)
}
