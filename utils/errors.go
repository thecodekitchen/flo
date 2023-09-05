package utils

import (
	"fmt"
	"log"
	"os"
)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func ExistsError(err error) {
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}

func PrintError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
