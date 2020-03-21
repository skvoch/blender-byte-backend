package export

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/skvoch/blender-byte-backend/internal/model"
)

// Exporter ...
type Exporter struct {
	jsonPath string
}

// New ...
func New(jsonPath string) *Exporter {
	return &Exporter{
		jsonPath: jsonPath,
	}
}

// Start ...
func (e *Exporter) Start() {
	jsonFile, err := os.Open(e.jsonPath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened books.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	books := make([]*model.Book, 0)

	err = json.Unmarshal(byteValue, &books)
	if err != nil {
		log.Printf("error decoding sakura response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
	}
}
