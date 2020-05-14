package models

type Post struct {
  ID       string `json:"id"`
  Title    string `json:"title"`
  Content  string `json:"content"`
  AuthorID string `json:"author_id"`
}