package jsonmap

// IsNil checks if value is null or undefined.
func IsNil(json *Json) bool {
	return json == nil || json.IsNil()
}

// Filter iterates over elements of collection, returning an array of all elements predicate returns truthy for.
func Filter(collection []*Json, predicate func(v *Json) bool) []*Json {
	var values []*Json

	if predicate == nil {
		return values
	}

	for _, j := range collection {
		if predicate(j) {
			copy := *j // filter must returns new array
			values = append(values, &copy)
		}
	}

	return values
}

// Map creates an array of values by running each element in collection thru iteratee
func Map(collection []*Json, iteratee func(v *Json) *Json) []*Json {
	var values []*Json

	if iteratee == nil {
		return values
	}

	for _, j := range collection {
		copy := *j // map must returns new array
		values = append(values, iteratee(&copy))
	}

	return values
}

// Assign Assigns own enumerable string keyed properties of source objects to the destination object.
func Assign(source *Json, json ...*Json) *Json {
	if IsNil(source) {
		source = New()
	}

	for _, j := range json {
		if !IsNil(j) {
			if o := j.AsObject(); o != nil {
				for k, v := range o {
					source.Set(k, v)
				}
			}
		}
	}

	return source
}
