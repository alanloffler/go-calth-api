package response

type PaginatedData[T any] struct {
	Result []T   `json:"result"`
	Total  int32 `json:"total"`
}
