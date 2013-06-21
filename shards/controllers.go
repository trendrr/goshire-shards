package shards

import (
	"fmt"
	"github.com/trendrr/goshire/cheshire"
	"github.com/trendrr/goshire/dynmap"
	"log"
	"time"
)

//the global shard manager
var sm *Manager

// This is the global shard manager.
// it is set up when you call shards.RegisterServiceControllers
func SM() *Manager {
	return sm
}

// Sets the partitioner and registers the necessary
// controllers
func RegisterServiceControllers(man *Manager) {
	sm = man

	//register the controllers.
	cheshire.RegisterApi(ROUTERTABLE_GET, "GET", GetRouterTable)
	cheshire.RegisterApi(ROUTERTABLE_SET, "POST", SetRouterTable)
	cheshire.RegisterApi(PARTITION_LOCK, "POST", Lock)
	cheshire.RegisterApi(PARTITION_UNLOCK, "POST", Unlock)
	cheshire.RegisterApi(CHECKIN, "GET", Checkin)
	cheshire.RegisterApi(DATA_PULL, "GET", DataPull)
}

func Checkin(txn *cheshire.Txn) {
	table, err := SM().RouterTable()

	revision := int64(0)
	if err == nil {
		revision = table.Revision
	}
	response := cheshire.NewResponse(txn)
	response.Put("rt_revision", revision)
	response.Put("ts", time.Now())
	txn.Write(response)
}

func GetRouterTable(txn *cheshire.Txn) {
	log.Println("GetRouterTable")
	tble, err := SM().RouterTable()
	if err != nil {
		cheshire.SendError(txn, 506, fmt.Sprintf("Error: %s", err))
		return
	}
	response := cheshire.NewResponse(txn)
	response.Put("router_table", tble.ToDynMap())
	txn.Write(response)
}

func SetRouterTable(txn *cheshire.Txn) {
	rtmap, ok := txn.Params().GetDynMap("router_table")
	if !ok {
		cheshire.SendError(txn, 406, "No router_table")
		return
	}

	rt, err := ToRouterTable(rtmap)
	if err != nil {
		cheshire.SendError(txn, 406, fmt.Sprintf("Unparsable router table (%s)", err))
		return
	}

	_, err = SM().SetRouterTable(rt)
	if err != nil {
		cheshire.SendError(txn, 406, fmt.Sprintf("Unable to set router table (%s)", err))
		return
	}
	response := cheshire.NewResponse(txn)
	txn.Write(response)
}

func Lock(txn *cheshire.Txn) {

	partition, ok := txn.Params().GetInt("partition")
	if !ok {
		cheshire.SendError(txn, 406, fmt.Sprintf("partition param missing"))
		return
	}

	err := SM().LockPartition(partition)
	if err != nil {
		//now send back an error
		cheshire.SendError(txn, 406, fmt.Sprintf("Unable to lock partitions (%s)", err))
		return
	}
	response := cheshire.NewResponse(txn)
	txn.Write(response)
}

func Unlock(txn *cheshire.Txn) {
	partition, ok := txn.Params().GetInt("partition")
	if !ok {
		cheshire.SendError(txn, 406, fmt.Sprintf("partition param missing"))
		return
	}

	err := SM().UnlockPartition(partition)
	if err != nil {
		//now send back an error
		cheshire.SendError(txn, 406, fmt.Sprintf("Unable to lock partitions (%s)", err))
		return
	}
	response := cheshire.NewResponse(txn)
	txn.Write(response)
}

func DataPull(txn *cheshire.Txn) {
	part, ok := txn.Params().GetInt("partition")
	if !ok {
		cheshire.SendError(txn, 406, fmt.Sprintf("partition param is manditory"))
		return
	}
	dataChan := make(chan *dynmap.DynMap, 10)
	finishedChan := make(chan int)
	errorChan := make(chan error)
	go func() {
		for {
			select {
			case data := <-dataChan:
				//send a data packet
				res := cheshire.NewResponse(txn)
				res.SetTxnStatus("continue")
				res.Put("data", data)
				txn.Write(res)
			case <-finishedChan:
				res := cheshire.NewResponse(txn)
				res.SetTxnStatus("complete")
				txn.Write(res)
				return
			case err := <-errorChan:
				cheshire.SendError(txn, 406, fmt.Sprintf("Unable to unlock (%s)", err))
				return
			}
		}
	}()
	SM().service.Data(part, dataChan, finishedChan, errorChan)
}
