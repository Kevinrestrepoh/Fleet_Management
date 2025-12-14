package http

import (
	"net/http"
	"strconv"

	"github.com/Kevinrestrepoh/control-api/grpc"
	"github.com/Kevinrestrepoh/control-api/types"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	ingestion *grpc.IngestionClient
}

func NewHandler(ingestion *grpc.IngestionClient) *Handler {
	return &Handler{ingestion: ingestion}
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

	err = h.ingestion.SendCommand(uint32(vehicleID), req, c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
	})
}
