package service

import "fmt"

func logError(value error) error {
	fmt.Printf("Got error %s", value.Error())
	return value
}
