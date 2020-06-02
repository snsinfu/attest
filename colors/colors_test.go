package colors

import (
	"strings"
	"testing"
)

func TestColors_reset(t *testing.T) {
	const resetCode = "\x1b[m"

	funcs := map[string](func(string) string){
		"Black":   Black,
		"Red":     Red,
		"Green":   Green,
		"Yellow":  Yellow,
		"Blue":    Blue,
		"Magenta": Magenta,
		"Cyan":    Cyan,
		"White":   White,
		"Gray":    Gray,
	}

	for name, fun := range funcs {
		out := fun("")

		if !strings.HasSuffix(out, resetCode) {
			t.Errorf("%v does not reset color: %x", name, []byte(out))
		}
	}
}
