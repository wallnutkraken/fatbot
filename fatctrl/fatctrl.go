package fatctrl

import (
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/wallnutkraken/fatbot/fatbrain"
	"github.com/wallnutkraken/fatbot/fatctrl/ctrltypes"
)

type HttpHandler struct {
	engine   *gin.Engine
	v1Router *gin.RouterGroup
	brain    *fatbrain.FatBotBrain
	addr     string
}

func New(addr string, brain *fatbrain.FatBotBrain) *HttpHandler {
	handlr := &HttpHandler{
		engine: gin.Default(),
		brain:  brain,
		addr:   addr,
	}
	handlr.v1Router = handlr.engine.Group("/v1")
	handlr.v1Router.PATCH("/training/status/:status", handlr.SetStatus)
	handlr.v1Router.GET("/training/status", handlr.GetStatus)

	return handlr
}

func (h *HttpHandler) SetStatus(ctx *gin.Context) {
	status := ctx.Params.ByName("status")
	switch status {
	case "train":
		var req ctrltypes.StartTrainingRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}
		go h.brain.TrainFor(time.Second * time.Duration(req.EndAfterSeconds))
		ctx.Status(http.StatusOK)
	case "stop":
		h.brain.StopTraining()
		ctx.Status(http.StatusOK)
	default:
		ctx.Status(http.StatusBadRequest)
	}
}

func (h *HttpHandler) GetStatus(ctx *gin.Context) {
	var netStatus string
	if h.brain.IsTraining() {
		netStatus = "training"
	} else {
		netStatus = "idle"
	}
	ctx.JSON(http.StatusOK, ctrltypes.StatusResponse{
		Network: netStatus,
	})
}

func (h *HttpHandler) Start() error {
	return h.engine.Run(h.addr)
}
