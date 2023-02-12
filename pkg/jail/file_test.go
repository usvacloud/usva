package jail

import (
	"context"
	"os"
	"testing"
)

func TestJail(t *testing.T) {
	var (
		values   = []string{"isovauvajeesus", "leego", "luuranko", "kolmannen-silmänkyyneleet"}
		testfile = os.TempDir() + "/usva-jail-test"
	)

	for _, value := range values {
		file := getDescriptor(t, testfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE)
		defer file.Close()

		var jail Jail = NewJailFS(file)
		err := jail.Ban(context.Background(), value)
		checkError(t, err)
	}

	for _, value := range values {
		file := getDescriptor(t, testfile, os.O_RDONLY)
		defer file.Close()

		var jail Jail = NewJailFS(file)
		found, _ := jail.IsAuthorized(context.Background(), value)
		if !found {
			t.Fatalf("test failed: %s not found", value)
		}
	}
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("test failed: %s", err)
	}
}

func getDescriptor(t *testing.T, path string, flag int) *os.File {
	fp, err := os.OpenFile(path, flag, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	return fp
}
