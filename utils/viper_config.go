package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type ViperConfiguration struct {
}

func (vc *ViperConfiguration) Get(key string) interface{} {
	return viper.Get(key)
}
func (vc *ViperConfiguration) GetBool(key string) bool {
	return viper.GetBool(key)
}
func (vc *ViperConfiguration) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}
func (vc *ViperConfiguration) GetInt(key string) int {
	return viper.GetInt(key)
}
func (vc *ViperConfiguration) GetString(key string) (string, error) {
	if !viper.IsSet(key) {
		return "", fmt.Errorf("Could not find key %s", key)
	}
	return viper.GetString(key), nil
}
func (vc *ViperConfiguration) GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}
func (vc *ViperConfiguration) GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}
func (vc *ViperConfiguration) GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}
func (vc *ViperConfiguration) GetTime(key string) time.Time {
	return viper.GetTime(key)
}
func (vc *ViperConfiguration) GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}
func (vc *ViperConfiguration) IsSet(key string) bool {
	return viper.IsSet(key)
}
func (vc *ViperConfiguration) Unmarshal(rawVal interface{}) error {
	return viper.Unmarshal(rawVal)
}
func (vc *ViperConfiguration) UnmarshalKey(key string, dest interface{}) error {
	return viper.UnmarshalKey(key, dest)
}
