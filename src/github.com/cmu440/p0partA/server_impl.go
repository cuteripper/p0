// Implementation of a KeyValueServer. Students should write their code in this file.

package p0partA

import (
	kvstore "github.com/cmu440/p0partA/kvstore"
	"net"
	"strconv"
)

const MaxWriteBuffer = 500

type keyValueServer struct {
	// TODO: implement this!
	store             kvstore.KVStore // KVStore API
	listener          net.Listener    // listener API
	activeNum         int             // the number of clients currently connected to the server
	droppedNum        int             // the number of clients dropped by the server
	addConnChan       chan int        // send a signal when newly added a connection
	dropConnChan      chan int        // send a signal when newly dropped a connection
	countActiveChan   chan int        // request activeNum
	countDroppedChan  chan int        // request droppedNum
	countResponseChan chan int        // activeNum/droppedNum request result
	closeAcceptChan   chan int        // channel for closing accept routing
	closeMainChan     chan int        // channel for closing main routing
}

// New creates and returns (but does not start) a new KeyValueServer.
func New(store kvstore.KVStore) KeyValueServer {
	// TODO: implement this!
	return &keyValueServer{
		store,
		nil,
		0,
		0,
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan int),
	}
}

func (kvs *keyValueServer) Start(port int) error {
	// TODO: implement this!
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	kvs.listener = ln
	go count(kvs)
	go accept(kvs)
	go main(kvs)
	return nil
}

func (kvs *keyValueServer) Close() {
	// TODO: implement this!
	kvs.closeAcceptChan <- 0
	kvs.closeMainChan <- 0
}

func (kvs *keyValueServer) CountActive() int {
	// TODO: implement this!
	kvs.countActiveChan <- 1
	return <-kvs.countResponseChan
}

func (kvs *keyValueServer) CountDropped() int {
	// TODO: implement this!
	kvs.countDroppedChan <- 1
	return <-kvs.countResponseChan
}

// TODO: add additional methods/functions below!
func count(kvs *keyValueServer) {
	for {
		select {
		case <-kvs.addConnChan:
			kvs.activeNum++
		case <-kvs.dropConnChan:
			kvs.droppedNum++
			kvs.activeNum--
		case <-kvs.countActiveChan:
			kvs.countResponseChan <- kvs.activeNum
		case <-kvs.countDroppedChan:
			kvs.countResponseChan <- kvs.droppedNum
		}
	}
}

func accept(kvs *keyValueServer) {
	for {
		select {
		case <-kvs.closeAcceptChan:
			kvs.listener.Close()
			return
		default:
			conn, err := kvs.listener.Accept()
			if err != nil {
				// handle error
				continue
			}
			go handleConnection(kvs, conn)
		}
	}
}

func handleConnection(kvs *keyValueServer, connection net.Conn) {
	// TODO: implement this!
	kvs.addConnChan <- 1
}

func main(kvs *keyValueServer) {
	// TODO: implement this!
}
