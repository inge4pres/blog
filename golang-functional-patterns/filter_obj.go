package filter

type Obj struct {
	Strings  []string
	criteria *criteria
}

type criteria struct {
	filter func(interface{}) bool
}

func (c *criteria) WithLengthBiggerThan(number int) *criteria {
	c.filter = func(in interface{}) bool {
		n := in.(int)
		return n > number
	}
	return c
}

func (o *Obj) Filter() {
	for i, s := range o.Strings {
		if o.criteria.filter(len(s)) {
			o.Strings = append(o.Strings[:i],o.Strings[i+1:]...)
		}
	}
}

func (o *Obj) CriteriaFactory() *criteria{
	c := &criteria{}
	o.criteria = c
	return c

}