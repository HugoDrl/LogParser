package analyser

import (
	"testing"

	"github.com/HugoDrl/zebra/internal/parser"
	"github.com/google/go-cmp/cmp"
)

func TestNewMetrics(t *testing.T) {
	metric := newMetrics()
	expected := &CollectionMetric{
		Lines:              map[parser.Level]int{},
		ServicePerformance: map[string]ServiceMetric{},
	}
	if diff := cmp.Diff(expected, metric); diff != "" {
		t.Fatal(diff)
	}
}
