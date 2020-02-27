// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package server

import (
	"github.com/aquasecurity/fanal/analyzer"
	"github.com/aquasecurity/fanal/cache"
	"github.com/aquasecurity/trivy-db/pkg/db"
	db2 "github.com/aquasecurity/trivy/pkg/db"
	"github.com/aquasecurity/trivy/pkg/detector/library"
	"github.com/aquasecurity/trivy/pkg/detector/ospkg"
	"github.com/aquasecurity/trivy/pkg/github"
	"github.com/aquasecurity/trivy/pkg/indicator"
	library2 "github.com/aquasecurity/trivy/pkg/rpc/server/library"
	ospkg2 "github.com/aquasecurity/trivy/pkg/rpc/server/ospkg"
	"github.com/aquasecurity/trivy/pkg/scanner/local"
	"github.com/aquasecurity/trivy/pkg/vulnerability"
	"k8s.io/utils/clock"
)

// Injectors from inject.go:

func initializeScanServer(localLayerCache cache.LocalImageCache) *ScanServer {
	applier := analyzer.NewApplier(localLayerCache)
	detector := ospkg.Detector{}
	driverFactory := library.DriverFactory{}
	libraryDetector := library.NewDetector(driverFactory)
	scanner := local.NewScanner(applier, detector, libraryDetector)
	config := db.Config{}
	client := vulnerability.NewClient(config)
	scanServer := NewScanServer(scanner, client)
	return scanServer
}

func initializeOspkgServer() *ospkg2.Server {
	detector := ospkg.Detector{}
	config := db.Config{}
	client := vulnerability.NewClient(config)
	server := ospkg2.NewServer(detector, client)
	return server
}

func initializeLibServer() *library2.Server {
	driverFactory := library.DriverFactory{}
	detector := library.NewDetector(driverFactory)
	config := db.Config{}
	client := vulnerability.NewClient(config)
	server := library2.NewServer(detector, client)
	return server
}

func initializeDBWorker(quiet bool) dbWorker {
	config := db.Config{}
	client := github.NewClient()
	progressBar := indicator.NewProgressBar(quiet)
	realClock := clock.RealClock{}
	dbClient := db2.NewClient(config, client, progressBar, realClock)
	serverDbWorker := newDBWorker(dbClient)
	return serverDbWorker
}
