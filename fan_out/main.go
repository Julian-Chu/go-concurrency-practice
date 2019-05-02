package main

func main() {
	Convert(make(chan int))
}

func Convert(ch <-chan int) {

}
func Split(ch <-chan int, n int) []<-chan int {
	var cs []chan int
	for i := 0; i < n; i++ {
		cs = append(cs, make(chan int))
	}

	distributeToChannels := func(ch <-chan int, cs []chan int) {
		defer func(cs []chan int) {
			for _, c := range cs {
				close(c)
			}
		}(cs)

		for {
			for _, c := range cs {
				select {
				case val, ok := <-ch:
					if !ok {
						return
					}
					c <- val
				}
			}
		}
	}

	go distributeToChannels(ch, cs)

	return func(cs []chan int) []<-chan int {
		chs := make([]<-chan int, 0)
		for _, c := range cs {
			chs = append(chs, c)
		}
		return chs
	}(cs)
}
