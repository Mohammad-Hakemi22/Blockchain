package utility

import "log"

func ErrorHandler(msg string, err error) {
	if err != nil {
		log.Fatalln(msg)
	}
}