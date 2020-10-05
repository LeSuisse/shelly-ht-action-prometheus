package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

const SensorName = "name"

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.GET("/sensors/:sensor_name", sensorAction)
	router.GET("/metrics", promHandler(promhttp.Handler()))

	log.Fatal(router.Run(":14380"))
}

var (
	temperatureSensor = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "shellyht_temperature_celsius",
		Help: "Current temperature in Celsius",
	}, []string{SensorName})
	humiditySensor = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "shellyht_humidity_percent",
		Help: "Current humidity level in %",
	}, []string{SensorName})
	lastSuccess = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "shellyht_last_report",
		Help: "Unixtime last time the sensor has reported a value",
	}, []string{SensorName})
)

type GetSensorActionParameters struct {
	Humidity    int     `form:"hum" binding:"required,min=0,max=100"`
	Temperature float64 `form:"temp" binding:"required"`
}

func sensorAction(c *gin.Context) {
	sensorName := c.Param("sensor_name")
	var actionParameters GetSensorActionParameters
	if err := c.BindQuery(&actionParameters); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Printf("%v%v\n", sensorName, actionParameters)
	labels := prometheus.Labels{SensorName: sensorName}
	temperatureSensor.With(labels).Set(actionParameters.Temperature)
	humiditySensor.With(labels).Set(float64(actionParameters.Humidity))
	lastSuccess.With(labels).SetToCurrentTime()
}

func promHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
