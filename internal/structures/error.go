package structures

import "fmt"

type CustomError struct {
	Text     string `default:""`
	HttpCode int    `default:"200"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%d - %s", e.HttpCode, e.Text)
}
