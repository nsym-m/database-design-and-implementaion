package log_test

import (
	"fmt"
	"testing"

	"github.com/nsym-m/simpledb/internal/file"
	"github.com/nsym-m/simpledb/internal/log"
)

func TestAppender(t *testing.T) {
	tmpDir := t.TempDir()
	bs, err := file.NewBlockStore(tmpDir, 400)
	if err != nil {
		t.Fatal(err)
	}
	lm, err := log.NewAppender(bs, "test")
	if err != nil {
		t.Fatal(err)
	}
	createRecords(t, lm, 1, 35)
	t.Logf("lm1: %+v\n", lm)
	printeds := printLogRecords(t, lm, "The log file now has these records:")
	if len(printeds) != 35 {
		t.Errorf("printedsが35でない: %d\n", len(printeds))
	}
	t.Logf("printeds: %+v\n", printeds)
	t.Logf("lm2: %+v\n", lm)
	createRecords(t, lm, 36, 70)
	t.Logf("lm3: %+v\n", lm)
	if err := lm.Flush(65); err != nil {
		t.Errorf("lm.Flush(65): %v\n", err)
	}
	t.Logf("lm4: %+v\n", lm)
	printeds2 := printLogRecords(t, lm, "The log file now has these records:")
	if len(printeds2) != 70 {
		t.Errorf("printeds2が70でない: %d\n", len(printeds2))
	}
	t.Logf("printeds2: %+v\n", printeds2)
	t.Logf("lm5: %+v\n", lm)
}

func printLogRecords(t *testing.T, lm log.Appender, msg string) []int {
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

func createRecords(t *testing.T, lm log.Appender, start, end int) {
	t.Log("creating records")
	for i := start; i <= end; i++ {
		rec := createLogRecord(t, fmt.Sprintf("record%d", i), i+100)
		if _, err := lm.Append(rec); err != nil {
			t.Fatal(err)
		}
	}
}

func createLogRecord(t *testing.T, s string, n int) []byte {
	npos := file.MaxLength(len(s))
	b := make([]byte, npos+file.Int32Bytes)
	p := file.NewPageFromBytes(b)
	if err := p.SetString(0, s); err != nil {
		t.Errorf("SetString(0, s) s: %v, error: %v", s, err)
	}
	if err := p.SetInt(npos, n); err != nil {
		t.Errorf("SetInt(npos, n) npos: %v, n: %v, error: %v", npos, n, err)
	}
	return b
}
