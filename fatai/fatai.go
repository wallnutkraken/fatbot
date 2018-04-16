package fatai

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"sync"

	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/unixpickle/serializer"
	"github.com/wallnutkraken/char-rnn"
)

const (
	ChunkSize = 128
)

type LSTMWrapper struct {
	network          charrnn.Model
	settings         LSTMSettings
	closeChan        chan struct{}
	stopTrainFor     chan bool
	closer           struct{}
	forTraining      bool
	forTrainingMutex *sync.Mutex
}

func (l *LSTMWrapper) IsForTraining() bool {
	l.forTrainingMutex.Lock()
	defer l.forTrainingMutex.Unlock()
	return l.forTraining
}

type LSTMSettings struct {
	SavePath  string
	WordCount string
}

func New(s LSTMSettings) (*LSTMWrapper, error) {
	wrapper := &LSTMWrapper{
		network:          &charrnn.LSTM{},
		settings:         s,
		closeChan:        make(chan struct{}, 0),
		stopTrainFor:     make(chan bool, 0),
		forTrainingMutex: &sync.Mutex{},
	}
	if _, err := os.Stat(s.SavePath); err == nil {
		file, err := ioutil.ReadFile(s.SavePath)
		if err != nil {
			return nil, err
		}
		deserialized, err := serializer.DeserializeWithType(file)
		if err != nil {
			return nil, err
		}
		var ok bool
		wrapper.network, ok = deserialized.(charrnn.Model)
		if !ok {
			return nil, errors.New("Brain model is not a charrnn model")
		}
	}

	return wrapper, nil
}

func (w *LSTMWrapper) Train(data []string) {
	start := time.Now()
	samples := w.loadSamples(data)
	w.network.Train(samples, w.closeChan)
	logrus.Infof("Finished training in %s", time.Since(start).String())
}

func (w *LSTMWrapper) Stop() {
	w.closeChan <- w.closer
	w.forTrainingMutex.Lock()
	w.forTraining = false
	w.forTrainingMutex.Unlock()
}

func (w *LSTMWrapper) TrainFor(data []string, duration time.Duration) {
	w.forTrainingMutex.Lock()
	if w.forTraining {
		w.Stop()
	} else {
		w.forTraining = true
	}
	w.forTrainingMutex.Unlock()

	go w.Train(data)

	select {
	case <-time.After(duration):
		w.Stop()
	case <-w.stopTrainFor:
		w.Stop()
	}
}

func (w *LSTMWrapper) StopTrainFor() {
	w.stopTrainFor <- true
}

func (w *LSTMWrapper) StartTraining(data []string, callback func()) {
	go func() {
		w.Train(data)
		callback()
	}()
}

func (w *LSTMWrapper) loadSamples(data []string) charrnn.SampleList {
	allText := strings.Join(data, "\n")
	var result charrnn.SampleList

	for index := 0; index < len(allText); index += ChunkSize {
		bs := ChunkSize
		if bs > len(allText)-index {
			bs = len(allText) - index
		}

		result = append(result, []byte(allText[index:index+bs]))
	}

	return result
}

func (w *LSTMWrapper) Save() error {
	encoded, err := serializer.SerializeWithType(w.network)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(w.settings.SavePath, encoded, os.ModePerm)
}

func (w *LSTMWrapper) Generate() string {
	w.network.GenerationFlags().Parse([]string{w.settings.WordCount})
	return w.network.Generate()
}
