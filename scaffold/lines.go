package scaffold

import (
	"io"
	"os"
	"strings"
)

func Lines() []string {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	return strings.Split(strings.TrimSpace(string(b)), "\n")
}
