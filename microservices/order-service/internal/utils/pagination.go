package utils

type PaginationQuery struct{ Page, PerPage int; SortBy, Order string }
func (q PaginationQuery) GetLimit() int { if q.PerPage <= 0 { return 20 }; if q.PerPage > 100 { return 100 }; return q.PerPage }
func (q PaginationQuery) GetOffset() int { p := q.Page; if p <= 1 { p = 1 }; return (p-1) * q.GetLimit() }
func BuildMeta(total, page, perPage int, sortBy, order string) Meta { if perPage <= 0 { perPage = 20 }; if page <= 0 { page = 1 }; return Meta{ Total: total, Page: page, PerPage: perPage, SortBy: sortBy, Order: order } }