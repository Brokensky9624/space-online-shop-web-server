package common

type Stoppable interface {
	Stop()
}

func StopAll(stoppers ...Stoppable) {
	for _, stopper := range stoppers {
		stopper.Stop()
	}
}
