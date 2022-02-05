package main

import (
	"errors"
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(`LinEEEE`, Line())
	fmt.Println(`LinEEEE 2222`, Line())
	fmt.Println(`ERRRR : `, test().Error())
}

func test() error {

	return errors.New(Line())
}

func Line() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf(`%s:%d`, file, line)
}
