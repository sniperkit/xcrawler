package ajii

import (
	"fmt"
	"github.com/op/go-logging"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type BaseService struct {
	Config     *EtcdConfig
	Name       string
	WorkerChan chan bool
	Logger     *logging.Logger
	Callback   func(*CallCtx) interface{}
}

type Service interface {
	Announce() (string, error)
	Unannounce(syscall os.Signal) error
	ExitHandler()
	Run()
	Init(*cli.Context)
}

func NewService(name string, callback func(*CallCtx) interface{}) *BaseService {
	return &BaseService{
		// Config:   c,
		Name:     name,
		Logger:   GetLogger("service"),
		Callback: callback,
	}
}

func (s *BaseService) Announce() (string, error) {
	go s.ExitHandler()
	s.Logger.Debugf("Registering %s on with value: %s", s.Name, s.Config.ServiceUrl)
	return s.Config.Set(s.Name, s.Config.ServiceUrl)
}

func (s *BaseService) Unannounce(syscall os.Signal) error {
	s.Logger.Debug("Unannouncing from registry")
	return s.Config.Delete(s.Name)
}

func (s *BaseService) Init(c *cli.Context) {
	s.Logger.Debug("Booting")
	s.Config.Port = c.Int("port")
	s.Config.ServiceUrl += fmt.Sprintf(":%d", s.Config.Port)
	_, err := s.Announce()
	if err != nil {
		log.Fatal(err)
	}
	s.Config.Port = c.Int("port")
	s.WorkerChan = make(chan bool)
	// go BackgroundWorker(s.Config, s.WorkerChan)

	port_config := fmt.Sprintf(":%d", c.Int("port"))
	s.Logger.Debugf("Entering main thingy, ready to serve backend at: %s", port_config)
	Backend(port_config, s.Config, s.Callback)
	// Frontend("0.0.0.0:8020", s.Config)
}

func (s *BaseService) Run() {
	config := NewConfig()
	s.Config = config
	cli := NewCli(config, s)
	cli.Run(os.Args)
}

func (s *BaseService) ExitHandler() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	sig := <-sigc // block until we get ze signal
	s.Logger.Info("Unregistering from registry, reason:", sig)
	s.Unannounce(sig)
	s.Logger.Info("Sending kill-flag to BackgroundWorker")
	// s.WorkerChan <- true

	os.Exit(0)
}
