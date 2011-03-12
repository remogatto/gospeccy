package spectrum

// Element of a linked-list of events
type Event interface {
	// Returns the T-state when this event occurred
	GetTState() uint

	// Returns an 'Event' of the same type as this 'Event',
	// or nil if this is the last element in the linked-list.
	//
	// Constraint: (T-state of the previous event) < (T-state of this event)
	GetPrevious_orNil() Event
}

type EventCondition interface {
	isTrue(e Event) bool
}

// Returns the number of events in the linked-list.
//
// If 'cond' is not nil, it will receive the individual list elements,
// and this function returns the number of elements in the linked-list
// up until an element E on which 'cond' returns 'false'. The count
// does not include E.
func EventListLength(head Event, cond EventCondition) int {
	n := 0

	if cond != nil {
		for e := head; (e != nil) && cond.isTrue(e); e = e.GetPrevious_orNil() {
			n++
		}
	} else {
		for e := head; e != nil; e = e.GetPrevious_orNil() {
			n++
		}
	}

	return n
}

type EventArray interface {
	Init(n int)
	Set(i int, e Event)
}

// Copies the events from the linked-list into the array.
// The events are sorted by T-state values, in ascending order.
//
// If 'cond' is not nil, it will receive the individual list elements,
// and the array will contain all elements up until an element E on which 'cond'
// returns 'false. The array does not contain E.
func EventListToArray_Ascending(head Event, array EventArray, cond EventCondition) {
	var n int = EventListLength(head, cond)

	array.Init(n)

	for e, i := head, (n - 1); i >= 0; {
		array.Set(i, e)

		i = i - 1
		e = e.GetPrevious_orNil()
	}

}
