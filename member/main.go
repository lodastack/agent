package member

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/memberlist"
)

// Member is global member list
var Member member

// Member struct
type member struct {
	list    *memberlist.Memberlist
	started bool
}

// Start member service
func (m *member) Start(nodes []string, key string) error {
	if Member.started {
		return nil
	}
	if keyBytes, err := base64.StdEncoding.DecodeString(key); err != nil {
		return fmt.Errorf("Invalid key, Decode failed: %s", err)
	} else {
		if err := memberlist.ValidateKey(keyBytes); err != nil {
			return fmt.Errorf("Invalid key, Validate failed: %s", err)
		}
	}

	var err error
	conf := memberlist.DefaultLANConfig()
	conf.TCPTimeout = 5 * time.Second
	conf.GossipInterval = 1 * time.Second
	loadKeyring(conf, []string{key})
	m.list, err = memberlist.Create(conf)
	if err != nil {
		return fmt.Errorf("Failed to create memberlist: %s", err)
	}
	m.started = true
	if len(nodes) == 0 {
		return nil
	}
	// Join an existing cluster by specifying at least one known member.
	_, err = m.list.Join(nodes)
	if err != nil {
		m.started = false
		return fmt.Errorf("Failed to join member cluster: %s", err)
	}
	return nil
}

func loadKeyring(c *memberlist.Config, keys []string) error {
	keysDecoded := make([][]byte, len(keys))
	for i, key := range keys {
		keyBytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return err
		}
		keysDecoded[i] = keyBytes
	}

	if len(keysDecoded) == 0 {
		return fmt.Errorf("no keys present in config keyring")
	}

	keyring, err := memberlist.NewKeyring(keysDecoded, keysDecoded[0])
	if err != nil {
		return err
	}
	c.Keyring = keyring
	return nil
}

func (m *member) List() []memberlist.Node {
	var res []memberlist.Node
	if m.list == nil {
		return res
	}
	for _, member := range m.list.Members() {
		res = append(res, *member)
	}
	return res
}
