package server

import (
	"github.com/btcsuite/winsvc/mgr"
	"github.com/btcsuite/winsvc/svc"
)

type WindowsServiceTools struct {
	i ServiceTools
}

func IsStart(name string) (st int, err error) {
	var m *mgr.Mgr
	m, err = mgr.Connect()
	if err != nil {
		return 0, err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	var ss svc.Status
	ss, err = s.Query()
	st = int(ss.State)
	return
}
