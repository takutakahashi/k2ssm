package output

import "fmt"

type Sops struct{}

func (o Sops) Execute() error {
	return fmt.Errorf("not implemented")
}
