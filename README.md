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
	handles := asyncio.Slice(sleep, asyncio.SingleArg(1, 2, 3, 4))
	h := asyncio.H[float64](handles).To()
	fmt.Printf("h: %v\n", h)
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

```c
// console
h: [1 1.4142135623730951 1.7320508075688772 2]
No.0 sleep() return [1]
No.1 sleep() return [1.4142135623730951]
No.2 sleep() return [1.7320508075688772]
No.3 sleep() return [2]
Hello Alice
My name is Alice and I'm glad to see you!
PASS
ok      github.com/Drelf2018/asyncio    6.044s
```