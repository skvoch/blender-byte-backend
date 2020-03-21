package export

// JSONBook ...
type JSONBook struct {
	Author string `json:"Author"`
	//Code        int64   `json:"Code"`
	Cost        float64 `json:"Cost"`
	Date        string  `json:"Date"`
	Description string  `json:"Description"`
	FullName    string  `json:"FullName"`
	ISBN        string  `json:"ISBN"`
	Name        string  `json:"Name"`
	Photo       string  `json:"Photo"`
	Publish     string  `json:"Publish"`
	//Series      int     `json:"Series"`
	Sheets int    `json:"Sheets"`
	Topic  string `json:"Topic"`
}
