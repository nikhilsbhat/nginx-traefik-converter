// Package render provides human-friendly output renderers for migration reports.
// It can print per-Ingress and global summaries in either table format (using
// tablewriter) or plain text format, depending on configuration.
package render

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/nikhilsbhat/nginx-traefik-converter/pkg/configs"
	"github.com/olekukonko/tablewriter"
)

// SummaryCounts represents aggregated counts of annotation conversion outcomes.
// It is used for both per-Ingress summaries and global summaries.
type SummaryCounts struct {
	// Converted is the number of annotations successfully converted.
	Converted int `yaml:"converted,omitempty" json:"converted,omitempty"`

	// Warnings is the number of annotations converted with warnings.
	Warnings int `yaml:"warnings,omitempty"  json:"warnings,omitempty"`

	// Skipped is the number of annotations that could not be converted and
	// require manual migration.
	Skipped int `yaml:"skipped,omitempty"   json:"skipped,omitempty"`

	// Ignored is the number of annotations that were intentionally ignored.
	Ignored int `yaml:"ignored,omitempty"   json:"ignored,omitempty"`
}

// Config controls how reports are rendered.
type Config struct {
	// Table determines whether output should be rendered as tables (true)
	// or as plain text (false).
	Table bool `yaml:"table,omitempty" json:"table,omitempty"`
}

// statusLabel maps annotation statuses to human-readable labels for display.
var statusLabel = map[configs.AnnotationStatus]string{
	configs.AnnotationConverted: "Converted",
	configs.AnnotationWarned:    "Warning",
	configs.AnnotationSkipped:   "Skipped",
	configs.AnnotationIgnored:   "Ignored",
}

const fixedStringLength = 80

// ---------------- Public API ----------------

// PrintIngressSummary renders the migration report for a single Ingress.
// The output format (table or text) is selected based on the Config.
func (cfg *Config) PrintIngressSummary(ingressReport configs.IngressReport) error {
	if cfg.Table {
		return cfg.printIngressReportTable(ingressReport)
	}

	cfg.printIngressReport(ingressReport)

	return nil
}

// PrintGlobalSummary renders the aggregated migration summary across all Ingresses.
// The output format (table or text) is selected based on the Config.
func (cfg *Config) PrintGlobalSummary(globalReport configs.GlobalReport) error {
	if cfg.Table {
		return cfg.printGlobalSummaryTable(globalReport)
	}

	cfg.printGlobalSummary(globalReport)

	return nil
}

// New returns a new Config with default settings.
func New() *Config {
	return &Config{}
}

// ---------------- Separators ----------------

func printSectionSeparator(title string) {
	line := strings.Repeat("=", fixedStringLength)
	fmt.Println(line)
	fmt.Println(color.HiCyanString(title))
	fmt.Println(line)
	fmt.Println()
}

func printSubSectionSeparator(title string) {
	line := strings.Repeat("-", fixedStringLength)
	fmt.Println(line)
	fmt.Println(title)
	fmt.Println(line)
	fmt.Println()
}

// ---------------- Table Renderers ----------------

// printIngressReportTable renders a single Ingress report in table format,
// including a detailed per-annotation table and a summary table.
func (cfg *Config) printIngressReportTable(ingressReport configs.IngressReport) error {
	printSectionSeparator(fmt.Sprintf("INGRESS: %s/%s", ingressReport.Namespace, ingressReport.Name))

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Annotation", "Status", "Message"})

	rows := make([][]string, 0, len(ingressReport.Entries))

	for _, entries := range ingressReport.Entries {
		msg := entries.Message
		if msg == "" {
			msg = "-"
		}

		rows = append(rows, []string{entries.Name, statusLabelColored(entries.Status), msg})
	}

	if err := table.Bulk(rows); err != nil {
		return err
	}

	if err := table.Render(); err != nil {
		return err
	}

	// Render per-Ingress summary table.
	printSubSectionSeparator("SUMMARY")

	return renderSummaryTable(summarizeIngress(ingressReport))
}

// printGlobalSummaryTable renders the global summary across all Ingresses
// in table format.
func (cfg *Config) printGlobalSummaryTable(globalReport configs.GlobalReport) error {
	printSectionSeparator("GLOBAL SUMMARY")

	return renderSummaryTable(summarizeGlobal(globalReport))
}

