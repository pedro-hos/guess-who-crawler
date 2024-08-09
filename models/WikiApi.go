package models

type Response struct {
	Batchcomplete bool  `json:"batchcomplete"`
	Query         Query `json:"query"`
}

type Query struct {
	Normalized []Normalized `json:"normalized"`
	Pages      []Page       `json:"pages"`
}

type Normalized struct {
	Fromencoded bool   `json:"fromencoded"`
	From        string `json:"from"`
	To          string `json:"to"`
}

type Page struct {
	Pageid    int64      `json:"pageid"`
	NS        int64      `json:"ns"`
	Title     string     `json:"title"`
	Revisions []Revision `json:"revisions"`
}

type Revision struct {
	Slots Slots `json:"slots"`
}

type Slots struct {
	Main Main `json:"main"`
}

type Main struct {
	Contentmodel  string `json:"contentmodel"`
	Contentformat string `json:"contentformat"`
	Content       string `json:"content"`
}
