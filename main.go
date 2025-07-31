package main

import (
	"fmt"

	"github.com/daozhonglee/go-util/errorutil"
	"github.com/daozhonglee/go-util/random"
)

func main() {
	defer errorutil.Recover()
	fmt.Println(random.Int(1, 100))
	fmt.Println(random.Int(1, 100))
	fmt.Println(random.Int(1, 100))
	fmt.Println(random.Int(1, 100))
}
