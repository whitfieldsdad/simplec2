package main

import (
	"fmt"

	"github.com/whitfieldsdad/simplec2/internal/util"
)

func main() {
	shell, _ := util.GetDefaultShell()
	fmt.Println(shell)
}
