package scaffold

import (
	"io/ioutil"
	"os"
	"strings"
)

func Lines() []string {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	return strings.Split(strings.TrimSpace(string(b)), "\n")
}
