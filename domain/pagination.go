package domain

type PaginateResult[V interface{}] struct {
	Items        *[]V  `json:"items"`
	TotalItems   int64 `json:"totalItems"`
	ItemsPerPage int64 `json:"itemsPerPage"`
	Page         int64 `json:"page"`
	ItemCount    int64 `json:"itemCount"`
}

type PaginationOptions struct {
	Page      int64
	Limit     int64
	Sort      map[string]string
	Filter    interface{}
	Collation interface{}
}
type PaginationService[V interface{}] interface {
	Paginate(paginateOption PaginationOptions) PaginateResult[V]
}
