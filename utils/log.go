package utils

import "fmt"

func Log(err error) {
	if err != nil {
		fmt.Println("Error is -> " + err.Error())
	}
}
