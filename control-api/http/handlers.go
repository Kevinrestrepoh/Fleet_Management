package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Kevinrestrepoh/control-api/grpc"
	"github.com/Kevinrestrepoh/control-api/types"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	client *grpc.Client
}

func NewHandler(client *grpc.Client) *Handler {
	return &Handler{client}
}

func (h *Handler) SendCommand(c echo.Context) error {
	vehicleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid vehicle id",
		})
	}

	var req types.CommandRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid json body",
		})
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	err = h.client.SendCommand(uint32(vehicleID), req, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
	})
}

func (h *Handler) MetricsSSE(c echo.Context) error {
	stream, err := h.client.Stream(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	res := c.Response()
	res.Header().Set(echo.HeaderContentType, "text/event-stream")
	res.Header().Set(echo.HeaderCacheControl, "no-cache")
	res.Header().Set(echo.HeaderConnection, "keep-alive")
	res.WriteHeader(http.StatusOK)

	for {
		metrics, err := stream.Recv()
		if err != nil {
			return nil
		}

		payload, err := json.Marshal(metrics)
		if err != nil {
			continue
		}

		res.Write([]byte("data: "))
		res.Write(payload)
		res.Write([]byte("\n\n"))

		if flusher, ok := res.Writer.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}
