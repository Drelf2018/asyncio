package asyncio_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/Drelf2018/asyncio"
)

func TestGlobal(t *testing.T) {
	second := time.Duration(3) * time.Second
	asyncio.RunForever(
		asyncio.CreateTask(
			func(sec time.Duration) {
				time.Sleep(sec)
				asyncio.GetEventLoop().Stop()
			},
			second,
		),
	)
	t.Log("Awake")
}

func sleep(second int, num string) float64 {
	time.Sleep(time.Duration(second) * time.Second)
	fmt.Printf("No.%v Awake.\n", num)
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
	handles := asyncio.Slice([][]any{{1, "1"}, {2, "2"}, {3, "3"}, {4, "4"}}, sleep)

	// need hint type
	a := asyncio.Results[float64](handles)
	fmt.Printf("a: %v\n", a)
	// auto infer
	b := make([]float64, len(handles))
	asyncio.Fill(handles, b)
	fmt.Printf("b: %v\n", b)

	for i, handle := range handles {
		fmt.Printf("No.%v sleep() return %v\n", i, handle.Result())
	}
}

func TestStruct(t *testing.T) {
	s := Student{"Alice"}
	asyncio.Loop(func(loop asyncio.EventLoop) {
		loop.CreateTask(s.Introduce, "I'm glad to see you!")
		loop.CreateTask(s.Hello)
		loop.CreateTask(s.Me)
	})
}

func TestRetry(t *testing.T) {
	count := 0
	res := asyncio.RetryTask(4, 0.5, asyncio.CreateTask(
		func(i *int) bool {
			defer func() {
				count++
			}()
			fmt.Printf("i: %v\n", *i)
			return *i == 3
		},
		&count,
	))
	fmt.Printf("res: %v\n", res)
}

func TestCreate(t *testing.T) {
	data := []int{1, 2, 3, 4}
	asyncio.GetEventLoop().Start()
	asyncio.CreateTask(func(i, j, k, l int) {
		fmt.Printf("i: %v\n", i)
		fmt.Printf("j: %v\n", j)
		fmt.Printf("k: %v\n", k)
		fmt.Printf("l: %v\n", l)
	}, reflect.Slice, data)
	// }, data...) -> cannot use data (variable of type []int) as []any value in argument to
	asyncio.RunUntilComplete()
}

func TestPromise(t *testing.T) {
	count := 0
	for {
		asyncio.Promise(func(i int) bool {
			return i == 3
		}, count).Then(func(res bool) {
			if res {
				println("succeed")
			} else {
				println("failed")
			}
		}).Finally(func() {
			count++
		}).Call()
		if count == 5 {
			break
		}
	}
}

func TestFunc(t *testing.T) {
	asyncio.ForFunc("abc", func(s string) {
		s += "d"
		fmt.Printf("s1: %v\n", s)
	}, func(s string) {
		s += s
		fmt.Printf("s2: %v\n", s)
	})
}

func TestDelay(t *testing.T) {
	t.Log("原神，")
	asyncio.Delay(3, func() { t.Log("启动！") })
	time.Sleep(4 * time.Second)
}
