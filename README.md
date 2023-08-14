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

func TestMain(t *testing.T) {
	handles := asyncio.Slice(sleep, asyncio.SingleArg(1, 2, 3, 4))
	h := asyncio.H[float64](handles).To()
	fmt.Printf("h: %v\n", h)
	for i, handle := range handles {
		fmt.Printf("No.%v sleep() return %v\n", i, handle.Result())
	}
}
```

```go
// console
h: [1 1.4142135623730951 1.7320508075688772 2]
No.0 sleep() return [1]
No.1 sleep() return [1.4142135623730951]
No.2 sleep() return [1.7320508075688772]
No.3 sleep() return [2]
PASS
ok      github.com/Drelf2018/asyncio    4.041s
```