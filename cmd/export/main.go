package main

import (
	export "github.com/skvoch/blender-byte-backend/internal/export"
)

func main() {
	exporter := export.New("/Volumes/Macintosh HD/Projects/blender-byte-backend/books.json")
	exporter.Start()
}
