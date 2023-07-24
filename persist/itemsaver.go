package persist

import "log"

func ItemSaver() chan interface{} {
	out := make(chan interface{})
	// 真正执行保存数据的操作
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("ItemSaver got item #%d: %v", itemCount, item)
			itemCount++
		}

	}()
	return out
}
