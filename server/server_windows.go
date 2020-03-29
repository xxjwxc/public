package server

import (
	"github.com/btcsuite/winsvc/mgr"
	"github.com/btcsuite/winsvc/svc"
)

type ServiceTools struct {
	//i IServiceTools
}

func (s *ServiceTools) IsStart(name string) (st int, err error) {
	var m *mgr.Mgr
	m, err = mgr.Connect()
	if err != nil {
		return 0, err
	}
	defer m.Disconnect()

	sv, err := m.OpenService(name) 
	if err != nil {
		return 0, err
	}
	defer sv.Close()

	var ss svc.Status
	ss, err = sv.Query()
	st = int(ss.State)
	return
}
