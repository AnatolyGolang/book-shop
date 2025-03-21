package models

type CartUpdateRequest struct {
	BookIds []int `json:"book_ids"`
}
