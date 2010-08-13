package perf

import (
	"os"
	"syscall"
	"unsafe"
)


type PerfCounter struct {
	attr Attr
	fd   map[int](*os.File) // File descriptors for each OS thread, initially empty
}


func newPerfCounter() *PerfCounter {
	return &PerfCounter{attr: Attr{}, fd: make(map[int](*os.File))}
}

func (attr *Attr) init_HW(event uint, exclude_user bool, exclude_kernel bool) {
	attr.Type = TYPE_HARDWARE
	attr.Size = ATTR_SIZE
	attr.Config = uint64(event)

	var flags uint64 = 0
	if exclude_user {
		flags |= FLAG_EXCLUDE_USER
	}
	if exclude_kernel {
		flags |= FLAG_EXCLUDE_KERNEL
	}
	flags |= FLAG_EXCLUDE_HV
	flags |= FLAG_EXCLUDE_IDLE
	attr.Flags = flags
}

func (attr *Attr) open(pid int) (counter *os.File, err os.Error) {
	return sys_perf_counter_open(attr, /*pid*/ pid, /*cpu*/ -1, /*group_fd*/ -1, /*flags*/ 0)
}


// Returns a new performance counter for counting CPU cycles
//
// @param user   Specifies whether to count cycles spent in user-space
// @param kernel Specifies whether to count cycles spent in kernel-space
func NewCounter_CpuCycles(user, kernel bool) *PerfCounter {
	counter := newPerfCounter()
	counter.attr.init_HW(HW_CPU_CYCLES, !user, !kernel)
	return counter
}

// Returns a new performance counter for counting retired instructions
//
// @param user   Specifies whether to count instructions executed in user-space
// @param kernel Specifies whether to count instructions executed in kernel-space
func NewCounter_Instructions(user, kernel bool) *PerfCounter {
	counter := newPerfCounter()
	counter.attr.init_HW(HW_INSTRUCTIONS, !user, !kernel)
	return counter
}

// Reads the current value of the performance counter
func (c *PerfCounter) Read() (n uint64, err os.Error) {
	var fd *os.File
	{
		tid := syscall.Gettid()

		var present bool
		fd, present = c.fd[tid]

		if !present {
			fd, err = c.attr.open(tid)
			if err != nil {
				return
			}

			c.fd[tid] = fd
		}
	}

	var buf [8]byte

	var num_read int
	num_read, err = fd.Read(buf[0:8])
	if err != nil {
		return
	}
	if num_read != 8 {
		panic("expected 8 bytes of data")
	}

	n = *(*uint64)(unsafe.Pointer(&buf[0]))
	return
}

func (c *PerfCounter) Close() os.Error {
	var err os.Error = nil

	for _, file := range c.fd {
		err2 := file.Close()
		if err2 != nil {
			// Report only the 1st error
			if err == nil {
				err = err2
			}
		}
	}

	// Clear 'c.fd'
	c.fd = make(map[int](*os.File))

	return err
}
