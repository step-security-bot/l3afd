package stats

import (
	"net/http"

	"tbd/go-shared/logs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	NFStartCount *prometheus.CounterVec
	NFStopCount  *prometheus.CounterVec
	NFRunning    *prometheus.GaugeVec
	NFStartTime	 *prometheus.GaugeVec
)

func SetupMetrics(hostname, daemonName, metricsAddr string) {

	nfStartCountVec := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: daemonName,
			Name:      "NFStartCount",
			Help:      "The count of network functions started",
		},
		[]string{"host", "network_function", "direction"},
	)

	NFStartCount = nfStartCountVec.MustCurryWith(prometheus.Labels{"host": hostname})

	nfStopCountVec := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: daemonName,
			Name:      "NFStopCount",
			Help:      "The count of network functions stopped",
		},
		[]string{"host", "network_function", "direction"},
	)

	NFStopCount = nfStopCountVec.MustCurryWith(prometheus.Labels{"host": hostname})

	nfRunningVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: daemonName,
			Name:      "NFRunning",
			Help:      "This value indicates network functions is running or not",
		},
		[]string{"host", "network_function", "direction"},
	)

	logs.IfWarningLogf(prometheus.Register(nfRunningVec), "Failed to register NFRunning metrics")

	NFRunning = nfRunningVec.MustCurryWith(prometheus.Labels{"host": hostname})

	nfStartTimeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: daemonName,
			Name:      "NFStartTime",
			Help:      "This value indicates start time of the network function since unix epoch in seconds",
		},
		[]string{"host", "network_function", "direction"},
	)

	logs.IfWarningLogf(prometheus.Register(nfStartTimeVec), "Failed to register NFStartTime metrics")

	NFStartTime = nfStartTimeVec.MustCurryWith(prometheus.Labels{"host": hostname})

	// Prometheus handler
	metricsHandler := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})

	// Adding web endpoint
	go func() {
		// Expose the registered metrics via HTTP.
		http.Handle("/metrics", metricsHandler)
		logs.IfFatalLogf(http.ListenAndServe(metricsAddr, nil), "Failed to launch prometheus metrics endpoint")
	}()
}

func Incr(counterVec *prometheus.CounterVec, networkFunction, direction string) {

	if counterVec == nil {
		logs.Warningf("Metrics: counter vector is nil and needs to be initialized")
		return
	}
	if nfCounter, err := counterVec.GetMetricWithLabelValues(networkFunction, direction); err == nil {
		nfCounter.Inc()
	}
}

func Set(value float64, gaugeVec *prometheus.GaugeVec, networkFunction, direction string) {

	if gaugeVec == nil {
		logs.Warningf("Metrics: gauge vector is nil and needs to be initialized")
		return
	}
	if nfGauge, err := gaugeVec.GetMetricWithLabelValues(networkFunction, direction); err == nil {
		nfGauge.Set(value)
	}
}
