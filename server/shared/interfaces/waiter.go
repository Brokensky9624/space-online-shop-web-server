package interfaces

type IWaiter interface {
	Wait()
}

func WaitAll(waiters ...IWaiter) {
	for _, waiter := range waiters {
		waiter.Wait()
	}
}
