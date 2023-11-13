package db

import (
	"encoding/json"
	"strconv"
)

type CompareOP int

const (
	Less CompareOP = iota
	LessEqual
	Greater
	GreaterEqual
	Equal
	NotEqual
)

type FilterOperation struct {
	Key   string    `json:"key"`
	Op    CompareOP `json:"op"`
	Value any       `json:"value"`
}

func FilterOperationsToSql(base int, fs []FilterOperation) (string, []any) {
	s := ""
	values := make([]any, len(fs))

	if len(fs) > 0 {
		s += "WHERE "
	} else {
		return "", nil
	}

	for i, e := range fs {
		if i > 0 && i != len(fs)-1 {
			s += " AND "
		}

		s += e.Key + e.Op.String() + "$" + strconv.Itoa(base+i+1)
		values[i] = e.Value
	}
	return s, values
}

func (op CompareOP) String() string {
	switch op {
	default:
		return "="
	case Less:
		return "<"
	case LessEqual:
		return "<="
	case Greater:
		return ">"
	case GreaterEqual:
		return ">="
	case Equal:
		return "="
	case NotEqual:
		return "<>"
	}
}

func (op *CompareOP) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	default:
		*op = Equal
	case "<":
		*op = Less
	case "<=":
		*op = LessEqual
	case ">":
		*op = Greater
	case ">=":
		*op = GreaterEqual
	case "=":
		*op = Equal
	case "<>":
		*op = NotEqual
	}

	return nil
}

func (op CompareOP) MarshalJSON() ([]byte, error) {
	var s string
	switch op {
	default:
		s = "="
	case Less:
		s = "<"
	case LessEqual:
		s = "<="
	case Greater:
		s = ">"
	case GreaterEqual:
		s = ">="
	case Equal:
		s = "="
	case NotEqual:
		s = "<>"
	}

	return json.Marshal(s)
}
