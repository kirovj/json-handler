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

func (h *Handler) handle(s string) string {

	// invalid json string
	if len(s) < 2 {
		return ""
	}

	// start push single rune to chan
	go func() {
		for _, char := range s {
			h.ch <- char
		}
		h.ch <- End
	}()

	// pull rune from chan
	for flag := h.innerHandle(); flag != End; {
	}
	return h.builder.String()
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
		// todo pull until , " } ]
	case Quote:
		// pull from chan until close quote
		return h.handleValue()
	case Colon:
		h.append(Space)
		h.ignoreSpace = true
	case Comma:
		h.appendCurrent()
	}
	return h.current
}

func (h *Handler) handleValue() int32 {
	for r := <-h.ch; ; {
		if r == End {
			return End
		}

		h.current = r
		// only in value backslash is useful
		if r == Backslash {
			if next := <-h.ch; next == UnicodeFlag {
				// start parse unicode like \u....
				var unicos []rune
				for i := 0; i < 4; i++ {
					char := <-h.ch
					unicos = append(unicos, char)
					if char == End {
						break
					}
				}
				h.current = unicos[len(unicos)-1]
				if len(unicos) == 4 && h.current != End {
					if val, err := parseUnicode(string(unicos)); err != nil {
						h.current = val
						h.appendCurrent()
						return h.current
					}
				}
				// if parse unicode failed or len of unicos != 4
				h.append(Backslash).append(UnicodeFlag).extend(&unicos)
			} else {
				h.append(next)
			}
		} else {
			h.appendCurrent()
			if r == Quote {
				if h.last != Backslash {
					break
				}
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
