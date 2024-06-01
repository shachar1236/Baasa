package api

type QueryMessage struct {
    QueryId int64 `json:"QueryId"`
    Session string `json:"Session"`
    QueryArgs map[string]any `json:"QuaryArgs"`
}
