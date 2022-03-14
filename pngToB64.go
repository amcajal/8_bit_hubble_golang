package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
)

func pngToBase64(inputFile string) string {
	// Read the content of the file as pure bytes
	imageData, _ := os.ReadFile(inputFile)

	// Simply turn the byte slice? array? into the string with the encoded data
	return base64.StdEncoding.EncodeToString(imageData)
}

func main() {
	f := os.Args[1]
	b64rep := pngToBase64(f)
	fmt.Printf("//%s\n%s\n", path.Base(f), b64rep)
}
