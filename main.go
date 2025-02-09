package main

import (
	"fmt"
)

func main() {
	fmt.Println("Started")

	downloadViaHash("aeb4a909b930228bfd62bd04322b6cc7861ad154", "1.21.4", "fabric", "output/")
	downloadViaHash("aeb4a909b930228bfd62bd04322b6cc7861ad154", "1.21.2", "fabric", "output/")
}
