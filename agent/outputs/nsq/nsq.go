package nsq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/lodastack/agent/agent/outputs"

	"github.com/lodastack/log"
)

const NAME = "nsq"
const maxBatchPoints = 30

type NSQ struct {
	Servers []string
}

const timeout = 2 * time.Second

func (n *NSQ) Description() string {
	return "Send measurements to NSQD"
}

func (n *NSQ) SetServers(servers []string) {
	n.Servers = servers
}

func merge(batch outputs.Data, data outputs.Data) (outputs.Data, bool) {
	if batch.Points == nil {
		return data, true
	}
	if data.Namespace == batch.Namespace && len(batch.Points.Points) < maxBatchPoints {
		batch.Points.Points = append(batch.Points.Points, data.Points.Points...)
		return batch, true
	}
	return batch, false
}

func (n *NSQ) Write(queue chan outputs.Data) {
	var batch outputs.Data
	for {
		data := <-queue

		// merge same NS data points, reduce TCP overhead
		batch, merged := merge(batch, data)
		if merged {
			continue
		}
		queue <- data
		data = batch
		batch = outputs.Data{}

		body, err := json.Marshal(data.Points)
		if err != nil {
			log.Error("marshal datapoint:", data, " failed. error:", err)
			continue
		}

		l := len(n.Servers)
		p := rand.Perm(l)
		for i, idx := range p {
			//log.Debug("send to " + Conf.NsqServers[idx])
			err = httpPost(n.Servers[idx], body, data.Namespace)
			if err == nil {
				break
			}
			log.Warning("Publish to nsq failed: ", err)
			if i == l-1 {
				if !strings.Contains(err.Error(), "connection refused") {
					log.Error("send to nsq failed:", err.Error(), "discard message, namespace: ", data.Namespace, " data: ", string(body))
				} else {
					select {
					case queue <- data:
					default:
						log.Error("queue is full, discard message, namespace: ", data.Namespace, " data: ", string(body))
					}
				}
			} else {
				time.Sleep(time.Millisecond * time.Duration(100*i))
			}
		}
	}
}

func httpPost(addr string, data []byte, namespace string) error {
	url := fmt.Sprintf("http://%s/put?topic=%s", addr, namespace)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(timeout)
				c, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("bad response status code: " + resp.Status)
	}
	return nil
}

func init() {
	outputs.Add(NAME, func() outputs.OutputInf {
		return &NSQ{}
	})
}
