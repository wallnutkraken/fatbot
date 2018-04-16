package fatcaller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"bytes"

	"github.com/pkg/errors"
	"github.com/wallnutkraken/fatbot/fatctrl/ctrltypes"
)

type Client struct {
	addr   string
	client *http.Client
}

func read(resp *http.Response, jsonPtr interface{}) error {
	defer resp.Body.Close()
	readBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(readBytes, jsonPtr)
}

func New(addr string) *Client {
	if !strings.HasSuffix(addr, "/") {
		addr += "/"
	}
	return &Client{
		addr:   addr,
		client: &http.Client{},
	}
}

func (c *Client) StopTraining() error {
	req, err := http.NewRequest("PATCH", c.addr+"training/status/stop", nil)
	if err != nil {
		errors.Wrap(err, "StopTraining: NewRequest")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "StopTraining: Do")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("StopTraining: Failed setting status to \"stop\"")
	}
	return nil
}

func (c *Client) StartTraining(request ctrltypes.StartTrainingRequest) error {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "StartTraining: Marshal")
	}
	req, err := http.NewRequest("PATCH", c.addr+"training/status/train", bytes.NewReader(reqBytes))
	if err != nil {
		return errors.Wrap(err, "StartTraining: NewRequest")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "StartTraining: Do")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("StartTraining: Failed setting status to \"train\"")
	}
	return nil
}

func (c *Client) GetStatus() (ctrltypes.StatusResponse, error) {
	req, err := http.NewRequest("GET", c.addr+"training/status", nil)
	if err != nil {
		return ctrltypes.StatusResponse{}, errors.Wrap(err, "GetStatus: NewRequest")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return ctrltypes.StatusResponse{}, errors.Wrap(err, "GetStatus: Do")
	}
	var status ctrltypes.StatusResponse
	err = read(resp, &status)
	if err != nil {
		err = errors.Wrap(err, "GetStatus: read")
	}
	return status, err
}
