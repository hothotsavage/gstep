package pagination

type Pagination[T any] struct {
	Total int `json:"total"`
	List  []T `json:"list"`
}
