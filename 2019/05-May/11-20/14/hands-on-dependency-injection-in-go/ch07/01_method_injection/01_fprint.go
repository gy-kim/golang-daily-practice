package method_injection

import (
	"fmt"
	"os"
)

func ExampleA() {
	fmt.Fprintf(os.Stdout, "Hello world")
}
