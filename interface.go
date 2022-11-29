package ws

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/umpc/go-sortedmap"
	"github.com/umpc/go-sortedmap/asc"
)

func InitWebSocket(h func(string) bool, addr string) {
	if wsserverInternal == nil {
		wsserverInternal = &wsserver{
			RWMutex:       sync.RWMutex{},
			internalWSMap: map[string]*websocket.Conn{},
		}
	}

	wsserverInternal.wsmapTTL = sortedmap.New(1, asc.Time)

	wsserverInternal.upgrader = websocket.Upgrader{
		HandshakeTimeout: time.Second * handshakeTimeout,
		ReadBufferSize:   readBufferSize,
		WriteBufferSize:  writeBufferSize,
		WriteBufferPool:  nil,
		Subprotocols:     []string{},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		},

		CheckOrigin: func(r *http.Request) bool {
			return true
		},

		EnableCompression: EnableCompression,
	}
	gin.SetMode(gin.ReleaseMode)

	wsserverInternal.ginEngine = gin.New()

	wsserverInternal.ginEngine.GET(defaultPath, wsserverInternal.getwshandler(h))
	wsserverInternal.ginEngine.DELETE(defaultPath, wsserverInternal.delwshandler(h))

	// Wait for gin server to initialize and run in background
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup, addr string) {
		wg.Done()
		wsserverInternal.ginEngine.Run(addr)

	}(wg, addr)
	wg.Wait()

	wsserverInternal.checkTTLofRecords()

}

func GetWebSocketSession(sessionID string) (ok bool, ws *websocket.Conn) {
	return wsserverInternal.getWebSocketSession(sessionID)
}

func NoOfWebSocketClients() int {
	return wsserverInternal.len()
}
