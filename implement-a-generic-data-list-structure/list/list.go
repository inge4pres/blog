package primitive

type List struct {
	data []interface{}
}

type filterFn func(interface{}) interface{}
type foldFn func([]interface{}) interface{}

func (m *List) Reverse() *List {
	var ret []interface{}
	for i := 1; i <= len(m.data); i++ {
		ret = append(ret, m.data[len(m.data)-i])
	}
	return &List{
		data: ret,
	}
}

func (m *List) Map() []interface{} {
	return m.data
}

func (m *List) Filter(filter filterFn) *List {
	var ret []interface{}
	for d := range m.data {
		if data := filter(m.data[d]); data != nil {
			ret = append(ret, data)
		}
	}
	return &List{
		data: ret,
	}
}

func (m *List) FoldLeft(fold foldFn) *List {
	var ret []interface{}
	ret = append(ret, fold(m.data))
	return &List{
		data: ret,
	}
}
