package csvsource

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mohammadalviyan/rpmp/backend/internal/domain"
)

var requiredHeaders = []string{
	"ProcessId",
	"ProcessName",
	"PackageName",
	"EnvironmentName",
	"FullyQualifiedName",
	"CountExecuting",
	"CountPending",
	"CountSuspended",
	"CountResumed",
	"CountSuccessful",
	"CountErrors",
	"CountStopped",
	"AverageDurationInSeconds",
	"AveragePendingTimeInSeconds",
	"TotalRows",
	"EntityId",
}

var utf8BOM = []byte{0xef, 0xbb, 0xbf}

type Source struct {
	path string
	now  func() time.Time
}

func New(path string) *Source {
	return &Source{path: path, now: time.Now}
}

func (s *Source) Read(ctx context.Context) (domain.AggregateSnapshot, error) {
	if err := ctx.Err(); err != nil {
		return domain.AggregateSnapshot{}, err
	}
	content, err := os.ReadFile(s.path)
	if err != nil {
		return domain.AggregateSnapshot{}, fmt.Errorf("read source CSV: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return domain.AggregateSnapshot{}, err
	}
	snapshot, err := parse(content, s.now().UTC())
	if err != nil {
		return domain.AggregateSnapshot{}, err
	}
	return snapshot, nil
}

type sourceRow struct {
	processID                   string
	processName                 string
	packageName                 string
	environmentName             *string
	fullyQualifiedName          string
	countExecuting              int32
	countPending                int32
	countSuspended              int32
	countResumed                int32
	countSuccessful             int32
	countErrors                 int32
	countStopped                int32
	averageDurationInSeconds    *float64
	averagePendingTimeInSeconds *float64
	totalRows                   int32
	entityID                    string
}

func parse(content []byte, importedAt time.Time) (domain.AggregateSnapshot, error) {
	reader := csv.NewReader(bytes.NewReader(bytes.TrimPrefix(content, utf8BOM)))
	reader.Comma = ','
	reader.FieldsPerRecord = len(requiredHeaders)
	reader.ReuseRecord = false

	header, err := reader.Read()
	if err != nil {
		return domain.AggregateSnapshot{}, fmt.Errorf("read CSV header: %w", err)
	}
	if err := validateHeader(header); err != nil {
		return domain.AggregateSnapshot{}, err
	}

	var rows []domain.ProcessAggregate
	processIDs := make(map[string]struct{})
	entityIDs := make(map[string]struct{})
	for line := 2; ; line++ {
		record, readErr := reader.Read()
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			return domain.AggregateSnapshot{}, fmt.Errorf("parse CSV row %d: %w", line, readErr)
		}
		source, parseErr := parseRow(record, line)
		if parseErr != nil {
			return domain.AggregateSnapshot{}, parseErr
		}
		if _, exists := processIDs[source.processID]; exists {
			return domain.AggregateSnapshot{}, fmt.Errorf("CSV row %d: duplicate ProcessId", line)
		}
		if _, exists := entityIDs[source.entityID]; exists {
			return domain.AggregateSnapshot{}, fmt.Errorf("CSV row %d: duplicate EntityId", line)
		}
		processIDs[source.processID] = struct{}{}
		entityIDs[source.entityID] = struct{}{}
		rows = append(rows, mapRow(source))
	}
	if len(rows) == 0 {
		return domain.AggregateSnapshot{}, errors.New("source CSV contains no data rows")
	}

	sum := sha256.Sum256(content)
	return domain.AggregateSnapshot{
		SourceSnapshotKey: "csv-sha256:" + hex.EncodeToString(sum[:]),
		ImportedAt:        importedAt,
		Rows:              rows,
	}, nil
}

func validateHeader(header []string) error {
	for i, required := range requiredHeaders {
		if header[i] != required {
			return fmt.Errorf("invalid CSV header at column %d: expected %s", i+1, required)
		}
	}
	return nil
}

