package fatbrain

import "sync"

type ChainStatus struct {
	isTraining bool
	trainMutex *sync.Mutex
}

func (c *ChainStatus) IsTraining() bool {
	c.trainMutex.Lock()
	defer c.trainMutex.Unlock()
	return c.isTraining
}

func (c *ChainStatus) SetTraining(status bool) {
	c.trainMutex.Lock()
	c.isTraining = status
	c.trainMutex.Unlock()
}

func newChainStatus() *ChainStatus {
	return &ChainStatus{
		trainMutex: &sync.Mutex{},
	}
}
