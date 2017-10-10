
# Agent [![CircleCI](https://circleci.com/gh/lodastack/agent.svg?style=svg&circle-token=b26ad578124d061da19fb8cd796bc12b0d1393bd)](https://circleci.com/gh/lodastack/agent)

## Build

    make build

## Getting Started

### Start agent
    
    ./agent start -f ${path_to_config_file}

### Stop agent

    ./agent stop

### Reset agent, clean all plugins data and local cache data

    ./agent stop -m clean



## Configuration

```
[agent]
	# HTTP API listen
	listen = "0.0.0.0:1232" 
	# which the network interface will be monitored
	ifaceprefix = [ "eth" ]
	# registry server address
	registryaddr = "https://registry.test.com"
	# which directory to store plugins 
	pluginsdir = "/usr/local/agent-plugins"
	# run plugins user, must exist on your system
	pluginsuser = "user"
	# your plugins git organizantion
	git = "git://git@git.test.com/orgname/%s.git"

[output]
	# message queue, now only support NSQ
	name = "nsq"
	# MQ addresses
	servers = [ "0.0.0.0:7777" ]
	# how many points cached in local memory
	buffersize = 1000

[log]
	# log directory
	logdir = "/tmp/agent/log"
	# log level
	loglevel = "DEBUG"
	# how many log files retention
	logrotatenum = 5
	# log file size (unit:byte)
	logrotatesize = 1887436800


```

## Plugin
- 示例：https://git.test.com/XXX/plugin-example
- 插件放在git上，每个插件的根目录下必须包含一个名为plugin的可执行文件
- agent主动从git拉取，放在本地的PluginsDir里面；agent不会更新已经缓存在本地的插件，提供手动更新插件的http接口
- 建议使用非root权限运行插件。运行时cd到对应目录下，执行plugin
- registry提供每个ns对应的插件的git地址和运行周期（单位秒）、运行参数（可选）
	- 运行周期为0：不是定期运行的脚本，由外部调用的接口运行。见http接口
	- 定期运行的脚本：
    	- 要求插件的输出是一个json序列化的metric的list，提交的时候将measurement修改为PLUGIN.插件的name.metric的name

```
    	       type Metric struct {
        	      Name      string            `json:"name"`
        	      Timestamp int64             `json:"timestamp"`
        	      Tags      map[string]string `json:"tags"`
        	      Value     interface{}       `json:"value"`
    	       }
```

- 提供disable和enable某个插件的http接口。如果disable之后一段时间（一天）没有enable，将会自动enable
- 使用ns和git项目名作为索引，相同ns下的插件项目**不能同名**
- 插件运行的结果也是一种metric，名字为plugin.Name，提交到对应的NS中


## HTTP API

- /post: push data to agent：curl -d ${metrics} "http://host:port/post?ns=xxx"
- /update: update local collect resource 
- /me/ns: NS list  
- /me/status: agent version
- /plugins/list: 获取当前的插件状态（是否enable）列表
- 下面各种接口都必须有两个参数ns和repo。repo是完整的gitlab地址，比如git@git.test.com:XXX/plugin-example.git。对应的插件配置必须已经在tree配置。新增加的插件可能会因为agent没有及时更新而报错（agent每隔十分钟从tree拉取一次）
- /plugins/update?ns=xxx&repo=xxx: 更新本地缓存的插件
- /plugins/run?ns=xxx&repo=xxx&timeout=20&param=xxx: 运行插件（一般是外部调用的插件，即在tree的配置中采集周期为0），timeout和param可选，timeout单位为秒，默认10s。如果调用时当前插件正在执行中，则忽略当前调用。
- /plugins/{enbable|disable}?ns=xxx&repo=xxx: 启用/禁用某个插件


## Other

### nsq && influxdb
- 按照influxdb的数据格式发送，database是namespace；对每个point添加host的tag；如果point的时间戳为0或者单位不是秒，改为当前的时间。
- timestamp的precision是秒