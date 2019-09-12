package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go-market-patterns/config"
	"go-market-patterns/model/core"
	"go-market-patterns/model/graph"
	"go-market-patterns/model/report"
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

func startProfile() {
	log.Info("Starting profile server...")
	log.Fatal(http.ListenAndServe("localhost:6060", nil))
}

func start(conf *config.AppConfig) {

	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./ui/build", true)))

	apiLatest := router.Group("/api/latest")
	apiLatest.GET("/predict/:symbol", handlePredict)
	apiLatest.GET("/tickers", handleTickers)
	apiLatest.GET("/series/:symbol", handleSeries)
	apiLatest.GET("/symbols", handleSymbols)
	apiLatest.GET("/edge-probabilities/:density", handleEdgeProbabilities)
	apiLatest.GET("/graph/pattern-density/:symbol", handlePatternDensity)
	apiLatest.GET("/graph/stock/:symbol", handlePeriodCloseSeries)

	log.Info("market-pattern server listening...")

	log.Fatal(router.Run(conf.Runtime.HttpServerUrl))
}

func handlePredict(ctx *gin.Context) {

	id := ctx.Param("symbol")
	if id == "undefined" || id == "" {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("no 'symbol' path parameter defined"))
	}
	lengthParam := ctx.Query("length")
	length, err := strconv.Atoi(lengthParam)
	if err != nil || length <= 0 {
		_ = ctx.AbortWithError(http.StatusBadRequest,
			fmt.Errorf("length parameter '%v' is not a valid parameter", lengthParam))
	}

	prediction, err := predictOne(id, length)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, prediction)
}

func handlePeriodCloseSeries(ctx *gin.Context) {

	id := ctx.Param("symbol")
	if id == "undefined" || id == "" {
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("no 'symbol' path parameter defined"))
	}

	series, err := Repos.GraphController.FindPeriodCloseSeries(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}
	ctx.JSON(http.StatusOK, series)
}

func handleTickers(ctx *gin.Context) {
	tickerNames := report.Tickers{Tickers: Repos.TickerRepo.FindSymbolCompanySliceSortAsc()}
	ctx.JSON(http.StatusOK, tickerNames)
}

func handleSymbols(ctx *gin.Context) {
	tickerNames := report.TickerSymbolNames{Names: Repos.TickerRepo.FindSymbols()}
	ctx.JSON(http.StatusOK, tickerNames)
}

func handleSeries(ctx *gin.Context) {

	id := ctx.Param("symbol")
	if id == "undefined" || id == "" {
		_ = ctx.AbortWithError(http.StatusBadRequest,
			errors.New("no 'symbol' path parameter defined"))
	}

	data := Repos.SeriesRepo.FindNameLengthSliceBySymbol(id)
	if data == nil {
		_ = ctx.AbortWithError(http.StatusNotExtended,
			fmt.Errorf("series with symbol '%v' not found", id))
	}
	ctx.JSON(http.StatusOK, report.Series{Series: data})
}

func handlePatternDensity(ctx *gin.Context) {

	id := ctx.Param("symbol")
	if id == "undefined" || id == "" {
		_ = ctx.AbortWithError(http.StatusBadRequest,
			errors.New("no 'symbol' path parameter defined"))
	}

	lengthParam := ctx.Query("length")
	length, err := strconv.Atoi(lengthParam)
	if err != nil || length <= 0 {
		_ = ctx.AbortWithError(http.StatusBadRequest,
			fmt.Errorf("length parameter '%v' is not a valid parameter", lengthParam))
	}
	data, err := Repos.GraphController.FindPatternDensities(id, length)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusNotFound, err)
	}

	companyName, err := Repos.TickerRepo.FindOneCompanyNameBySymbol(id)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError,
			errors.Wrapf(err, "problem getting company name for symbol %v", id))
	}
	graphData := graph.PatternDensityGraph{Symbol: id, CompanyName: companyName, Graph: data}

	ctx.JSON(http.StatusOK, graphData)
}

func handleEdgeProbabilities(ctx *gin.Context) {

	densityStr := ctx.Param("density")
	if densityStr == "undefined" || densityStr == "" {
		_ = ctx.AbortWithError(http.StatusBadRequest,
			errors.New("no 'density' path parameter defined"))
	}

	density, err := core.PatternDensityFromString(densityStr)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest,
			errors.New("invalid 'density' path parameter"))
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
