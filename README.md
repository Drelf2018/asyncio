# asyncio

并发异步库

### 使用

```go
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
	handles := asyncio.Slice(asyncio.SingleArg(1, 2, 3, 4), sleep)

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
```

#### 控制台

```c
// TestSleep
a: [1 1.4142135623730951 1.7320508075688772 2]
b: [1 1.4142135623730951 1.7320508075688772 2]
No.0 sleep() return [1]
No.1 sleep() return [1.4142135623730951]
No.2 sleep() return [1.7320508075688772]
No.3 sleep() return [2]
// TestStruct
Hello Alice
My name is Alice and I'm glad to see you!
```