package humans

import (
	"fmt"
	"strings"
)

func Simple(ID int64, name string, age int, registered bool) string {
	nameWithNoSpaces := strings.Replace(name, " ", "-", -1)
	return fmt.Sprintf("%d-%s-%d-%t", ID, nameWithNoSpaces, age, registered)
}
