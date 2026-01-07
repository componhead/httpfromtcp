package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if len(data) >= 2 && data[0] == '\r' && data[1] == '\n' {
		return 0, true, nil
	}
	crlfIndex := -1
	for i := 0; i < len(data)-1; i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			crlfIndex = i
			break
		}
	}
	if crlfIndex == -1 {
		return 0, false, nil
	}
	line := data[0:crlfIndex]
	s := strings.SplitN(string(line), ":", 2)
	keyToTrim := s[0]
	if len(s) != 2 || len(keyToTrim) == 0 || keyToTrim[len(keyToTrim)-1] == ' ' {
		return 0, false, fmt.Errorf("error: malformed key header")
	}
	key := strings.Trim(keyToTrim, " ")
	if key == "" {
		return 0, false, fmt.Errorf("error: malformed key header")
	} else {
		value := strings.Trim(s[1], " ")
		h[key] = value
		return len(line) + 2, false, nil
	}
}
