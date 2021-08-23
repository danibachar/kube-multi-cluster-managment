package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func createAWSCluster(c echo.Context) error {
	log.Info("provider is aws")

	config := AWSClusterConfig{}
	if err := c.Bind(&config); err != nil {
		log.Fatalf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	log.Info(config)
	if err := CreateAWSCluster(config); err != nil {
		log.Fatalf("Failed creating cluster %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	return c.JSON(http.StatusOK, config)
}

func createGCPCluster(c echo.Context) error {
	log.Info("provider is gcp")

	config := GCPClusterConfig{}
	if err := c.Bind(&config); err != nil {
		log.Fatalf("Failed reading the request body %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	log.Info(config)
	if err := CreateGCPCluster(config); err != nil {
		log.Fatalf("Failed creating cluster %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	return c.JSON(http.StatusOK, config)
}

func main() {

	e := echo.New()
	// e.Logger.SetLevel(log.DEBUG)
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.GET("/:name", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! 445 <3 11")
	})

	e.PUT("/:name", func(c echo.Context) error {
		defer c.Request().Body.Close()
		provider := c.Param("name")
		switch provider {
		case "aws":
			return createAWSCluster(c)
		case "gcp":
			return createGCPCluster(c)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "no supported provider")

	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
