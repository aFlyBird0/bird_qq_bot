package bot

import (
	"bird_qq_bot/config"
	"time"
)

const prefix = "modules." // 子模块配置文件的顶层前缀

// GetModConfigPath 根据模块与键获得完整子模块配置文件路径
func GetModConfigPath(m Module, key string) string {
	return prefix + m.GetModuleInfo().ID.String() + "." + key
}

// GetModConfigString 根据模块与键获得String类型配置信息
func GetModConfigString(m Module, key string) string {
	return config.GlobalConfig.GetString(GetModConfigPath(m, key))
}

// GetModConfigInt 根据模块与键获得Int类型配置信息
func GetModConfigInt(m Module, key string) int {
	return config.GlobalConfig.GetInt(GetModConfigPath(m, key))
}

func GetModConfigInt64(m Module, key string) int64 {
	return config.GlobalConfig.GetInt64(GetModConfigPath(m, key))
}

func GetModConfigBool(m Module, key string) bool {
	return config.GlobalConfig.GetBool(GetModConfigPath(m, key))
}

func GetModConfigFloat64(m Module, key string) float64 {
	return config.GlobalConfig.GetFloat64(GetModConfigPath(m, key))
}

func GetModConfigTime(m Module, key string) time.Time {
	return config.GlobalConfig.GetTime(GetModConfigPath(m, key))
}

func GetModConfigDuration(m Module, key string) time.Duration {
	return config.GlobalConfig.GetDuration(GetModConfigPath(m, key))
}

func GetModConfigStringSlice(m Module, key string) []string {
	return config.GlobalConfig.GetStringSlice(GetModConfigPath(m, key))
}

func GetModConfigIntSlice(m Module, key string) []int {
	return config.GlobalConfig.GetIntSlice(GetModConfigPath(m, key))
}

func GetModConfigInt64Slice(m Module, key string) []int64 {
	nums := config.GlobalConfig.GetIntSlice(GetModConfigPath(m, key))
	numsInt64 := make([]int64, len(nums))
	for i, num := range nums {
		numsInt64[i] = int64(num)
	}
	return numsInt64
}
