package env

import "reflect"

// Publish an object for other agents to find.
//
// In case the environment already contains an object of the same type,
// this function returns an error.
func Publish(object interface{}) (PublishedObject, error) {
	if object == nil {
		panic("nil object")
	}

	resultCh := make(chan cmd_publish_result)
	commandChannel <- cmd_publish{
		name_orEmpty: "",
		object:       object,
		resultCh:     resultCh,
	}
	result := <-resultCh
	return result.pub, result.err
}

type PublishedObject interface {
	// Removes a published object from the environment.
	// This method can be called at most once per PublishedObject.
	Remove()
}

// Publish an object which other agents can find by its name.
// The name cannot be empty, and the object cannot be nil.
//
// In case the environment already contains an object with the same name,
// this function returns an error.
func PublishName(name string, object interface{}) (PublishedObject, error) {
	if name == "" {
		panic("empty name")
	}
	if object == nil {
		panic("nil object")
	}

	resultCh := make(chan cmd_publish_result)
	commandChannel <- cmd_publish{
		name_orEmpty: name,
		object:       object,
		resultCh:     resultCh,
	}
	result := <-resultCh
	return result.pub, result.err
}

// Find a published object of the specified type.
//
// If there is currently no object of such a type in the environment,
// this function returns nil.
func Find(objectType reflect.Type) (object_orNil interface{}) {
	if objectType == nil {
		panic("nil objectType")
	}

	resultCh := make(chan cmd_find_result)
	commandChannel <- cmd_find{
		objectType: objectType,
		resultCh:   resultCh,
	}
	result := <-resultCh
	return result.object_orNil
}

// Find a published object with the specified name.
//
// If there is currently no object with such a name in the environment,
// this function returns nil.
func FindName(name string) (object_orNil interface{}) {
	if name == "" {
		panic("empty name")
	}

	resultCh := make(chan cmd_find_result)
	commandChannel <- cmd_findName{
		name:     name,
		resultCh: resultCh,
	}
	result := <-resultCh
	return result.object_orNil
}

// Wait until an object of the specified type appears in the environment.
//
// If the environment already contains the object,
// this function returns immediately.
//
// This function never returns nil.
func Wait(objectType reflect.Type) (object interface{}) {
	if objectType == nil {
		panic("nil objectType")
	}

	resultCh := make(chan cmd_wait_result)
	commandChannel <- cmd_wait{
		objectType: objectType,
		resultCh:   resultCh,
	}
	result := <-resultCh
	return result.object
}

// Wait until an object with the specified name appears in the environment.
//
// If the environment already contains the object,
// this function returns immediately.
//
// This function never returns nil.
func WaitName(name string) (object interface{}) {
	if name == "" {
		panic("empty name")
	}

	resultCh := make(chan cmd_wait_result)
	commandChannel <- cmd_waitName{
		name:     name,
		resultCh: resultCh,
	}
	result := <-resultCh
	return result.object
}

// Asynchronously wait for an object of the specified type to appear in the environment.
// When the object gets published, it will be sent to the specified channel.
//
// This function return immediately.
func WaitAsync(objectType reflect.Type, ch chan<- interface{}) {
	if objectType == nil {
		panic("nil objectType")
	}

	go func() {
		ch <- Wait(objectType)
	}()
}

// Asynchronously wait for an object with the specified name to appear in the environment.
// When the object gets published, it will be sent to the specified channel.
//
// This function return immediately.
func WaitNameAsync(name string, ch chan<- interface{}) {
	if name == "" {
		panic("empty name")
	}

	go func() {
		ch <- WaitName(name)
	}()
}
