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
