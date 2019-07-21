package main

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"market-patterns/model"
	"market-patterns/model/report"
	"net/http"
	_ "net/http/pprof"
	"sort"
)

func startProfile() {
	log.Info("Starting profile server...")
	log.Fatal(http.ListenAndServe("localhost:6060", nil))
}

func start() {

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./ui/build", true)))

	apiLatest := router.Group("/api/latest")
	apiLatest.GET("/predict/:id", handlePredict)
	apiLatest.GET("/tickers", handleTickers)
	apiLatest.GET("/series/:id", handleSeries)
	apiLatest.GET("/symbols", handleSymbols)
	apiLatest.GET("/edge-probabilities/:density", handleEdgeProbabilities)
	apiLatest.GET("/graph/pattern-density/:id", handlePatternDensity)
	apiLatest.GET("/graph/stock/:id", handlePeriodCloseSeries)

	log.Info("market-pattern server listening...")

	log.Fatal(router.Run(":7666"))
}

func handlePredict(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "undefined" {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("no 'id' path parameter defined"))
	}

	prediction, err := predict(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, prediction)
}

func handlePeriodCloseSeries(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "undefined" {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("no 'id' path parameter defined"))
	}

	series, err := Repos.GraphController.FindPeriodCloseSeries(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, series)
}

func handleTickers(ctx *gin.Context) {
	tickerNames := report.TickerNames{Names: Repos.TickerRepo.FindSymbolsAndCompany()}
	ctx.JSON(http.StatusOK, tickerNames)
}

func handleSymbols(ctx *gin.Context) {
	tickerNames := report.SymbolNames{Names: Repos.TickerRepo.FindSymbols()}
	ctx.JSON(http.StatusOK, tickerNames)
}

func handleSeries(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "undefined" {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("no 'id' path parameter defined"))
	}

	series, err := Repos.SeriesRepo.FindBySymbol(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, "problem access Series data"))
	}

	seriesAry := make([]string, len(series))
	for idx, val := range series {
		seriesAry[idx] = val.Name
	}

	sort.Strings(seriesAry)
	data := report.SeriesNames{Names: seriesAry}
	ctx.JSON(http.StatusOK, data)
}

func handlePatternDensity(ctx *gin.Context) {

	id := ctx.Param("id")
	if id == "undefined" {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("no 'id' path parameter defined"))
	}

	data, err := Repos.GraphController.FindPatternDensities(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusNotFound, err)
	}
	ctx.JSON(http.StatusOK, data)
}

func handleEdgeProbabilities(ctx *gin.Context) {

	densityStr := ctx.Param("density")
	if densityStr == "undefined" {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("no 'density' path parameter defined"))
	}

	density, err := model.PatternDensityFromString(densityStr)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid 'density' path parameter"))
	}

	upHighProd, err := Repos.PatternRepo.FindHighestUpProbability(density)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError,
			errors.Wrap(err, "problem generating highest up prob"))
	}
	downHighProd, err := Repos.PatternRepo.FindHighestDownProbability(density)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError,
			errors.Wrap(err, "problem generating highest down prob"))
	}
	ncHighProd, err := Repos.PatternRepo.FindHighestNoChangeProbability(density)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError,
			errors.Wrap(err, "problem generating highest nochange prob"))
	}
	upLowProd, err := Repos.PatternRepo.FindLowestUpProbability(density)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError,
			errors.Wrap(err, "problem generating lowest up prob"))
	}
	downLowProd, err := Repos.PatternRepo.FindLowestDownProbability(density)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError,
			errors.Wrap(err, "problem generating lowest down prob"))
	}
	ncLowProd, err := Repos.PatternRepo.FindLowestNoChangeProbability(density)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError,
			errors.Wrap(err, "problem generating lowest nochange prob"))
	}

	probs := report.ProbabilityEdges{BestDownHigh: downHighProd, BestDownLow: downLowProd, BestNoChangeHigh: ncHighProd,
		BestNoChangeLow: ncLowProd, BestUpHigh: upHighProd, BestUpLow: upLowProd}

	ctx.JSON(http.StatusOK, probs)
}
