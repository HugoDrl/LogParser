package filter

import (
	"time"

	"github.com/HugoDrl/zebra/internal/parser"
)

type Filters struct {
	StartDate time.Time
	EndDate   time.Time
	Level     parser.Level
	Service   string
}
