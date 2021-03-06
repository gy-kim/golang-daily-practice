package humans

import (
	"bytes"
	"strconv"
	"strings"
)

func NoSoSimple(ID int64, name string, age int, registered bool) string {
	out := &bytes.Buffer{}
	out.WriteString(strconv.FormatInt(ID, 10))
	out.WriteString("-")
	out.WriteString(strings.Replace(name, " ", "-", -1))
	out.WriteString(strconv.Itoa(age))
	out.WriteString("-")
	out.WriteString(strconv.FormatBool(registered))
	return out.String()
}
