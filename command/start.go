package command

import (
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/lodastack/agent/agent/agent"
	"github.com/lodastack/agent/config"
	"github.com/lodastack/log"

	"github.com/oiooj/cli"
)

var logBackend *log.FileBackend

var CmdStart = cli.Command{
	Name:        "start",
	Usage:       "启动客户端",
	Description: "启动Agent客户端",
	Action:      runStart,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "f",
			Value: "/etc/agent.conf",
			Usage: "配置文件路径，默认位置：/etc/agent.conf",
		},
	},
}

func runStart(c *cli.Context) {
	//parse config file
	err := config.ParseConfig(c.String("f"))
	if err != nil {
		log.Fatalf("Parse Config File Error: %s", err.Error())
	}
	//init log setting
	initLog()
	//save pid to file
	ioutil.WriteFile(config.PID, []byte(strconv.Itoa(os.Getpid())), 0744)
	go Notify()

	//start agent module
	a, err := agent.New(config.C)
	if err != nil {
		log.Fatalf("New agent Error: %s", err.Error())
	}
	a.Start()
	// Print sweet Agent logo.
	PrintLogo()
	select {}
}

func initLog() {
	var err error
	logBackend, err = log.NewFileBackend(config.C.Log.Dir)
	if err != nil {
		log.Fatalf("failed to new log backend")
	}
	log.SetLogging(config.C.Log.Level, logBackend)
	log.Rotate(config.C.Log.Logrotatenum, config.C.Log.Logrotatesize)
}

func Notify() {
	message := make(chan os.Signal, 1)

	signal.Notify(message, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGKILL, os.Interrupt)
	<-message
	log.Info("receive signal, exit...")
	logBackend.Flush()
	stopProfile()
	os.Exit(0)
}

func PrintLogo() {
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
