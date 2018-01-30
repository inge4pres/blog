package experiments

// SliceFunctor accepts a func and passes an argument slice to it
func SliceFunctor(functionToRun func([]interface{}) interface{}, args []interface{}) interface{} {
	return functionToRun(args)
}

func lastItem(args []interface{}) interface{} {
	return args[(len(args) - 1)]
}
