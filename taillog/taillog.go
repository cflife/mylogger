package taillog

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"project1/kafka"
)

var (
	tailObj *tail.Tail
)

type TailTask struct {
	path       string
	topic      string
	instance   *tail.Tail
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path, topic string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:       path,
		topic:      topic,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	tailObj.Init()
	return
}

func (t *TailTask) Init() {
	config := tail.Config{
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		ReOpen:    true,  //重新打开
		MustExist: false, //不存在报错
		Poll:      true,
		Follow:    true, //是否跟随
	}

	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println("tail file failed,err:", err)
	}
}

func (t *TailTask) run() {
	for {
		select {
		case <-t.ctx.Done():
			fmt.Printf("tail task:%s_%s 结束了...\n", t.path, t.topic)
		case line := <-t.instance.Lines:
			fmt.Printf("get log data from %s success, log:%v\n", t.path, line.Text)
			kafka.SendToChan(t.topic, line.Text)
		}
	}
}
