package main

import (
	"context"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	v1alpha1 "sigs.k8s.io/mcs-api/pkg/apis/v1alpha1"
	mcsClientset "sigs.k8s.io/mcs-api/pkg/client/clientset/versioned"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func prepareClients() (kubernetes.Interface, mcsClientset.Interface, error) {
	config, err := restclient.InClusterConfig()
	if err != nil {
		return nil, err
	}
	mcsClientSet, mcsError := mcsClientset.NewForConfig(config)
	if mcsError != nil {
		return nil, nil, mcsError
	}
	kubeClientSet, kubeError := kubernetes.NewForConfig(config)
	if kubeError != nil {
		return nil, mcsClientSet, kubeError
	}
	return kubeClientSet, mcsClientSet, nil
}

func export(services Exports) error {
	kubeClientSet, mcsClientSet, err := prepareClients()
	if err != nil {
		return err
	}
	// validate service exists locally
	_, err = kubeClientSet.CoreV1().Services(serviceNamespace).Get(context.TODO(), svcName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	// Create the export
	// mcsServiceExport := &v1alpha1.ServiceExport{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:        svcName,
	// 		Namespace:   serviceNamespace,
	// 		Annotations: withAnnotations,
	// 	},
	// }
	for _, svcExport := range services.servicesToExport {
		_, err := mcsClientSet.MulticlusterV1alpha1.ServiceExport(svcExport.metadata.Namespace).Create(context.TODO(), svcExport, metav1.CreateOption{})
		if err != nil {
			return err
		}
	}
	return nil
}

type Exports struct {
	servicesToExport []v1alpha1.ServiceExport `json:"servicesToExport,omitempty"`
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.PUT("/export", func(c echo.Context) error {
		log.Info("1")
		exports := new(Exports)
		if err := c.Bind(exports); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err := export(*exports)
		log.Info("2")
		if err != nil {
			log.Info("3")
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		log.Info("4")
		return c.JSON(http.StatusOK, exports)
	})

	httpPort := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(":" + httpPort))
}
