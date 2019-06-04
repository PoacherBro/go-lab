package test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 主要是为了测试Cancel Context在父子关系中如何表现
// 测试结果为：父context cancel后，子context也会被标记为Done；反过来则不行。
func TestCancelContext(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.Background())
	child, cancelChild := context.WithCancel(parent)

	var wg sync.WaitGroup
	finish := make(chan bool)
	inited := make(chan bool, 2)

	ctxMsg := map[context.Context]string{
		parent: "Parent is done",
		child:  "Child is done",
	}

	for ctx, msg := range ctxMsg {
		c := ctx
		m := msg
		wg.Add(1)
		go func() {
			inited <- true
			defer wg.Done()
			select {
			case <-c.Done():
				fmt.Println(m)
			case <-finish:
				fmt.Println("Finished")
			}

			time.Sleep(time.Second) // make sure log printed
		}()
	}

	i := 0
	for {
		<-inited
		i++
		if i == 2 {
			break
		}
	}
	close(inited)

	// 通过调整两个cancel的顺序，测试子context是否可以cancel父context。经测试，是不行的
	// 但是父context可以cancel子context

	//fmt.Println("cancelParent")
	//cancelParent()

	//time.Sleep(time.Second)

	fmt.Println("cancelChild")
	cancelChild()

	time.Sleep(time.Second)

	fmt.Println("cancelParent")
	cancelParent()

	close(finish)
	wg.Wait()
	time.Sleep(time.Second)
}
