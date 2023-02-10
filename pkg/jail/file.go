package jail

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

type filejail struct {
	file *os.File
}

func NewJailFS(file *os.File) filejail {
	return filejail{
		file: file,
	}
}

func (j filejail) Ban(ctx context.Context, value string) error {
	_, err := j.file.WriteString(fmt.Sprintf("%s\n", value))
	return err
}

func (j filejail) IsAuthorized(ctx context.Context, value string) (bool, error) {
	scanner := bufio.NewScanner(j.file)
	for scanner.Scan() {
		if scanner.Text() == value {
			return true, nil
		}
	}

	return false, nil
}

func (j filejail) GetFileDescriptor() *os.File {
	return j.file
}
