package mathmod

import "fmt"

func errNArgs(want, got int) error {
	return fmt.Errorf("[moonBASIC] Runtime Error: expects %d argument(s), got %d", want, got)
}

func errNArgsRange(want string, got int) error {
	return fmt.Errorf("[moonBASIC] Runtime Error: expects %s argument(s), got %d", want, got)
}
