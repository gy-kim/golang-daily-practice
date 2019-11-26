package gin

import (
	"fmt"
	"strings"
)

const ginSupportMonGoVer = 10

func IsDebugging() bool {
	return ginMode == debugCode
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(DefaultWriter, "[GIN-debug] "+format, values...)
	}
}
