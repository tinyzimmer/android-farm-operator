package rethinkdb

import (
	rdb "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

const (
	StatusOffline      = 1
	StatusUnauthorized = 2
)

type RethinkDBSession interface {
	GetAllDevicesForProvider(provider string) ([]string, error)
	GetDevicesForProviderByStatus(provider string, status int) ([]string, error)
	Close() error
}

type rethinkDBSession struct {
	session *rdb.Session
}

func NewSession(addr string) (RethinkDBSession, error) {
	session, err := rdb.Connect(rdb.ConnectOpts{
		Address: addr,
	})
	if err != nil {
		return nil, err
	}
	return &rethinkDBSession{session: session}, nil
}

func (r *rethinkDBSession) Close() error {
	return r.session.Close()
}

func (r *rethinkDBSession) GetAllDevicesForProvider(provider string) ([]string, error) {
	res, err := rdb.DB("stf").
		Table("devices").
		Filter(func(uu rdb.Term) rdb.Term {
			return uu.Field("provider").Field("name").Eq(provider)
		}).
		Field("serial").
		Run(r.session)
	if err != nil {
		return nil, err
	}
	if res.IsNil() {
		return []string{}, nil
	} else if res.Err() != nil {
		return nil, res.Err()
	}
	defer res.Close()
	devices := make([]string, 0)
	ch := make(chan string)
	res.Listen(ch)
	for x := range ch {
		devices = append(devices, x)
	}
	return devices, nil
}

func (r *rethinkDBSession) GetDevicesForProviderByStatus(provider string, status int) ([]string, error) {
	res, err := rdb.DB("stf").
		Table("devices").
		Filter(func(uu rdb.Term) rdb.Term {
			return uu.And(uu.Field("provider").Field("name").Eq(provider), uu.Field("status").Eq(status))
		}).
		Field("serial").
		Run(r.session)
	if err != nil {
		return nil, err
	}
	if res.IsNil() {
		return []string{}, nil
	} else if res.Err() != nil {
		return nil, res.Err()
	}
	defer res.Close()
	devices := make([]string, 0)
	ch := make(chan string)
	res.Listen(ch)
	for x := range ch {
		devices = append(devices, x)
	}
	return devices, nil
}
