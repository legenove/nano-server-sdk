package servers

func safeClose(ch chan bool) {
	defer func() {
		if recover() != nil {
		}
	}()
	close(ch)
}

func safeSend(ch chan bool, value bool) {
	defer func() {
		if recover() != nil {
		}
	}()
	ch <- value
}
