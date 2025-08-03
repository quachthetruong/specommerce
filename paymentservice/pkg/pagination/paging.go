package pagination

import "math"

type Direction string

const (
	DirectionAsc      Direction = "ASC"
	DirectionDesc     Direction = "DESC"
	DefaultPageSize   uint      = 20
	MaximumPageSize   uint      = 200
	DefaultPageNumber uint      = 1
)

type Order struct {
	Direction  Direction
	ColumnName string
}

type Orders []Order

func (oo *Orders) Contain(columnName string) bool {
	for _, o := range *oo {
		if o.ColumnName == columnName {
			return true
		}
	}
	return false
}

// Add merge given orders with skipping duplicated items
func (oo *Orders) Add(orders ...Order) {
	for i := range orders {
		order := &orders[i]
		if !oo.Contain(order.ColumnName) {
			*oo = append(*oo, *order)
		}
	}
}

func (oo *Orders) Strings() []string {
	res := make([]string, 0, len(*oo))
	for _, o := range *oo {
		res = append(res, o.ColumnName+" "+string(o.Direction))
	}
	return res
}

// Paging request
type Paging struct {
	Sort   Orders
	Size   uint
	Number uint
}

func (p *Paging) Orders() Orders {
	return p.Sort
}

func (p *Paging) Limit() int {
	return int(p.Size)
}

func (p *Paging) Offset() int {
	return int((p.Number - 1) * p.Size)
}

func (p *Paging) TotalPages(totalRecords int) uint {
	if p.Size == 0 {
		return 1
	}
	return uint(math.Ceil(float64(totalRecords) / float64(p.Size)))
}

type MetaData struct {
	Total      int
	PageSize   uint
	PageNumber uint
	TotalPages uint
}

type Page[T any] struct {
	Data     []T
	Metadata MetaData
}