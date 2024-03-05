package color

import "fmt"

var (
	reset   = "\033[0m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	gray    = "\033[37m"
	white   = "\033[97m"
)

// Red adds red color to str in terminal.
func Red(str string) string {
	return fmt.Sprintf("%s%s%s", red, str, reset)
}

// Green adds green color to str in terminal.
func Green(str string) string {
	return fmt.Sprintf("%s%s%s", green, str, reset)
}

// Yellow adds yellow color to str in terminal.
func Yellow(str string) string {
	return fmt.Sprintf("%s%s%s", yellow, str, reset)
}

// Blue adds blue color to str in terminal.
func Blue(str string) string {
	return fmt.Sprintf("%s%s%s", blue, str, reset)
}

// Magenta adds magenta color to str in terminal.
func Magenta(str string) string {
	return fmt.Sprintf("%s%s%s", magenta, str, reset)
}

// Cyan adds cyan color to str in terminal.
func Cyan(str string) string {
	return fmt.Sprintf("%s%s%s", cyan, str, reset)
}

// Gray adds gray color to str in terminal.
func Gray(str string) string {
	return fmt.Sprintf("%s%s%s", gray, str, reset)
}

// White adds white color to str in terminal.
func White(str string) string {
	return fmt.Sprintf("%s%s%s", white, str, reset)
}
