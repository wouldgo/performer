package prom

import (
	"net/http"
	"performer/client"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	reg *prometheus.Registry

	download prometheus.Gauge
	upload   prometheus.Gauge

	dataChan chan *client.BitsPerSecond
}

func NewMetrics(dataChan chan *client.BitsPerSecond) *Metrics {
	reg := prometheus.NewRegistry()

	m := &Metrics{
		reg: reg,
		download: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "download_speed",
			Help: "download speed",
		}),
		upload: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "upload_speed",
			Help: "upload speed",
		}),
		dataChan: dataChan,
	}
	reg.MustRegister(m.download)
	reg.MustRegister(m.upload)
	return m
}

func (m *Metrics) handler() http.Handler {
	return promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{
		Registry: m.reg,
	})
}

func (m *Metrics) setUpload(newValue float64) {
	m.upload.Set(newValue)
}

func (m *Metrics) setDownload(newValue float64) {
	m.download.Set(newValue)
}

func (m *Metrics) Handle() error {
	go func() {
		for aReport := range m.dataChan {

			m.setDownload(aReport.Download)
			m.setUpload(aReport.Upload)
		}
	}()

	http.Handle("/metrics", m.handler())

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}
	return nil
}
