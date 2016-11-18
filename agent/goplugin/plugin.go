package goplugin

import (
	"github.com/lodastack/agent/agent/common"
)

var funcs map[string]func(map[string]interface{}) ([]*common.Metric, error)

func init() {
	// add plugin name to function here
	// please delete "test" if you add some real function
	funcs = map[string]func(map[string]interface{}) ([]*common.Metric, error){
		"test": test,
	}
}

func test(params map[string]interface{}) (L []*common.Metric, err error) {
	L = append(L, &common.Metric{Name: "test", Value: 100})
	return
}
