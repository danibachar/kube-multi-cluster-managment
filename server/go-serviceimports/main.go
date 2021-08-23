package main

import (
	"context"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	restclient "k8s.io/client-go/rest"
	v1alpha1 "sigs.k8s.io/mcs-api/pkg/apis/v1alpha1"
	mcsClientset "sigs.k8s.io/mcs-api/pkg/client/clientset/versioned"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func prepareClient() (mcsClientset.Interface, error) {
	config, err := restclient.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return mcsClientset.NewForConfig(config)
}

func getAllImportedServicesIn(namespace string) (*v1alpha1.ServiceImportList, error) {
	clientSet, err := prepareClient()
	if err != nil {
		return nil, err
	}
	return clientSet.MulticlusterV1alpha1().ServiceImports(namespace).List(context.TODO(), metav1.ListOptions{})
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		log.Info("1")
		imports, err := getAllImportedServicesIn(metav1.NamespaceAll)
		log.Info("2")
		if err != nil {
			log.Info("3")
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		log.Info("4")
		return c.JSON(http.StatusOK, imports)
	})

	e.GET("/:namespace", func(c echo.Context) error {
		// imports, err := getAllImportedServicesIn(c.Param("namespace"))
		// if err != nil {
		// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
		// }
		// return c.JSON(http.StatusOK, imports)
		services := "hello"
		return c.JSON(http.StatusOK, services)
	})

	httpPort := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(":" + httpPort))
}
