package asyncio_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/Drelf2018/asyncio"
)

func sleep(second int) float64 {
	time.Sleep(time.Duration(second) * time.Second)
	return math.Sqrt(float64(second))
}

type Student struct {
	Name string
}

func (s Student) Introduce(more string) {
	time.Sleep(time.Duration(2) * time.Second)
	println("My name is", s.Name, "and", more)
}

func (s Student) Hello() {
	print("Hello ")
}

func (s Student) Me() {
	time.Sleep(time.Second)
	println(s.Name)
}

func TestSleep(t *testing.T) {
	handles := asyncio.Slice(sleep, asyncio.SingleArg(1, 2, 3, 4))

	// need hint type
	a := asyncio.To[[]float64](handles)
	fmt.Printf("a: %v\n", a)
	// auto infer
	b := make([]float64, handles.Len())
	asyncio.To(handles, b)
	fmt.Printf("b: %v\n", b)

	for i, handle := range handles {
		fmt.Printf("No.%v sleep() return %v\n", i, handle.Result())
	}
}

func TestStruct(t *testing.T) {
	s := Student{"Alice"}
	coro := asyncio.C(s.Introduce, "I'm glad to see you!")
	coros := asyncio.NoArgsFunc(s.Hello, s.Me)
	asyncio.Await(append(coros, coro)...)
}
