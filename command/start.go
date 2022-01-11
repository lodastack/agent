package command

import (
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"

	"github.com/lodastack/agent/agent/agent"
	"github.com/lodastack/agent/config"
	"github.com/lodastack/agent/member"
	"github.com/lodastack/agent/trace"
	"github.com/lodastack/log"

	"github.com/oiooj/cli"
)

var logBackend *log.FileBackend

var CmdStart = cli.Command{
	Name:        "start",
	Usage:       "start agent",
	Description: "start agent client",
	Action:      runStart,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "f",
			Value: "/etc/agent.conf",
			Usage: "default config file：/etc/agent.conf",
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

func runStart(c *cli.Context) {
	//parse config file
	err := config.ParseConfig(c.String("f"))
	if err != nil {
		log.Fatalf("Parse Config File Error: %s", err.Error())
	}
	start(c, config.C)
}

func start(c *cli.Context, conf *config.Config) {
	//init log setting
	initLog(conf.Log)

	//start agent module
	a, err := agent.New(conf)
	if err != nil {
		log.Fatalf("New agent Error: %s", err)
	}
	if err := a.Start(); err != nil {
		log.Fatalf("agent start failed: %s", err)
	}
	// Print sweet Agent logo.
	printLogo()
	if runtime.GOOS != "windows" {
		//save pid to file
		ioutil.WriteFile(config.PID, []byte(strconv.Itoa(os.Getpid())), 0644)
		go notify()
	}
	// trace module
	if config.C.Trace.Enable {
		err = trace.Start(config.C.Trace.Collector, config.C.Log.Dir)
		if err != nil {
			log.Errorf("trace module start failed: %s", err)
		}
		log.Info("trace module started")
	}
	// member module
	if config.C.Member.Enable {
		if err := member.Member.Start(config.C.Member.Nodes, config.C.Member.Key); err != nil {
			log.Errorf("member module start failed: %s", err)
		}
		log.Info("member module started")
	}
	//pprof
	startProfile(c.String("cpuprofile"), c.String("memprofile"))
	select {}
}

func initLog(conf config.LogConfig) {
	var err error
	logBackend, err = log.NewFileBackend(conf.Dir)
	if err != nil {
		log.Fatalf("failed to new log backend:%v\n", err)
	}
	log.SetLogging(conf.Level, logBackend)
	logBackend.Rotate(conf.Logrotatenum, conf.Logrotatesize)
}

func notify() {
	message := make(chan os.Signal, 1)

	signal.Notify(message, syscall.SIGINT, syscall.SIGKILL, os.Interrupt)
	<-message
	log.Info("receive signal, exit...")
	logBackend.Flush()
	stopProfile()
	os.Exit(0)
}

func printLogo() {
	log.Printf(logo)
}

const logo = `
                                                                                                               
                                                                                                               
               AAA                                                                               tttt          
              A:::A                                                                           ttt:::t          
             A:::::A                                                                          t:::::t          
            A:::::::A                                                                         t:::::t          
           A:::::::::A           ggggggggg   ggggg    eeeeeeeeeeee    nnnn  nnnnnnnn    ttttttt:::::ttttttt    
          A:::::A:::::A         g:::::::::ggg::::g  ee::::::::::::ee  n:::nn::::::::nn  t:::::::::::::::::t    
         A:::::A A:::::A       g:::::::::::::::::g e::::::eeeee:::::een::::::::::::::nn t:::::::::::::::::t    
        A:::::A   A:::::A     g::::::ggggg::::::gge::::::e     e:::::enn:::::::::::::::ntttttt:::::::tttttt    
       A:::::A     A:::::A    g:::::g     g:::::g e:::::::eeeee::::::e  n:::::nnnn:::::n      t:::::t          
      A:::::AAAAAAAAA:::::A   g:::::g     g:::::g e:::::::::::::::::e   n::::n    n::::n      t:::::t          
     A:::::::::::::::::::::A  g:::::g     g:::::g e::::::eeeeeeeeeee    n::::n    n::::n      t:::::t          
    A:::::AAAAAAAAAAAAA:::::A g::::::g    g:::::g e:::::::e             n::::n    n::::n      t:::::t    tttttt
   A:::::A             A:::::Ag:::::::ggggg:::::g e::::::::e            n::::n    n::::n      t::::::tttt:::::t
  A:::::A               A:::::Ag::::::::::::::::g  e::::::::eeeeeeee    n::::n    n::::n      tt::::::::::::::t
 A:::::A                 A:::::Agg::::::::::::::g   ee:::::::::::::e    n::::n    n::::n        tt:::::::::::tt
AAAAAAA                   AAAAAAA gggggggg::::::g     eeeeeeeeeeeeee    nnnnnn    nnnnnn          ttttttttttt  
                                          g:::::g                                                              
                              gggggg      g:::::g                                                              
                              g:::::gg   gg:::::g                                                              
                               g::::::ggg:::::::g                                                              
                                gg:::::::::::::g                                                               
                                  ggg::::::ggg                                                                 
                                     gggggg                                                                    

`
