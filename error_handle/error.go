package error_handle

import "fmt"

type CustomError struct {
   Code int
   Msg string
}
func (e *CustomError) Error() string {
    return fmt.Sprintf("custom error %d: %s", e.Code, e.Msg)
}