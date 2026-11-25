package domain

import (
	"fmt"
	"strings"

	"github.com/subhadip0x539/arc-dex-core-srv/src/internal/dto"
)

type SortField string

const (
	SortFieldName      SortField = "name"
	SortFieldCreatedAt SortField = "created_at"
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

func NewSortField(field string) SortField {
	switch f := SortField(field); f {
	case SortFieldName:
		return SortFieldName
	default:
		return SortFieldCreatedAt
	}
}

func NewSortOrder(order string) SortOrder {
	switch o := SortOrder(order); o {
	case SortOrderASC, SortOrderDESC:
		return o
	default:
		return SortOrderDESC
	}
}

func NewSort(sort dto.Sort) Sort {
	return Sort{
		By:    NewSortField(sort.By),
		Order: NewSortOrder(sort.Order),
	}
}

type FilterField string

const (
	FilterFieldName      FilterField = "name"
	FilterFieldCreatedAt FilterField = "created_at"
	FilterFieldDex       FilterField = "dex_id"
)

type FilterOperator string

const (
	OperatorEq   FilterOperator = "eq"
	OperatorNeq  FilterOperator = "neq"
	OperatorGt   FilterOperator = "gt"
	OperatorGte  FilterOperator = "gte"
	OperatorLt   FilterOperator = "lt"
	OperatorLte  FilterOperator = "lte"
	OperatorLike FilterOperator = "like"
)

type Filter struct {
	Field    FilterField
	Operator FilterOperator
	Values   []any
}

func NewFilterField(field string) FilterField {
	switch FilterField(field) {
	case FilterFieldName,
		FilterFieldCreatedAt,
		FilterFieldDex:
		return FilterField(field)
	default:
		return FilterFieldName
	}
}

func NewFilterOperator(operator string) FilterOperator {
	switch FilterOperator(strings.ToLower(operator)) {
	case OperatorEq, OperatorNeq, OperatorGt, OperatorGte, OperatorLt, OperatorLte, OperatorLike:
		return FilterOperator(strings.ToLower(operator))
	default:
		return OperatorEq
	}
}

func NewFilter(filter dto.Filter) Filter {
	return Filter{
		Field:    NewFilterField(filter.Field),
		Operator: NewFilterOperator(filter.Operator),
		Values:   filter.Values,
	}
}

type Query struct {
	Sort    Sort
	Filters []Filter
}

func NewQuery(sort dto.Sort, filters []dto.Filter) *Query {
	st := NewSort(sort)
	fs := make([]Filter, 0, len(filters))

	for _, filter := range filters {
		fs = append(fs, NewFilter(filter))
	}

	return &Query{
		Sort:    st,
		Filters: fs,
	}
}

func NewCondition(alias string, sort Sort, filters []Filter) (string, []any, error) {
	conditions := make([]string, 0, len(filters))
	args := make([]any, 0)

	for _, filter := range filters {
		field := fmt.Sprintf("%s.%s", alias, filter.Field)

		if len(filter.Values) == 0 {
			return "", nil, fmt.Errorf("filter values cannot be empty")
		}

		var operator string

		switch filter.Operator {
		case OperatorEq:
			operator = "="
		case OperatorNeq:
			operator = "<>"
		case OperatorGt:
			operator = ">"
		case OperatorGte:
			operator = ">="
		case OperatorLt:
			operator = "<"
		case OperatorLte:
			operator = "<="
		case OperatorLike:
			operator = "ILIKE"
		default:
			return "", nil, fmt.Errorf("unsupported filter operator: %s", filter.Operator)
		}

		args = append(args, filter.Values[0])

		conditions = append(conditions,
			fmt.Sprintf("%s %s $%d", field, operator, len(args)),
		)

	}

	var builder strings.Builder

	if len(conditions) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(strings.Join(conditions, " AND "))
	}

	if sort.By != "" {
		fmt.Fprintf(&builder, " ORDER BY %s.%s %s", alias, sort.By, sort.Order)
	}

	return builder.String(), args, nil
}
