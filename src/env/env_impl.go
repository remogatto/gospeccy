package env

import (
	"errors"
	"reflect"
)

type objectInfo_t struct {
	name_orEmpty string

	objectType reflect.Type

	// An object of type 'objectType', or nil
	object_orNil interface{}

	// List of waiters waiting for 'object_orNil' to become non-nil.
	// If (object_orNil != nil) then (len(waiters) == 0).
	waiters []chan<- cmd_wait_result
}

func (i *objectInfo_t) Remove() {
	resultCh := make(chan cmd_remove_result)
	commandChannel <- cmd_remove{
		name_orEmpty: i.name_orEmpty,
		objectType:   i.objectType,
		resultCh:     resultCh,
	}
	result := <-resultCh
	if result.err != nil {
		panic(result.err.Error())
	}
}

// Set of all published objects
var objects = make(map[reflect.Type]*objectInfo_t)

// Set of all published named objects
var namedObjects = make(map[string]*objectInfo_t)

var commandChannel = make(chan interface{})

// Publish
type cmd_publish struct {
	name_orEmpty string
	object       interface{}
	resultCh     chan<- cmd_publish_result
}

type cmd_publish_result struct {
	pub PublishedObject
	err error
}

func do_publish(cmd cmd_publish) {
	objectType := reflect.TypeOf(cmd.object)

	var info *objectInfo_t
	if cmd.name_orEmpty != "" {
		info = namedObjects[cmd.name_orEmpty]
	} else {
		info = objects[objectType]
	}

	if info == nil {
		info = &objectInfo_t{
			name_orEmpty: cmd.name_orEmpty,
			objectType:   objectType,
			object_orNil: cmd.object,
			waiters:      nil,
		}

		if cmd.name_orEmpty != "" {
			namedObjects[cmd.name_orEmpty] = info
		} else {
			objects[objectType] = info
		}
	} else {
		if info.object_orNil != nil {
			var err error
			if cmd.name_orEmpty != "" {
				err = errors.New("conflict with an already published object (name \"" + cmd.name_orEmpty + "\")")
			} else {
				err = errors.New("conflict with an already published object (type " + objectType.String() + ")")
			}
			cmd.resultCh <- cmd_publish_result{
				pub: nil,
				err: err,
			}
			return
		}

		info.object_orNil = cmd.object

		// Unblock all waiters
		waiters := info.waiters
		info.waiters = nil
		for _, w := range waiters {
			w <- cmd_wait_result{
				object: cmd.object,
			}
		}
	}

	cmd.resultCh <- cmd_publish_result{
		pub: info,
		err: nil,
	}
}

// Find
type cmd_find struct {
	objectType reflect.Type
	resultCh   chan<- cmd_find_result
}

type cmd_findName struct {
	name     string
	resultCh chan<- cmd_find_result
}

type cmd_find_result struct {
	object_orNil interface{}
}

func do_find(cmd cmd_find) {
	var object_orNil interface{} = nil
	{
		info := objects[cmd.objectType]
		if info != nil {
			object_orNil = info.object_orNil
		}
	}

	cmd.resultCh <- cmd_find_result{
		object_orNil: object_orNil,
	}
}

func do_findName(cmd cmd_findName) {
	var object_orNil interface{} = nil
	{
		info := namedObjects[cmd.name]
		if info != nil {
			object_orNil = info.object_orNil
		}
	}

	cmd.resultCh <- cmd_find_result{
		object_orNil: object_orNil,
	}
}

// Wait
type cmd_wait struct {
	objectType reflect.Type
	resultCh   chan<- cmd_wait_result
}

type cmd_waitName struct {
	name     string
	resultCh chan<- cmd_wait_result
}

type cmd_wait_result struct {
	object interface{}
}

func do_wait(cmd cmd_wait) {
	info := objects[cmd.objectType]
	if info == nil {
		// Wait for the object to be published
		info = &objectInfo_t{
			objectType:   cmd.objectType,
			object_orNil: nil,
			waiters:      []chan<- cmd_wait_result{cmd.resultCh},
		}
		objects[cmd.objectType] = info
	} else if info.object_orNil == nil {
		// Wait for the object to be published
		info.waiters = append(info.waiters, cmd.resultCh)
	} else {
		// Object is already published - no need to wait
		cmd.resultCh <- cmd_wait_result{
			object: info.object_orNil,
		}
	}
}

func do_waitName(cmd cmd_waitName) {
	info := namedObjects[cmd.name]
	if info == nil {
		// Wait for the object to be published
		info = &objectInfo_t{
			name_orEmpty: cmd.name,
			objectType:   nil,
			object_orNil: nil,
			waiters:      []chan<- cmd_wait_result{cmd.resultCh},
		}
		namedObjects[cmd.name] = info
	} else if info.object_orNil == nil {
		// Wait for the object to be published
		info.waiters = append(info.waiters, cmd.resultCh)
	} else {
		// Object is already published - no need to wait
		cmd.resultCh <- cmd_wait_result{
			object: info.object_orNil,
		}
	}
}

// PublishedResult.Remove
type cmd_remove struct {
	name_orEmpty string
	objectType   reflect.Type
	resultCh     chan<- cmd_remove_result
}

type cmd_remove_result struct {
	err error
}

func do_remove(cmd cmd_remove) {
	if cmd.name_orEmpty != "" {
		info := namedObjects[cmd.name_orEmpty]
		if info != nil {
			if len(info.waiters) > 0 {
				panic("this cannot happen")
			}

			delete(namedObjects, cmd.name_orEmpty)
			cmd.resultCh <- cmd_remove_result{
				err: nil,
			}
		} else {
			cmd.resultCh <- cmd_remove_result{
				err: errors.New("no such object (name \"" + cmd.name_orEmpty + "\")"),
			}
		}
	} else {
		info := objects[cmd.objectType]
		if info != nil {
			if len(info.waiters) > 0 {
				panic("this cannot happen")
			}

			delete(objects, cmd.objectType)
			cmd.resultCh <- cmd_remove_result{
				err: nil,
			}
		} else {
			cmd.resultCh <- cmd_remove_result{
				err: errors.New("no such object (type " + cmd.objectType.String() + ")"),
			}
		}
	}
}

// The main command loop, running in a separate goroutine
func commandLoop() {
	for untypedCommand := range commandChannel {
		switch cmd := untypedCommand.(type) {
		case cmd_publish:
			do_publish(cmd)
		case cmd_find:
			do_find(cmd)
		case cmd_findName:
			do_findName(cmd)
		case cmd_wait:
			do_wait(cmd)
		case cmd_waitName:
			do_waitName(cmd)
		case cmd_remove:
			do_remove(cmd)
		}
	}
}

func init() {
	go commandLoop()
}
