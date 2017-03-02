package nsq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/lodastack/agent/agent/outputs"

	"github.com/lodastack/log"
)

const NAME = "nsq"

type NSQ struct {
	Servers []string
}

func (n *NSQ) Description() string {
	return "Send measurements to NSQD"
}

func (n *NSQ) SetServers(servers []string) {
	n.Servers = servers
}

func (n *NSQ) Write(queue chan outputs.Data) {
	for {
		data := <-queue
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
	resp, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(data))
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
