package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/pkg/versioning"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	HealthCheckController struct {
		dbClient *mongo.Client
		group    *echo.Group
	}
	HealthCheckOutput struct {
		NoSQLDBStatus string `json:"no_sqldb_status,omitempty"`
		Version       string `json:"version,omitempty"`
	}
)

func NewHealthCheckController(
	dbClient *mongo.Client,
	e *echo.Echo, //nolint:varnamelen
) *HealthCheckController {
	return &HealthCheckController{
		dbClient: dbClient,
		group:    e.Group("/health_check"),
	}
}

func (c *HealthCheckController) Register() {
	c.group.GET("", c.health_check)
}

// @Summary     Checks service health
// @Description It does a ping to db and returns service version
// @Tags        healt check
// @Produce     json
// @Success     202            {object} HealthCheckOutput
// @Failure     424 {object} HealthCheckOutput
// @Failure     500
// @Router      /health_check [get]
func (c *HealthCheckController) health_check(ctx echo.Context) error {
	var res = HealthCheckOutput{
		NoSQLDBStatus: "connected",
		Version:       versioning.Version,
	}
	if err := c.dbClient.Ping(
		ctx.Request().Context(),
		nil,
	); err != nil {
		res.NoSQLDBStatus = "disconnected"
		return ctx.JSON(http.StatusFailedDependency, res)
	}
	return ctx.JSON(http.StatusOK, res)
}
