package export

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"sync/atomic"

	"github.com/skvoch/blender-byte-backend/internal/model"
	psqlstore "github.com/skvoch/blender-byte-backend/internal/store/psql"
)

// Exporter ...
type Exporter struct {
	jsonPath string

	types map[string]*model.Type
	store *psqlstore.PSQLStore
	mux   sync.Mutex
	index uint64
}

// New ...
func New(jsonPath string) (*Exporter, error) {
	store, err := psqlstore.NewTest()

	if err != nil {

		return nil, err
	}

	return &Exporter{
		jsonPath: jsonPath,
		types:    make(map[string]*model.Type),
		store:    store,
	}, nil
}

func (e *Exporter) handleBook(ch chan JSONBook, wg *sync.WaitGroup) {
	for {
		book := <-ch
		if book.Topic == "" {
			continue
		}
		_type := &model.Type{}

		e.mux.Lock()
		if _t, ok := e.types[book.Topic]; ok {
			_type = _t
		} else {
			createdType, err := e.store.AddType(&model.Type{
				Name: book.Topic,
			})

			if err != nil {
				fmt.Println("Create type ERROR", err)
				continue
			}
			e.types[book.Topic] = createdType
			_type = e.types[book.Topic]
		}
		e.mux.Unlock()

		storeBook, err := e.store.AddBook(&model.Book{
			Name:        book.Name,
			Author:      book.Author,
			Cost:        book.Cost,
			Date:        book.Date,
			Description: book.Description,
			FullName:    book.FullName,
			ISBN:        book.ISBN,
			Photo:       book.Photo,
			Publish:     book.Publish,
		})
		if err != nil {
			fmt.Println("Create book ERROR", err)
			continue
		}

		if err != e.store.AssingBookToType(storeBook, _type) {
			fmt.Println("Book asign ERROR", err)
			continue
		}

		fmt.Println(e.index, "created ")
		atomic.AddUint64(&e.index, 1)
		wg.Done()
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

	books := make([]JSONBook, 0)

	err = json.Unmarshal(byteValue, &books)
	if err != nil {
		log.Printf("error decoding sakura response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
	}
	wg := sync.WaitGroup{}

	ch := make(chan JSONBook, 128)

	for i := 0; i < 30; i++ {
		go e.handleBook(ch, &wg)
	}
	wg.Add(len(books))

	for _, book := range books {
		ch <- book
	}
	wg.Wait()
}
