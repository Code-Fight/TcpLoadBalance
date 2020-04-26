package units

import (
	"fmt"
	"testing"
)

func TestGetBackend(t *testing.T) {
	ConfigInit()
	a:=GetBackend()
	fmt.Print(len(a))
}
