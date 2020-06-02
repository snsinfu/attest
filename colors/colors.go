// Package colors contain helper functions to decorate text with ANSI escape
// sequences.
package colors

func Black(s string) string {
	return "\x1b[30m" + s + "\x1b[m"
}

func Red(s string) string {
	return "\x1b[31m" + s + "\x1b[m"
}

func Green(s string) string {
	return "\x1b[32m" + s + "\x1b[m"
}

func Yellow(s string) string {
	return "\x1b[33m" + s + "\x1b[m"
}

func Blue(s string) string {
	return "\x1b[34m" + s + "\x1b[m"
}

func Magenta(s string) string {
	return "\x1b[35m" + s + "\x1b[m"
}

func Cyan(s string) string {
	return "\x1b[36m" + s + "\x1b[m"
}

func White(s string) string {
	return "\x1b[37m" + s + "\x1b[m"
}

func Gray(s string) string {
	return "\x1b[90m" + s + "\x1b[m"
}
