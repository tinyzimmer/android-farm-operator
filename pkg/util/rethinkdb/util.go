package rethinkdb

import (
	"fmt"
	"strings"

	androidv1alpha1 "github.com/tinyzimmer/android-farm-operator/pkg/apis/android/v1alpha1"
	stfutil "github.com/tinyzimmer/android-farm-operator/pkg/util/stf"
	"github.com/go-logr/logr"
	rdb "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func EnsureRethinkDBReplicas(reqLogger logr.Logger, instance *androidv1alpha1.AndroidFarm) error {
	if *instance.STFConfig().RethinkDBReplicas() == 1 {
		return nil
	}
	session, err := rdb.Connect(rdb.ConnectOpts{
		Address: strings.TrimPrefix(stfutil.RethinkDBProxyEndpoint(instance), "tcp://"),
	})
	if err != nil {
		return err
	}
	defer session.Close()
	res, err := rdb.DB("stf").TableList().Run(session)
	if err != nil {
		return err
	}
	var tables []string
	if err = res.All(&tables); err != nil {
		return err
	}
	desiredShards := *instance.STFConfig().RethinkDBShards()
	desiredReplicas := *instance.STFConfig().RethinkDBReplicas()
	for _, table := range tables {
		shards, replicas, err := getTableConfig(session, table)
		if err != nil {
			return err
		}
		reqLogger.Info(fmt.Sprintf("RDB table %s has %d shards with %d replicas", table, shards, replicas))
		if replicas != desiredReplicas || shards != desiredShards {
			reqLogger.Info(fmt.Sprintf("Updating RDB table %s to have %d shards with %d replicas", table, desiredShards, desiredReplicas))
			if _, err = rdb.DB("stf").Table(table).Reconfigure(rdb.ReconfigureOpts{
				Replicas: desiredReplicas,
				Shards:   desiredShards,
			}).Run(session); err != nil {
				return err
			}
		}
	}
	return nil
}

func getTableConfig(session *rdb.Session, table string) (shards, replicas int32, err error) {
	var res *rdb.Cursor
	res, err = rdb.DB("stf").
		Table(table).
		Config().
		Field("shards").
		Nth(0).
		Field("replicas").
		Count().
		Run(session)
	if err != nil {
		return
	}
	if err = res.One(&replicas); err != nil {
		return
	}
	res, err = rdb.DB("stf").
		Table(table).
		Config().
		Field("shards").
		Count().
		Run(session)
	if err != nil {
		return
	}
	err = res.One(&shards)
	return
}
