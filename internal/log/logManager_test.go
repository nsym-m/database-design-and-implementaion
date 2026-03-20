package log_test

import (
	"fmt"
	"testing"

	"github.com/nsym-m/simpledb/internal/file"
	"github.com/nsym-m/simpledb/internal/log"
)

func TestLogManager(t *testing.T) {
	tmpDir := t.TempDir()
	fm, err := file.NewFileManager(tmpDir, 400)
	if err != nil {
		t.Fatal(err)
	}
	lm, err := log.NewLogManager(fm, "test")
	if err != nil {
		t.Fatal(err)
	}
	createRecords(t, lm, 1, 35)
	printeds := printLogRecords(t, lm, "The log file now has these records:")
	if len(printeds) != 35 {
		t.Errorf("printedsが35でない: %d\n", len(printeds))
	}
	createRecords(t, lm, 36, 70)
	if err := lm.Flush(65); err != nil {
		t.Errorf("lm.Flush(65): %v\n", err)
	}
	printeds2 := printLogRecords(t, lm, "The log file now has these records:")
	if len(printeds2) != 5 {
		t.Errorf("printeds2が5でない: %d\n", len(printeds2))
	}
}

func printLogRecords(t *testing.T, lm *log.LogManager, msg string) []int {
	t.Log(msg)
	res := []int{}
	for rec, err := range lm.All() {
		if err != nil {
			t.Fatal(err)
		}
		p := file.NewPageFromBytes(rec)
		s := p.GetString(0)
		res = append(res, p.GetInt(file.MaxLength(len(s))))
	}
	return res
}

func createRecords(t *testing.T, lm *log.LogManager, start, end int) {
	t.Log("creating records")
	for i := start; i <= end; i++ {
		rec := createLogRecord(t, fmt.Sprintf("record%d", i), i+100)
		lsn, err := lm.Append(rec)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("lsn: %d\n", lsn)
	}
}

func createLogRecord(t *testing.T, s string, n int) []byte {
	npos := file.MaxLength(len(s))
	b := make([]byte, npos+log.Bytes)
	p := file.NewPageFromBytes(b)
	if err := p.SetString(0, s); err != nil {
		t.Errorf("SetString(0, s) s: %v, error: %v", s, err)
	}
	if err := p.SetInt(npos, n); err != nil {
		t.Errorf("SetInt(npos, n) npos: %v, n: %v, error: %v", npos, n, err)
	}
	return b
}
