package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const SensorName = "name"

var DefaultAddressMetrics = ""
var DefaultAddressSensor = ""

func main() {
	gin.SetMode(gin.ReleaseMode)

	routerSensor := gin.Default()
	routerSensor.Use(gin.BasicAuth(gin.Accounts{"sensor": getEnv("SENSOR_PASSWORD", "")}))
	routerSensor.GET("/sensors/:sensor_name", sensorAction)

	routerMetrics := gin.Default()
	routerMetrics.GET("/metrics", promHandler(promhttp.Handler()))
	go func() {
		log.Fatal(routerMetrics.Run(getEnv("ADDRESS_METRICS", DefaultAddressMetrics)))
	}()

	log.Fatal(routerSensor.Run(getEnv("ADDRESS_SENSOR", DefaultAddressSensor)))
}

func getEnv(name string, defaultValue string) string {
	value, exist := os.LookupEnv(name)
	if exist {
		return value
	}
	if defaultValue != "" {
		return defaultValue
	}
	log.Fatalf("%s environment variable is missing\n", name)
	return value
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