// renderSummaryTable renders a generic summary table given a title and summary counts.
func renderSummaryTable(summaryCounts SummaryCounts) error {
	summary := tablewriter.NewWriter(os.Stdout)
	summary.Header([]string{"Metric", "Count"})

	rows := [][]string{
		{"Converted", color.HiGreenString(strconv.Itoa(summaryCounts.Converted))},
		{"Warnings", color.HiYellowString(strconv.Itoa(summaryCounts.Warnings))},
		{"Skipped", color.HiRedString(strconv.Itoa(summaryCounts.Skipped))},
		{"Ignored", color.HiBlueString(strconv.Itoa(summaryCounts.Ignored))},
		{"Result", resultLabel(summaryCounts)},
	}

	if err := summary.Bulk(rows); err != nil {
		return err
	}

	return summary.Render()
}

// ---------------- Text Renderers ----------------

// printIngressReport renders a single Ingress report in plain text format.
func (cfg *Config) printIngressReport(ingressReport configs.IngressReport) {
	printSectionSeparator(fmt.Sprintf("INGRESS: %s/%s", ingressReport.Namespace, ingressReport.Name))

	for _, entries := range ingressReport.Entries {
		switch entries.Status {
		case configs.AnnotationConverted:
			fmt.Printf("  ✅ %s\n", entries.Name)
		case configs.AnnotationWarned:
			fmt.Printf("  ⚠️  %s\n      → %s\n", entries.Name, entries.Message)
		case configs.AnnotationSkipped:
			fmt.Printf("  ❌ %s\n      → %s\n", entries.Name, entries.Message)
		case configs.AnnotationIgnored:
			fmt.Printf("  ℹ️  %s\n", entries.Name)
		}
	}

	printSubSectionSeparator("SUMMARY")
	printSummaryText(
		fmt.Sprintf("Summary for %s/%s", ingressReport.Namespace, ingressReport.Name),
		summarizeIngress(ingressReport),
	)
}

// printGlobalSummary renders the aggregated global summary in plain text format.
func (cfg *Config) printGlobalSummary(globalReport configs.GlobalReport) {
	printSectionSeparator("GLOBAL SUMMARY")
	printSummaryText("Global Summary", summarizeGlobal(globalReport))
}

// printSummaryText prints a human-readable summary block in plain text.
func printSummaryText(title string, summaryCounts SummaryCounts) {
	fmt.Printf("%s\n", color.HiCyanString(title))
	fmt.Printf("Converted: %s\n", color.HiGreenString(strconv.Itoa(summaryCounts.Converted)))
	fmt.Printf("Warnings:  %s\n", color.HiYellowString(strconv.Itoa(summaryCounts.Warnings)))
	fmt.Printf("Skipped:   %s\n", color.HiRedString(strconv.Itoa(summaryCounts.Skipped)))
	fmt.Printf("Ignored:   %s\n", color.HiBlueString(strconv.Itoa(summaryCounts.Ignored)))
	fmt.Printf("Result:    %s\n\n", resultLabel(summaryCounts))
}

// ---------------- Helpers ----------------

// resultLabel returns a human-readable overall result string based on summary counts.
func resultLabel(summaryCounts SummaryCounts) string {
	if summaryCounts.Skipped > 0 {
		return color.HiRedString("Manual action required")
	}

	if summaryCounts.Warnings > 0 {
		return color.HiYellowString("Review recommended")
	}

	return color.HiGreenString("Clean migration")
}

// summarizeIngress computes summary counts for a single Ingress report.
func summarizeIngress(ingressReport configs.IngressReport) SummaryCounts {
	var summaryCounts SummaryCounts

	for _, entries := range ingressReport.Entries {
		switch entries.Status {
		case configs.AnnotationConverted:
			summaryCounts.Converted++
		case configs.AnnotationWarned:
			summaryCounts.Warnings++
		case configs.AnnotationSkipped:
			summaryCounts.Skipped++
		case configs.AnnotationIgnored:
			summaryCounts.Ignored++
		}
	}

	return summaryCounts
}

// summarizeGlobal computes aggregated summary counts across all Ingress reports.
func summarizeGlobal(globalReport configs.GlobalReport) SummaryCounts {
	var total SummaryCounts

	for _, ir := range globalReport.Ingresses {
		summarizedIngress := summarizeIngress(ir)

		total.Converted += summarizedIngress.Converted
		total.Warnings += summarizedIngress.Warnings
		total.Skipped += summarizedIngress.Skipped
		total.Ignored += summarizedIngress.Ignored
	}

	return total
}

// statusLabelColored returns a colored label for a given annotation status.
func statusLabelColored(annotationStatus configs.AnnotationStatus) string {
	switch annotationStatus {
	case configs.AnnotationConverted:
		return color.HiGreenString("Converted")
	case configs.AnnotationWarned:
		return color.HiYellowString("Warning")
	case configs.AnnotationSkipped:
		return color.HiRedString("Skipped")
	case configs.AnnotationIgnored:
		return color.HiBlueString("Ignored")
	default:
		return statusLabel[annotationStatus]
	}
}
