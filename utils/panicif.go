package utils

import (
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
