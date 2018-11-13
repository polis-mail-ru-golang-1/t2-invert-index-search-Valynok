package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/t2-invert-index-search-Valynok/config"
	"github.com/t2-invert-index-search-Valynok/handlers"
	"github.com/t2-invert-index-search-Valynok/invertindex"
	mapUtils "github.com/t2-invert-index-search-Valynok/utils"
	"github.com/t2-invert-index-search-Valynok/view"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logCfg := zap.NewProductionConfig()
	logCfg.OutputPaths = []string{
		cfg.LogFileName,
	}

	var debugLevel zapcore.Level
	debugLevel.Set(cfg.LogLevel)

	logCfg.Level.SetLevel(debugLevel)
	logger, err := logCfg.Build()
	if err != nil {
		panic(err)
	}
	Logger = logger.Sugar()
	handlers.Logger = Logger
	index, fileNames := IndexDirectory(cfg.DirectoryPath)

	handlers.MainIndex = index
	handlers.FileNames = fileNames

	v, _ := view.New()
	h := handlers.New(v)

	http.HandleFunc("/", h.IndexHandler)
	http.HandleFunc("/search", h.SearchHandler)
	Logger.Infof("starting server at %s", cfg.Listen)
	http.ListenAndServe(cfg.Listen, nil)
}

func IndexDirectory(directory string) (invertindex.IndexType, []string) {
	fileNames := GetFileNames(directory)

	return IndexFiles(directory, fileNames), fileNames
}

func IndexFile(directoryPath string, fileName string) invertindex.IndexType {

	fileContent, err := ioutil.ReadFile(directoryPath + "/" + fileName)

	if err != nil {
		Logger.Error(err)
	}

	text := string(fileContent)
	return invertindex.GetIndex(text, fileName)
}

func GetFileNames(directoryRelativePath string) []string {

	files, err := ioutil.ReadDir(directoryRelativePath)
	if err != nil {
		Logger.Error(err)
	}

	fileNames := mapUtils.Map(files, func(fi os.FileInfo) string { return fi.Name() })
	fileNames = mapUtils.FilterFiles(fileNames, func(fn string) bool {
		return strings.HasSuffix(fn, ".txt")
	})
	return fileNames
}

func IndexFiles(filesDirectoryPath string, fileNames []string) invertindex.IndexType {
	filesIndex := make(invertindex.IndexType)
	fileIndexChannel := make(chan invertindex.IndexType, len(fileNames))
	for f := 0; f < len(fileNames); f++ {
		go func(fileName string) {
			fileIndex := IndexFile(filesDirectoryPath, fileName)
			Logger.Info("found ", len(fileIndex), "words in ", fileName)
			fileIndexChannel <- fileIndex

		}(fileNames[f])
		Logger.Debug("go routine for ", fileNames[f], " started")
	}

	for f := 0; f < len(fileNames); f++ {
		fileRes := <-fileIndexChannel
		Logger.Debug("got from pipe", len(fileRes))
		invertindex.MergeIndex(filesIndex, fileRes)
	}

	return filesIndex
}
