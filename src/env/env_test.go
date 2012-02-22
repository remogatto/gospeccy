package env

import (
	"reflect"
	"sync"
	"testing"
)

type T struct{}

var waitGroup sync.WaitGroup

func add() {
	waitGroup.Add(1)
	Wait(reflect.TypeOf(T{}))
	waitGroup.Done()
}

func publish_remove() error {
	var t T
	pub, err := Publish(t)
	if err != nil {
		return err
	}
	waitGroup.Wait()
	pub.Remove()
	return nil
}

func Test1(t *testing.T) {
	go add()
	go add()

	err := publish_remove()
	if err != nil {
		t.Fatal(err)
	}
}
