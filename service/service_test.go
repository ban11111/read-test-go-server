package service

import (
	"fmt"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	create, err := os.Create("./sf/fdsf/sss.txt")
	fmt.Println(err)
	fmt.Println(create)
}
