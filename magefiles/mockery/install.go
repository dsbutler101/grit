package mockery

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

const (
	Version = "v2.43.0"
)

func Install() error {
	fmt.Println("Installing mockery")
	return sh.Run("go", "install", fmt.Sprintf("github.com/vektra/mockery/v2@%s", Version))
}
