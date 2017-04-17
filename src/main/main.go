package main

import (
	"gochat"
	"fmt"
)

type MyListener struct {

}

func (this *MyListener) Handle(event gochat.Event) error {
	if event.Type == gochat.MsgEvent{
		msgEventData, ok := event.Data.(gochat.MsgEventData)
		if ok {
			fmt.Println(msgEventData.Content)
		}
	}
	return nil
}

func main() {

	wechat := gochat.NewWechat(gochat.Option{
		StorageDirPath: "D:/develop/",
	})
	var listener gochat.Listener = new(MyListener)
	err := wechat.SetListener(listener).Run()
	if err != nil {
		fmt.Println(err)
	}
}
