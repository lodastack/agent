
# Agent [![CircleCI](https://circleci.com/gh/lodastack/agent.svg?style=svg&circle-token=b26ad578124d061da19fb8cd796bc12b0d1393bd)](https://circleci.com/gh/lodastack/agent)

## Install

    make instll
    
## Start agent
    
    agent start -f ${path_to_config_file}

## Stop agent

    agent stop

## Interface

### 从registry获得信息
- namespaces
- 插件信息
- 日志采集项信息
- 需要监控的端口和进程

### scheduler
- 每个scheduler定时采集信息，发送到nsq。一个scheduler可能是：
    - 系统信息
    - 端口或者进程信息
    - 一个插件
    - 一个日志采集项
- 定时（目前是每隔10分钟）从registry拉取各种信息；提供手动更新的http接口

#### 系统信息
- 向每个namespace都发送一份

#### 插件
- 示例：https://git.test.com/XXX/plugin-example
- 插件放在git上，每个插件的根目录下必须包含一个名为plugin.sh的bash脚本
- agent主动从git拉取，放在本地的PluginsDir里面；agent不会更新已经缓存在本地的插件，提供手动更新插件的http接口
- 使用非root权限运行插件。运行时cd到对应目录下，执行plugin.sh
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
- 插件运行的结果也是一种metric，名字为plugin.result，提交到 collect.plugin.monitor.didi.com, tags: {"type": "regular"（常规采集脚本）/"call"（外部调用）, "name": 项目名, "ns": namespace}；value: 0正常运行/1运行中遇到错误

#### 内置插件
- 内置插件的代码放在src/goplugin下，每个插件是一个函数 func(map[string]interface{}) (error, []*common.Metric)，goplugin.Init()建立名字到函数的映射。agent调用对应函数，并且将返回值发送到nsq（修改measurement同插件）

#### 日志

### nsq && influxdb
- 按照influxdb的数据格式发送，database是namespace；对每个point添加host的tag；如果point的时间戳为0或者单位不是秒，改为当前的时间。
- timestamp的precision是秒
- TODO: 打包发送

### http接口
- 插件相关
    - 获取当前的插件状态（是否enable）列表：plugins/list
    - 下面各种接口都必须有两个参数ns和repo。repo是完整的gitlab地址，比如git@git.ifengidc.com:XXX/plugin-example.git。对应的插件配置必须已经在tree配置。新增加的插件可能会因为agent没有及时更新而报错（agent每隔十分钟从tree拉取一次）
        - 更新本地缓存的插件：/plugins/update?ns=xxx&repo=xxx
        - 运行插件（一般是外部调用的插件，即在tree的配置中采集周期为0）：/plugins/run?ns=xxx&repo=xxx&timeout=20&param=xxx，timeout和param可选，timeout单位为秒，默认10s。如果调用时当前插件正在执行中，则忽略当前调用。
        - 启用/禁用某个插件：plugins/{enbable|disable}?ns=xxx&repo=xxx
- 推送数据到agent：curl -d ${metrics} "http://host:port/post?ns=xxx"
- 手动更新本地采集项配置: /update   
- 获取本机在服务树的节点列表: /me/ns   
- 获取agent版本信息: /me/status
- 日志相关  
    - 从文件的偏移量开始读取内容. log/offset?fpath=xxx&offset=yyy&num=zzz  (fpath: 日志完整路径, offset: 读取文件内容偏移量, num: 读取行数)  
