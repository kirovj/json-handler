package json_handler

import (
	"strings"
)

const (
	Space        = 32  // 空格
	Colon        = 58  // 冒号
	Comma        = 44  // 逗号
	DoubleQuote  = 34  // 双引号
	SingleQuote  = 39  // 单引号
	OpenBracket  = 91  // 左中括号
	CloseBracket = 93  // 右中括号
	OpenCurly    = 123 // 左花括号
	CloseCurly   = 125 // 右花括号
	Backslash    = 92  // 反斜杠
	NewlineN     = 10  // \n
	NewlineR     = 13  // \r\n
)

//func init() {
//	switch runtime.GOOS {
//	case "windows":
//		Bash = BashWin
//	case "linux":
//		Bash = BashLinux
//	}
//}

type Handler struct {
	//idx     int
	ch          chan int32
	current     int32
	last        int32
	result      string
	unicos      []int32
	ignoreSpace bool
	insideQuote bool
	builder     strings.Builder
}

func NewHandler() *Handler {
	return &Handler{
		ch:      make(chan int32),
		builder: strings.Builder{},
	}
}

func (h *Handler) appendCurrent() {
	h.builder.WriteRune(h.current)
}

func (h *Handler) appendLast() {
	h.builder.WriteRune(h.last)
}

func (h *Handler) append(r rune) {
	h.builder.WriteRune(r)
}

func (h *Handler) handle(s string) {
	// go func() {
	// 	for _, char := range s {
	// 		h.push(char)
	// 	}
	// 	h.ch <- -1
	// }

}

func (h *Handler) handleSingle() {
	switch h.last {
	case OpenCurly | OpenBracket:
		h.append(NewlineN)
		h.ignoreSpace = true
	case CloseCurly | CloseBracket:
		h.ignoreSpace = true
	case DoubleQuote | SingleQuote:
		h.ignoreSpace = false
	case Colon:
		h.append(Space)
		h.ignoreSpace = true
	case Space:
		if !h.ignoreSpace {
			h.appendCurrent()
		}
	case NewlineN | NewlineR:
	default:
		h.appendCurrent()
	}
}

func main(s string) {
	h := NewHandler()
	h.handle(s)
}
