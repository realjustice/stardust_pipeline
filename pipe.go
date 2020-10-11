package pipeline

import "sync"

type PipeChan chan interface{}

// 入口参数
type Entries []interface{}

type PipeCmdFunc func(in PipeChan) PipeChan

// 用户传入的自定义func，用户需要保证每个管道入口的参数和自定义func参数的一致性
type CustomFunc func(interface{}) interface{}

type Pipe struct {
	Cmd   PipeCmdFunc
	Count int
}

func NewPipe() *Pipe {
	return &Pipe{Count: 1}
}

func (p *Pipe) SetPipeCmd(f CustomFunc, count int) {
	p.Cmd = func(in PipeChan) PipeChan {
		out := make(PipeChan)
		go func() {
			defer close(out)
			for inValue := range in {
				out <- f(inValue) // 业务执行
			}
		}()
		return out
	}
	p.Count = count
}


func Run(args Entries, pipes ...*Pipe) PipeChan {
	// 获得入口PipeChan
	in := func(args Entries) PipeChan {
		out := make(PipeChan)
		go func() {
			defer close(out)
			for _, arg := range args {
				out <- arg
			}
		}()
		return out
	}(args)

	for _, pipe := range pipes {
		// 多路复用模块
		in = MultipartExec(in, pipe)
	}

	return in
}

func MultipartExec(in PipeChan, pipe *Pipe) (out PipeChan) {
	out = make(PipeChan)
	wg := sync.WaitGroup{}
	for i := 0; i < pipe.Count; i++ {
		getChan := pipe.Cmd(in)
		wg.Add(1)
		go func(pipeChan PipeChan) {
			defer wg.Done()
			for value := range pipeChan {
				out <- value
			}
		}(getChan)
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}
