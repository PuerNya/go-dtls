package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"html"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type metric struct {
	unit   string
	values []float64
}

type benchmark struct {
	name    string
	samples int
	metrics []*metric
	byUnit  map[string]*metric
}

type report struct {
	goVersion  string
	goos       string
	goarch     string
	cpu        string
	wolfSSL    string
	benchmarks []*benchmark
	byName     map[string]*benchmark
}

func main() {
	input := flag.String("input", "", "go test benchmark output")
	output := flag.String("output", "", "Markdown report path")
	commit := flag.String("commit", "", "tested commit")
	wolfSSL := flag.String("wolfssl", "", "wolfSSL commit and build")
	generated := flag.String("generated", time.Now().UTC().Format(time.RFC3339), "generation time")
	flag.Parse()
	if *input == "" || *output == "" {
		flag.Usage()
		os.Exit(2)
	}

	in, err := os.Open(*input) // #nosec G304 -- path is an explicit local command argument.
	if err != nil {
		fatal(err)
	}
	result, err := parseReport(in)
	closeErr := in.Close()
	if err != nil {
		fatal(err)
	}
	if closeErr != nil {
		fatal(closeErr)
	}

	out, err := os.Create(*output) // #nosec G304 -- path is an explicit local command argument.
	if err != nil {
		fatal(err)
	}
	err = writeReport(out, result, *commit, *wolfSSL, *generated)
	closeErr = out.Close()
	if err != nil {
		fatal(err)
	}
	if closeErr != nil {
		fatal(closeErr)
	}
}

func fatal(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func parseReport(reader io.Reader) (*report, error) {
	result := &report{byName: make(map[string]*benchmark)}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "go version ") && result.goVersion == "":
			result.goVersion = line
		case strings.HasPrefix(line, "goos: ") && result.goos == "":
			result.goos = strings.TrimSpace(strings.TrimPrefix(line, "goos: "))
		case strings.HasPrefix(line, "goarch: ") && result.goarch == "":
			result.goarch = strings.TrimSpace(strings.TrimPrefix(line, "goarch: "))
		case strings.HasPrefix(line, "cpu: ") && result.cpu == "":
			result.cpu = strings.TrimSpace(strings.TrimPrefix(line, "cpu: "))
		case strings.HasPrefix(line, "wolfssl: ") && result.wolfSSL == "":
			result.wolfSSL = strings.TrimSpace(strings.TrimPrefix(line, "wolfssl: "))
		case strings.HasPrefix(line, "Benchmark"):
			if err := result.addBenchmarkLine(lineNumber, line); err != nil {
				return nil, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read benchmark output: %w", err)
	}
	if len(result.benchmarks) == 0 {
		return nil, errors.New("benchmark output contains no results")
	}
	return result, nil
}

func (r *report) addBenchmarkLine(lineNumber int, line string) error {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return fmt.Errorf("line %d: malformed benchmark result %q", lineNumber, line)
	}
	if _, err := strconv.ParseUint(fields[1], 10, 64); err != nil {
		return fmt.Errorf("line %d: parse iteration count: %w", lineNumber, err)
	}
	if len(fields[2:])%2 != 0 {
		return fmt.Errorf("line %d: metric value/unit pairs are incomplete", lineNumber)
	}

	item := r.byName[fields[0]]
	if item == nil {
		item = &benchmark{name: fields[0], byUnit: make(map[string]*metric)}
		r.byName[item.name] = item
		r.benchmarks = append(r.benchmarks, item)
	}
	item.samples++
	for index := 2; index < len(fields); index += 2 {
		value, err := strconv.ParseFloat(fields[index], 64)
		if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
			return fmt.Errorf("line %d: invalid metric value %q", lineNumber, fields[index])
		}
		unit := fields[index+1]
		entry := item.byUnit[unit]
		if entry == nil {
			entry = &metric{unit: unit}
			item.byUnit[unit] = entry
			item.metrics = append(item.metrics, entry)
		}
		entry.values = append(entry.values, value)
	}
	return nil
}

func writeReport(writer io.Writer, result *report, commit, wolfSSL, generated string) error {
	if wolfSSL == "" {
		wolfSSL = result.wolfSSL
	}
	if _, err := fmt.Fprintln(writer, "# Automated benchmark results"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}
	platform := result.goos + "/" + result.goarch
	if result.goos == "" && result.goarch == "" {
		platform = ""
	}
	if platform != "" && result.cpu != "" {
		platform += ", "
	}
	platform += result.cpu
	metadata := [][2]string{
		{"Commit", commit},
		{"Generated", generated},
		{"Go", result.goVersion},
		{"Platform", platform},
		{"wolfSSL", wolfSSL},
	}
	for _, item := range metadata {
		if item[1] == "" {
			continue
		}
		if _, err := fmt.Fprintf(writer, "- %s: `%s`\n", item[0], markdownText(item[1])); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer, "Values are medians of the samples emitted by the final benchmark run. Units are reported by Go's benchmark harness."); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer, "| Benchmark | Samples | Median metrics |\n| --- | ---: | --- |"); err != nil {
		return err
	}
	for _, item := range result.benchmarks {
		metrics := make([]string, 0, len(item.metrics))
		for _, entry := range item.metrics {
			metrics = append(metrics, formatNumber(median(entry.values))+" "+markdownText(entry.unit))
		}
		if _, err := fmt.Fprintf(writer, "| `%s` | %d | %s |\n", markdownText(item.name), item.samples, strings.Join(metrics, "<br>")); err != nil {
			return err
		}
	}
	return nil
}

func median(values []float64) float64 {
	ordered := append([]float64(nil), values...)
	sort.Float64s(ordered)
	middle := len(ordered) / 2
	if len(ordered)%2 != 0 {
		return ordered[middle]
	}
	return (ordered[middle-1] + ordered[middle]) / 2
}

func formatNumber(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func markdownText(value string) string {
	value = html.EscapeString(value)
	value = strings.ReplaceAll(value, "|", "&#124;")
	return strings.ReplaceAll(value, "`", "&#96;")
}
