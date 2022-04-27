package globalconfig

import "github.com/chack93/scrumpoker_api/internal/domain/common"

// GlobalConfig defines model for GlobalConfig.
type GlobalConfig struct {
	common.BaseModel `yaml:",inline"`
	Key              *string `json:"key,omitempty"`
	Value            *string `json:"Value,omitempty"`
}

// GlobalConfigList defines model for GlobalConfigList.
type GlobalConfigList []GlobalConfig
