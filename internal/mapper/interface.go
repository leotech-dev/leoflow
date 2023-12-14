package mapper

import (
	"github.com/leotech-dev/leoflow/internal/common"
	"github.com/numaproj/numaflow-go/pkg/mapper"
)

type Mapper interface {
	mapper.Mapper
	common.Initializable
}
