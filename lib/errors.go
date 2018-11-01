package lib

import "fmt"

type SimpleError string
func (err SimpleError) Error() string {
	return string(err)
}

func SimpleErrorf(format string, args ...interface{}) SimpleError {
	return SimpleError(fmt.Sprintf(format, args...))
}
