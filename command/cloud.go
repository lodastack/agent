package command

import (
	"encoding/json"
	"fmt"

	"github.com/lodastack/agent/agent/loda"
	"github.com/lodastack/agent/config"
	"github.com/lodastack/log"

	"github.com/oiooj/cli"
)

var CmdCloudStart = cli.Command{
	Name:        "cloud",
	Usage:       "start agent with cloud.lodastack.com service",
	Description: "start agent client",
	Action:      runCloudStart,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "id",
			Value: "",
			Usage: "instance ID from https://cloud.lodastack.com",
		},
		cli.StringFlag{
			Name:  "token",
			Value: "",
			Usage: "token from https://cloud.lodastack.com",
		},
		cli.StringFlag{
			Name:  "cpuprofile",
			Value: "",
			Usage: "CPU pprof file",
		},
		cli.StringFlag{
			Name:  "memprofile",
			Value: "",
			Usage: "Memory pprof file",
		},
	},
}

func runCloudStart(c *cli.Context) {
	instanceID := c.String("id")
	token := c.String("token")
	if instanceID == "" || token == "" {
		log.Fatalf("Invild instanceID:%s or token:%s, see: https://cloud.lodastack.com\n", instanceID, token)
	}
	conf, err := agentConf(instanceID, token)
	if err != nil {
		log.Fatalf("fetch cloud agent config: %v\n", err)
	}
	start(c, conf)
}

func agentConf(instanceID, token string) (*config.Config, error) {
	url := fmt.Sprintf("%s/api/v1/agent/conf/%s/%s", "https://cloud.lodastack.com", instanceID, token)
	b, err := loda.Get(url)
	if err != nil {
		return nil, err
	}

	type ResponseRes struct {
		Code int           `json:"httpstatus"`
		Data config.Config `json:"data"`
	}
	var response ResponseRes
	err = json.Unmarshal(b, &response)
	if err != nil {
		return nil, err
	}
	c := response.Data
	return &c, nil
}
