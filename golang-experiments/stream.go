package experiments

type Iterator interface {
	Close()
	Stream() *Iterator
	Collect(func(interface{}) bool) []interface{}
}

type Iterable struct {
	list     []interface{}
	ingress  chan (*iterand)
	filtered chan (*iterand)
	done     chan bool
	output   []interface{}
}

type iterand struct {
	obj interface{}
}

// Stream creates the data stream through a channel for Iterable
func (i Iterable) Stream() *Iterable {
	i.output = make([]interface{}, 0)
	i.ingress = make(chan (*iterand), len(i.list))
	i.filtered = make(chan (*iterand))
	for _, item := range i.list {
		i.ingress <- &iterand{obj: item}
	}
	return &i
}

func (i Iterable) Collect(filterFn func(interface{}) bool) *Iterable {
	// go func() {
	// 	for filtered := range i.filtered {
	// 		output = append(output, filtered.obj)
	// 	}
	// }()
	// for input := range i.ingress {
	// 	go func(in *iterand) {
	// 		if filterFn(in.obj) {
	// 			i.filtered <- in
	// 		}
	// 	}(input)
	// }

	select {
	case filtered := <-i.filtered:
		i.output = append(i.output, filtered.obj)
	case input := <-i.ingress:
		go func(in *iterand) {
			if filterFn(in.obj) {
				i.filtered <- in
			}
		}(input)
	case <-i.done:
		break
	}
	return &i
}

func (i Iterable) Close() {
	i.done <- true
	close(i.ingress)
	close(i.filtered)
}
