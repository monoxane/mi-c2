package api

import (
	"encoding/json"
	"fmt"
	"mi-c2/internal/controller"
	"mi-c2/internal/env"
	"mi-c2/internal/logging"
	"mi-c2/internal/model"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	agentWSConnections = make(map[*model.WebsocketConnection]uuid.UUID)
	mutex              sync.Mutex
	log                logging.Logger
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Listen(stop chan os.Signal) {
	log = logging.Log.With().Str("module", "api").Logger()
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "OK")
	})

	r.GET("v1/ca", handleClusterWS)

	go r.Run(fmt.Sprintf(":%s", env.PORT))

	<-stop
}

func handleClusterWS(c *gin.Context) {
	w := c.Writer
	r := c.Request

	socket, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "failed to upgrade ws connection", "error": err})

		return
	}

	connection := &model.WebsocketConnection{
		Socket: socket,
		Mux:    new(sync.Mutex),
	}

	mutex.Lock()
	agentWSConnections[connection] = uuid.Must(uuid.NewRandom())
	mutex.Unlock()

	defer connection.Socket.Close()
	for {
		_, message, err := connection.Socket.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				if websocket.IsCloseError(
					err,
					websocket.CloseNormalClosure,
					websocket.CloseNoStatusReceived,
					websocket.CloseGoingAway,
				) {
					mutex.Lock()
					delete(agentWSConnections, connection)
					mutex.Unlock()

					return
				}
			}

			controller.RemoveCluster(connection.Cluster)

			mutex.Lock()
			delete(agentWSConnections, connection)
			mutex.Unlock()

			break
		} else {
			var clusterUpdate model.Cluster
			err := json.Unmarshal(message, &clusterUpdate)
			if err != nil {
				log.Error().Err(err).Msgf("unable to unmarshal cluster update: %s", message)
			} else {
				connection.Cluster = clusterUpdate.Name
				clusterUpdate.Connection = connection
				controller.UpdateClusterStatus(&clusterUpdate)
			}
		}
	}
}
