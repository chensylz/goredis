package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Bind string
	Port int
}

func Setup(filePath string) *Config {
	config := &Config{
		Bind: "127.0.0.1",
		Port: 6379,
	}
	properties, err := parseConfig(filePath)
	if err != nil {
		log.Panicf("load local config error: %v", err)
	}
	config.Bind = properties["bind"]
	config.Port, _ = strconv.Atoi(properties["port"])
	return config
}

func parseConfig(filePath string) (map[string]string, error) {
	config := make(map[string]string)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 忽略注释和空行
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			config[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}
