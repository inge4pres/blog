package filter

type check func(interface{}) bool

func Strings(in []string, fn check) []string {
	r := make([]string, 0)
	for _, s := range in {
		if fn(len(s)) {
			continue
		}
		r = append(r, s)
	}
	return r
}

func BiggerThan(b int) check {
	return func(i interface{}) bool {
		return i.(int) > b
	}
}