func parseRow(record []string, line int) (sourceRow, error) {
	processID, err := parseID(record[0], "ProcessId", line)
	if err != nil {
		return sourceRow{}, err
	}
	processName, err := requiredText(record[1], "ProcessName", line)
	if err != nil {
		return sourceRow{}, err
	}
	packageName, err := requiredText(record[2], "PackageName", line)
	if err != nil {
		return sourceRow{}, err
	}
	environmentName, err := nullableText(record[3], "EnvironmentName", line)
	if err != nil {
		return sourceRow{}, err
	}
	fullyQualifiedName, err := requiredText(record[4], "FullyQualifiedName", line)
	if err != nil {
		return sourceRow{}, err
	}
	values := make([]int32, 0, 8)
	for index, name := range []string{
		"CountExecuting", "CountPending", "CountSuspended", "CountResumed",
		"CountSuccessful", "CountErrors", "CountStopped", "TotalRows",
	} {
		column := 5 + index
		if name == "TotalRows" {
			column = 14
		}
		value, parseErr := parseNonnegativeInt(record[column], name, line)
		if parseErr != nil {
			return sourceRow{}, parseErr
		}
		values = append(values, value)
	}
	averageDuration, err := parseNullableNonnegativeFloat(record[12], "AverageDurationInSeconds", line)
	if err != nil {
		return sourceRow{}, err
	}
	averagePending, err := parseNullableNonnegativeFloat(record[13], "AveragePendingTimeInSeconds", line)
	if err != nil {
		return sourceRow{}, err
	}
	entityID, err := parseID(record[15], "EntityId", line)
	if err != nil {
		return sourceRow{}, err
	}
	return sourceRow{
		processID: processID, processName: processName, packageName: packageName,
		environmentName: environmentName, fullyQualifiedName: fullyQualifiedName,
		countExecuting: values[0], countPending: values[1], countSuspended: values[2],
		countResumed: values[3], countSuccessful: values[4], countErrors: values[5],
		countStopped: values[6], averageDurationInSeconds: averageDuration,
		averagePendingTimeInSeconds: averagePending, totalRows: values[7], entityID: entityID,
	}, nil
}

func mapRow(row sourceRow) domain.ProcessAggregate {
	return domain.ProcessAggregate{
		SourceProcessKey: row.processID, UseCaseSourceKey: row.fullyQualifiedName,
		UseCaseName: row.fullyQualifiedName, ProcessName: row.processName,
		PackageName: row.packageName, EnvironmentName: row.environmentName,
		ExecutingCount: row.countExecuting, PendingCount: row.countPending,
		SuspendedCount: row.countSuspended, ResumedCount: row.countResumed,
		SuccessfulCount: row.countSuccessful, ErrorCount: row.countErrors,
		StoppedCount: row.countStopped, AverageDurationSeconds: row.averageDurationInSeconds,
		AveragePendingSeconds: row.averagePendingTimeInSeconds,
		SourceTotalRows:       row.totalRows, SourceEntityKey: row.entityID,
	}
}

func parseID(value, name string, line int) (string, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return "", fmt.Errorf("CSV row %d: %s must be a positive integer", line, name)
	}
	return strconv.FormatInt(parsed, 10), nil
}

func parseNonnegativeInt(value, name string, line int) (int32, error) {
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("CSV row %d: %s must be a nonnegative integer", line, name)
	}
	return int32(parsed), nil
}

func parseNullableNonnegativeFloat(value, name string, line int) (*float64, error) {
	if value == "NULL" {
		return nil, nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || parsed < 0 || math.IsNaN(parsed) || math.IsInf(parsed, 0) {
		return nil, fmt.Errorf("CSV row %d: %s must be NULL or a nonnegative number", line, name)
	}
	return &parsed, nil
}

func requiredText(value, name string, line int) (string, error) {
	if strings.TrimSpace(value) == "" || value == "NULL" {
		return "", fmt.Errorf("CSV row %d: %s is required", line, name)
	}
	return value, nil
}

func nullableText(value, name string, line int) (*string, error) {
	if value == "NULL" {
		return nil, nil
	}
	if strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("CSV row %d: %s must be NULL or nonempty", line, name)
	}
	copy := value
	return &copy, nil
}
