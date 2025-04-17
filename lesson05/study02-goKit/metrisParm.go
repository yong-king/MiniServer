package main

import (
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
)

var (
	// instrumentation
	fieldKeys = []string{"method", "error"}

	requestCount = kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "add_srv",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "add_srv",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	countResult = kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "add_srv",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

)