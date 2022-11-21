package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/pkg/versioning"
	"gorm.io/gorm"
)

const (
	connected      = "connected"
	dbDisconnected = "disconnected"
)

type (
	HealthCheckController struct {
		conn  *gorm.DB
		group *echo.Group
	}
	HealthCheckOutput struct {
		DBStatus string `json:"db_status,omitempty"`
		Version  string `json:"version,omitempty"`
	}
)

func NewHealthCheckController(
	conn *gorm.DB,
	e *echo.Echo, //nolint:varnamelen
) *HealthCheckController {
	return &HealthCheckController{
		conn:  conn,
		group: e.Group("/health_check"),
	}
}

func (c *HealthCheckController) Register() {
	c.group.GET("", c.health_check)
}

// @Summary     Checks service health
// @Description It does a ping to db and returns service version
// @Tags        healt check
// @Produce     json
// @Success     202 {object} HealthCheckOutput
// @Failure     424 {object} HealthCheckOutput
// @Failure     500
// @Router      /health_check [get]
func (c *HealthCheckController) health_check(ctx echo.Context) error {
	var res = HealthCheckOutput{
		DBStatus: connected,
		Version:  versioning.Version,
	}
	db, err := c.conn.DB()
	if err != nil {
		res.DBStatus = dbDisconnected
		return ctx.JSON(http.StatusFailedDependency, res)
	}
	if err := db.Ping(); err != nil {
		res.DBStatus = dbDisconnected
		return ctx.JSON(http.StatusFailedDependency, res)
	}
	return ctx.JSON(http.StatusOK, res)
}
