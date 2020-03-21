package main

import (
	"fmt"

	export "github.com/skvoch/blender-byte-backend/internal/export"
)

func main() {
	exporter, err := export.New("/Volumes/Macintosh HD/Projects/blender-byte-backend/books.json")

	if err != nil {
		fmt.Print("Error in main:", err)
	}
	exporter.Start()
}
