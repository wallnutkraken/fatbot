package fatai


import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/unixpickle/serializer"
	"github.com/wallnutkraken/char-rnn"
)

const (
	ChunkSize = 128
)

type LSTMWrapper struct {
	network *charrnn.LSTM
	settings LSTMSettings
}

type LSTMSettings struct {
	SavePath string
	WordCount string
}

func New(s LSTMSettings) *LSTMWrapper {
	wrapper :=  &LSTMWrapper{
		network: &charrnn.LSTM{},
		settings: s,
	}

	return wrapper
}

func (w *LSTMWrapper) Train(data []string) {
	start := time.Now()
	samples := w.loadSamples(data)
	w.network.Train(samples)
	logrus.Infof("Finished training in %s", time.Since(start).String())
}

func (w *LSTMWrapper) loadSamples(data []string) charrnn.SampleList {
	allText := strings.Join(data, "\n")
	var result charrnn.SampleList

	for index := 0; index < len(allText); index += ChunkSize {
		bs := ChunkSize
		if bs > len(allText) - index {
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