package request

import (
	"github.com/daidr/doulog-core/lib/conf"
	"github.com/guonaihong/gout"
	"github.com/guonaihong/gout/dataflow"
)

func HTTP() *dataflow.DataFlow {
	return gout.New().Debug(conf.C.Debug)
}
