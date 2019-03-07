package server

import (
	"fmt"
	"os"
	"time"

	"data/config"

	"github.com/jander/golog/logger"
	"github.com/kardianos/service"
)

func OnStart(callBack func()) {
	name, displayName, desc := config.GetServiceConfig()
	p := &program{callBack}
	sc := &service.Config{
		Name:        name,
		DisplayName: displayName,
		Description: desc,
	}
	s, err := service.New(p, sc)
	//var s, err = service.NewService(name, displayName, desc)
	if err != nil {
		fmt.Printf("%s unable to start: %s", displayName, err)
		return
	}

	fmt.Printf("Service \"%s\" do.\n", displayName)

	if len(os.Args) > 1 {
		var err error
		verb := os.Args[1]
		switch verb {
		case "install":
			{
				err = s.Install()
				if err != nil {
					fmt.Printf("Failed to install: %s\n", err)
					return
				}
				fmt.Printf("Service \"%s\" installed.\n", displayName)
			}
		case "remove":
			{
				err = s.Uninstall()
				if err != nil {
					fmt.Printf("Failed to remove: %s\n", err)
					return
				}
				fmt.Printf("Service \"%s\" removed.\n", displayName)
			}
		case "run":
			{
				err = s.Run()
				if err != nil {
					fmt.Printf("Failed to run: %s\n", err)
					return
				}
				fmt.Printf("Service \"%s\" run.\n", displayName)
			}
		case "start":
			{
				err = s.Start()
				if err != nil {
					fmt.Printf("Failed to start: %s\n", err)
					return
				}
				fmt.Println("starting check service:", displayName)

				ticker := time.NewTicker(1 * time.Second)
				<-ticker.C

				st, err := IsStart(name)
				if err != nil {
					fmt.Println(err)
					return
				} else {
					if st == Stopped || st == StopPending {
						fmt.Printf("Service \"%s\" is Stopped.\n", displayName)
						fmt.Println("can't to start service.")
						return
					}
				}
				fmt.Printf("Service \"%s\" started.\n", displayName)
			}
		case "stop":
			{
				err = s.Stop()
				if err != nil {
					fmt.Printf("Failed to stop: %s\n", err)
					return
				}

				st, err := IsStart(name)
				if err != nil {
					fmt.Println(err)
					return
				} else {
					if st == Running || st == StartPending {
						fmt.Printf("Service \"%s\" is Started.\n", displayName)
						fmt.Println("can't to stop service.")
						return
					}
				}
				fmt.Printf("Service \"%s\" stopped.\n", displayName)
			}
		}
		return
	} else {
		fmt.Print("Failed to read args\n")
		//return
	}

	if err = s.Run(); err != nil {
		logger.Error(err)
	}
}

type program struct {
	callBack func()
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	p.callBack()
}
func (p *program) Stop(s service.Service) error {
	return nil
}

type ServiceTools interface {
	IsStart(name string) (status int, err error)
}
