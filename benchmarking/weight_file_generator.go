package benchmarking

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/shirou/gopsutil/v3/cpu"
)

type benchmarkInfo struct {
	Date          string
	Steps         int
	Repeat        int
	DbCache       int
	HeapPages     int
	HostName      string
	CpuName       string
	Gc            string
	TinyGoVersion string
	Target        string
	Summary       string
	PackageName   string
	FunctionName  string
}

func generateExtrinsicWeightFile(outputPath string, analysisResult analysis) error {
	data := struct {
		benchmarkInfo
		UsedComponents                                      []string
		BaseWeight, BaseReads, BaseWrites, MinExtrinsicTime uint64
		ComponentWeights, ComponentReads, ComponentWrites   []componentSlope
	}{}

	data.Date = time.Now().String()
	data.Steps = Config.Steps
	data.Repeat = Config.Repeat
	data.DbCache = Config.DbCache
	data.HeapPages = Config.HeapPages
	data.Gc = Config.GC
	data.TinyGoVersion = Config.TinyGoVersion
	data.Target = Config.Target

	data.Summary = analysisResult.String()

	paths := strings.Split(filepath.Dir(outputPath), "/")
	data.PackageName = paths[len(paths)-1]

	data.FunctionName = strcase.ToLowerCamel(strings.TrimSuffix(filepath.Base(outputPath), ".go")) // formats outputPath to weightFn name

	hostName, err := os.Hostname()
	if err != nil {
		return err
	}
	data.HostName = hostName

	if c, err := cpu.Info(); err == nil && len(c) > 0 {
		data.CpuName = fmt.Sprintf("%s(%d cores, %d mhz)", c[0].ModelName, c[0].Cores, int(c[0].Mhz))
	}

	data.UsedComponents = analysisResult.usedComponents
	data.BaseWeight = analysisResult.baseExtrinsicTime
	data.BaseReads = analysisResult.baseReads
	data.BaseWrites = analysisResult.baseWrites
	data.MinExtrinsicTime = analysisResult.minimumExtrinsicTime
	data.ComponentWeights = analysisResult.componentExtrinsicTimes
	data.ComponentReads = analysisResult.componentReads
	data.ComponentWrites = analysisResult.componentWrites

	// create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	// get current working directory
	_, cwd, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("error getting cwd from runtime")
	}

	// generate weight file
	templatePath := filepath.Join(filepath.Dir(cwd), "weight_file_extrinsic.go.tmpl")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(outputFile, data); err != nil {
		return err
	}

	return nil
}

func generateOverheadWeightFile(outputPath string, stats StatsResult) error {
	data := struct {
		benchmarkInfo
		BaseWeight uint64
	}{}

	data.Date = time.Now().String()
	data.Steps = Config.Steps
	data.Repeat = Config.Repeat
	data.DbCache = Config.DbCache
	data.HeapPages = Config.HeapPages
	data.Gc = Config.GC
	data.TinyGoVersion = Config.TinyGoVersion
	data.Target = Config.Target

	data.Summary = stats.String()

	paths := strings.Split(filepath.Dir(outputPath), "/")
	data.PackageName = paths[len(paths)-1]

	data.FunctionName = strcase.ToLowerCamel(strings.TrimSuffix(filepath.Base(outputPath), ".go")) // formats outputPath to weightFn name

	hostName, err := os.Hostname()
	if err != nil {
		return err
	}
	data.HostName = hostName

	if c, err := cpu.Info(); err == nil && len(c) > 0 {
		data.CpuName = fmt.Sprintf("%s(%d cores, %d mhz)", c[0].ModelName, c[0].Cores, int(c[0].Mhz))
	}

	data.BaseWeight = uint64(stats.Mean)

	// create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	// get current working directory
	_, cwd, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("error getting cwd from runtime")
	}

	// generate weight file
	templatePath := filepath.Join(filepath.Dir(cwd), "weight_file_overhead.go.tmpl")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(outputFile, data); err != nil {
		return err
	}

	return nil
}
