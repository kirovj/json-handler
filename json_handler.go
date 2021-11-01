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

type Handler struct {
	//idx     int
	ch          chan int32
	current     int32
	last        int32
	result      string
	unicos      []int32
	ignoreSpace bool
	insideQuote bool
	escape      bool
	builder     strings.Builder
}

func NewHandler() *Handler {
	return &Handler{
		ch:          make(chan int32),
		builder:     strings.Builder{},
		ignoreSpace: true,
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
	go func() {
		for _, char := range s {
			h.ch <- char
		}
		h.ch <- -1
	}()

	for {
		if flag := h.innerHandle(); flag == -1 {
			break
		}
	}
}

func (h *Handler) innerHandle() int32 {
	defer func() {
		h.last = h.current
	}()

	h.current = <-h.ch
	if h.escape {
		if h.current == 'u' {
			// todo
		} else {
			h.appendLast()
			h.appendCurrent()
			return h.current
		}
	}

	switch h.last {
	case OpenCurly | OpenBracket:
		h.append(NewlineN)
	case CloseCurly | CloseBracket:
	case DoubleQuote | SingleQuote:
		if !h.insideQuote && !h.escape {
			h.insideQuote = true
			h.ignoreSpace = false
			h.appendCurrent()
			break
		}
		if h.insideQuote {

		}
	case Colon:
		h.append(Space)
		h.ignoreSpace = true
	case Space:
		if !h.ignoreSpace {
			h.appendCurrent()
		}
	case NewlineN | NewlineR:
	case Backslash:
		h.escape = true
	}
	return h.current
}

func main(s string) {
	h := NewHandler()
	h.handle(s)
}
