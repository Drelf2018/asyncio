package asyncio_test

import (
	"fmt"
	"testing"

	"github.com/Drelf2018/asyncio"
)

func TestAsyncEvent(t *testing.T) {
	// a := make(asyncio.AsyncEvent[string])
	a := make(asyncio.AsyncEvent[int])

	// a.OnCommand("danmaku114", asyncio.OnlyData(func(data int) {
	// 	fmt.Printf("data: %v(%T)\n", data, data)
	// }))

	a.OnCommand(1, asyncio.OnlyData(func(data int) {
		fmt.Printf("data: %v(%T)\n", data, data)
	}))

	// a.OnRegexp(`danmaku\d`,
	// 	asyncio.WithData(
	// 		func(e *asyncio.Event, data int) {
	// 			if data&1 == 0 {
	// 				e.Store("sin", math.Sin(float64(data)))
	// 			}
	// 		},
	// 	),
	// 	func(e *asyncio.Event) {
	// 		var num float64
	// 		err := e.Get("sin", &num, -1.0)
	// 		utils.PanicErr(err)
	// 		if num == -1.0 {
	// 			println("sin: Didn't store the value of sin(data)")
	// 		} else {
	// 			fmt.Printf("sin: %v(%T)\n", num, num)
	// 		}
	// 	},
	// )

	a.OnAll(
		func(e *asyncio.Event) { fmt.Printf("%v\n", e) },
		func(e *asyncio.Event) { e.Abort() },
		func(e *asyncio.Event) { fmt.Println("Why still running!?") },
	)

	asyncio.Heartbeat(0, 2, asyncio.WithData(func(e *asyncio.Event, count int) {
		println()
		// a.Dispatch("danmaku114", count)
		a.Dispatch(1, count)
		if count == 5 {
			e.Abort()
		}
	}))
}
