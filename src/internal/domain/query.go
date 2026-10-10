package domain

type SortField string

const (
	SortFieldName      SortField = "name"
	SortFieldHeight    SortField = "height"
	SortFieldAge       SortField = "age"
	SortFieldBirthdate SortField = "birthdate"
)

type SortOrder string

const (
	SortOrderASC  SortOrder = "ASC"
	SortOrderDESC SortOrder = "DESC"
)

type Sort struct {
	By    SortField
	Order SortOrder
}

type Query struct {
	Sort Sort
}

func NewQuery(by string) Query {
	b := SortFieldName
	o := SortOrderASC

	switch by {
	case string(SortFieldHeight):
		b = SortFieldHeight
	case string(SortFieldAge):
		b = SortFieldBirthdate
		o = SortOrderDESC

	}

	return Query{
		Sort: Sort{
			By:    b,
			Order: o,
		},
	}
}
