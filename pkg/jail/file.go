package jail

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

type Filejail struct {
	file *os.File
}

func NewJailFS(file *os.File) *Filejail {
	return &Filejail{
		file: file,
	}
}

func (j Filejail) Ban(ctx context.Context, value string) error {
	_, err := j.file.WriteString(fmt.Sprintf("%s\n", value))
	return err
}

func (j Filejail) IsAuthorized(ctx context.Context, value string) (bool, error) {
	scanner := bufio.NewScanner(j.file)
	for scanner.Scan() {
		if scanner.Text() == value {
			return true, nil
		}
	}

	return false, nil
}
