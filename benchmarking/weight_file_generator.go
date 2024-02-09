package benchmarking

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/shirou/gopsutil/v3/cpu"
)

func generateWeightFile(template *weightFileTemplate, outputPath, weightSummary string, refTime, reads, writes uint64) error {
	if err := template.SetWeightValues(refTime, reads, writes); err != nil {
		return err
	}

	weightFn := strcase.ToLowerCamel(strings.TrimSuffix(filepath.Base(outputPath), ".go")) // formats outputPath to weightFn name

	if err := template.SetWeightFnName(weightFn); err != nil {
		return err
	}

	alertGeneratedFile := "THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE"
	hostName, _ := os.Hostname()
	cpuInfo := ""
	if c, err := cpu.Info(); err == nil && len(c) > 0 {
		cpuInfo = fmt.Sprintf("%s(%d cores, %d mhz)", c[0].ModelName, c[0].Cores, int(c[0].Mhz))
	}
	info := fmt.Sprintf(
		"%s\nDATE: %s, STEPS: %d, REPEAT: %d, DBCACHE: %d, HEAPPAGES: %d, HOSTNAME: %s, CPU: %s, GC: %s, TINYGO VERSION: %s, TARGET: %s\nSummary:\n%s",
		alertGeneratedFile, time.Now(), Config.Steps, Config.Repeat, Config.DbCache, Config.HeapPages, hostName, cpuInfo, Config.GC, Config.TinyGoVersion, Config.Target, weightSummary,
	)

	if err := template.SetInfoComment(info); err != nil {
		return err
	}

	paths := strings.Split(filepath.Dir(outputPath), "/")
	packageName := paths[len(paths)-1]
	if err := template.SetPackageName(packageName); err != nil {
		return err
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	if err := template.WriteGeneratedFile(outputFile); err != nil {
		return err
	}

	return nil
}
