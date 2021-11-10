package json_handler

import (
	"strconv"
	"strings"
)

const (
	Space        = ' '
	Colon        = ':'
	Comma        = ','
	Quote        = '"'
	OpenBracket  = '['
	CloseBracket = ']'
	OpenCurly    = '{'
	CloseCurly   = '}'
	Backslash    = '\\'
	NewlineN     = '\n'
	NewlineR     = '\r'
	UnicodeFlag  = 'u'
	End          = -1 // end of json
)

func isWhiteSpace(r rune) bool {
	return r == Space || r == NewlineN || r == NewlineR
}

func isClose(r rune) bool {
	return r == CloseCurly || r == CloseBracket
}

type Handler struct {
	ch      chan rune
	level   int
	pos     int
	current rune
	last    rune
	result  string
	builder strings.Builder
}

func NewHandler() *Handler {
	return &Handler{
		ch:      make(chan rune),
		builder: strings.Builder{},
	}
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

func (h *Handler) next() rune {
	h.pos++
	return <-h.ch
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
	for flag := h.innerHandle(); flag != End; flag = h.innerHandle() {
		println(flag)
	}
	return h.builder.String()
}

func (h *Handler) innerHandle() rune {
	defer func() {
		h.last = h.current
	}()

	h.current = h.next()

	switch h.current {
	case OpenCurly | OpenBracket:
		h.level++
		h.append(h.current).append(NewlineN)
		for {
			r := h.next()
			if r == Quote || r == End {
				h.current = r
				break
			}
		}
	case CloseCurly | CloseBracket:
		h.level--
		// todo pull until , " } ]
	case Quote:
		// pull from chan until close quote
		return h.handleValue(true)
	case Colon:
		h.append(Space)
		for {
			r := h.next()
			if isWhiteSpace(r) {
				continue
			}
			h.current = r
			h.append(r)
			if !isClose(r) {
				h.handleValue(false)
			}
			break
		}
	case Comma:
		h.append(h.current)
		h.append(NewlineN)
	}
	return h.current
}

func (h *Handler) handleValue(inQuote bool) rune {
	for {
		r := h.next()
		if r == End {
			return End
		}

		h.current = r
		if inQuote {
			// only in quote backslash is useful
			if r == Backslash {
				if next := h.next(); next == UnicodeFlag {
					// start parse unicode like \u....
					var unicos []rune
					for i := 0; i < 4; i++ {
						char := h.next()
						unicos = append(unicos, char)
						if char == End {
							break
						}
					}
					h.current = unicos[len(unicos)-1]
					if len(unicos) == 4 && h.current != End {
						if val, err := parseUnicode(string(unicos)); err != nil {
							h.current = val
							h.append(val)
							return h.current
						}
					}
					// if parse unicode failed or len of unicos != 4
					h.append(Backslash).append(UnicodeFlag).extend(&unicos)
				} else {
					h.append(next)
				}
			} else {
				h.append(h.current)
				if r == Quote {
					if h.last != Backslash {
						break
					}
				}
			}
		} else {
			// value is num or bool which not in quote
			if r == Comma || r == CloseCurly || isWhiteSpace(r) {
				break
			}
			h.append(r)
		}
	}
	return h.current
}

func parseUnicode(hex string) (rune, error) {
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
