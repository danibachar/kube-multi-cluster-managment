package main

import (
	"net/http"
	"os"

	"github.com/danibachar/kube-multi-cluster-managment/server/providers/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func gcpDeployBroker(c echo.Context) error {
	log.Info("provider is gcp")

	config := models.GCPClusterConfig{}
	if err := c.Bind(&config); err != nil {
		log.Errorf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := GCPDeploySubmarinerBrokerOn(config); err != nil {
		log.Errorf("Failed join broker on cluster %s with error", config.ClusterName, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, config)
}

func gcpJoinBroker(c echo.Context) error {
	log.Info("provider is gcp")

	config := models.GCPClusterConfig{}
	if err := c.Bind(&config); err != nil {
		log.Errorf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if err := GCPJoinClusterToBroker(config); err != nil {
		log.Errorf("Failed deply broker on cluster %s with error", config.ClusterName, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, config)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/:name", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.DELETE("/:name", func(c echo.Context) error {
		provider := c.Param("name")
		switch provider {
		case "aws":
			break
		case "gcp":
			break
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "no supported provider")
	})

	e.PUT("set/:name", func(c echo.Context) error {
		defer c.Request().Body.Close()
		provider := c.Param("name")
		switch provider {
		case "aws":
			break
		case "gcp":
			return gcpDeployBroker(c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "no supported provider")
	})

	e.PUT("join/:name", func(c echo.Context) error {
		defer c.Request().Body.Close()
		provider := c.Param("name")
		switch provider {
		case "aws":
			break
		case "gcp":
			return gcpJoinBroker(c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "no supported provider")
	})

	httpPort := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(":" + httpPort))
}
