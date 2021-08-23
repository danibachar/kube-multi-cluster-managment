package main

import (
	"net/http"
	"os"

	"github.com/danibachar/kube-multi-cluster-managment/server/golang-pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func deleteGCPCluster(c echo.Context) error {
	log.Info("provider is gcp")

	config := models.GCPClusterConfig{}
	if err := c.Bind(&config); err != nil {
		log.Errorf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	log.Info(config)
	if err := GCPDeleteCluster(config); err != nil {
		log.Errorf("Failed deleting cluster %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, config)
}

func createGCPCluster(c echo.Context) error {
	log.Info("provider is gcp")

	config := models.GCPClusterConfig{}
	if err := c.Bind(&config); err != nil {
		log.Errorf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	log.Info(config)
	if err := GCPCreateOrUpdateCluster(config); err != nil {
		log.Errorf("Failed creating cluster %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, config)
}

func getGCPRegions(c echo.Context) error {
	log.Info("provider is gcp")

	project := c.QueryParam("project")
	err, regions := GetGCPRegions(project)
	if err != nil {
		log.Errorf("failed getting regions for projecy %s with error ", project, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, regions)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("regions/:name", func(c echo.Context) error {
		provider := c.Param("name")
		switch provider {
		case "aws":
			break
		case "gcp":
			return getGCPRegions(c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "no supported provider")
	})

	e.GET("/:name", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.DELETE("/:name", func(c echo.Context) error {
		provider := c.Param("name")
		switch provider {
		case "aws":
			break
		case "gcp":
			return deleteGCPCluster(c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "no supported provider")
	})

	e.PUT("/:name", func(c echo.Context) error {
		defer c.Request().Body.Close()
		provider := c.Param("name")
		switch provider {
		case "aws":
			break
		case "gcp":
			return createGCPCluster(c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "no supported provider")
	})

	httpPort := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(":" + httpPort))

	// Run one time firewall config
}
