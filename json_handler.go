package json_handler

import (
	"strconv"
	"strings"
)

const (
	Space        = 32  // space
	Colon        = 58  // :
	Comma        = 44  // ,
	Quote        = 34  // "
	OpenBracket  = 91  // [
	CloseBracket = 93  // ]
	OpenCurly    = 123 // {
	CloseCurly   = 125 // }
	Backslash    = 92  // \
	NewlineN     = 10  // \n
	NewlineR     = 13  // \r
	UnicodeFlag  = 'u' // u
	End          = -1  // end of json
)

type Handler struct {
	ch          chan int32
	current     int32
	last        int32
	result      string
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

func (h *Handler) appendCurrent() *Handler {
	h.builder.WriteRune(h.current)
	return h
}

func (h *Handler) appendLast() *Handler {
	h.builder.WriteRune(h.last)
	return h
}

func (h *Handler) append(r rune) *Handler {
	h.builder.WriteRune(r)
	return h
}

func (h *Handler) extend(runes *[]rune) *Handler {
	for _, r := range *runes {
		h.append(r)
	}
	return h
}

func (h *Handler) handle(s string) {
	go func() {
		for _, char := range s {
			h.ch <- char
		}
		h.ch <- End
	}()

	for {
		if flag := h.innerHandle(); flag == End {
			break
		}
	}
}

func (h *Handler) innerHandle() int32 {
	defer func() {
		h.last = h.current
	}()

	h.current = <-h.ch

	switch h.last {
	case OpenCurly | OpenBracket:
		h.append(NewlineN)
	case CloseCurly | CloseBracket:
	case Quote:
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
		if h.insideQuote {
			// only in quotes backslash is useful
			if h.current == UnicodeFlag {
				// start parse unicode like \u....
				var unicos []rune
				for i := 0; i < 4; i++ {
					char := <-h.ch
					unicos = append(unicos, char)
					if char == End {
						break
					}
				}
				if len(unicos) == 4 {
					if val, err := parseUnicode(string(unicos)); err != nil {
						h.current = val
						h.appendCurrent()
						return h.current
					}
				}
				// if parse unicode failed or len of unicos != 4
				h.append(Backslash).append(UnicodeFlag).extend(&unicos)
				h.current = unicos[len(unicos)-1]
			}
		}
	}
	return h.current
}

func parseUnicode(hex string) (int32, error) {
	val, err := strconv.ParseUint(hex, 16, 32)
	if err == nil {
		return rune(val), err
	}
	return 0, err
}

func main(s string) {
	h := NewHandler()
	h.handle(s)
}
