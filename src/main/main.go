package main

import (
	"gochat"
	"fmt"
)

type myListener struct {

}

func (this *myListener) handle(event gochat.Event) error {
	fmt.Println(event.Type)
	return nil
}

func main() {

	v := interface{}(&myListener{})
	h, ok := v.(gochat.Listener)
	if ok {
		fmt.Println("a")
	} else {
		wechat := gochat.NewWechat(gochat.Option{
			StorageDirPath: "",
		})
		err := wechat.SetListener(h).Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}
