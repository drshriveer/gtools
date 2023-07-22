package config

type Config struct {
}

func GetEnv[T ~int | ~uint](cfg *Config, key string) T {

}

func MustGet[T any](cfg *Config, key string) T {

}
