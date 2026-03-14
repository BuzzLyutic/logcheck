package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// Config описывает настройки линтера.
// Загружается из JSON-файла, переданного через флаг -config.
type Config struct {
	Rules RulesConfig `json:"rules"`
}

// RulesConfig управляет включением/выключением отдельных правил.
type RulesConfig struct {
	Lowercase    bool            `json:"lowercase"`
	English      bool            `json:"english"`
	SpecialChars bool            `json:"special_chars"`
	Sensitive    SensitiveConfig `json:"sensitive"`
}

// SensitiveConfig — расширенные настройки правила sensitive data.
type SensitiveConfig struct {
	Enabled  bool     `json:"enabled"`
	Keywords []string `json:"keywords"` // Пустой список = использовать дефолтные.
	Patterns []string `json:"patterns"` // Регулярные выражения для поиска чувствительных данных.
}

// defaultConfig возвращает конфигурацию по умолчанию: все правила включены,
// ключевые слова и паттерны — дефолтные.
func defaultConfig() *Config {
	return &Config{
		Rules: RulesConfig{
			Lowercase:    true,
			English:      true,
			SpecialChars: true,
			Sensitive: SensitiveConfig{
				Enabled: true,
			},
		},
	}
}

// loadConfigFromFile читает конфигурацию из JSON-файла.
func loadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("logcheck: cannot read config %s: %w", path, err)
	}

	cfg := defaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("logcheck: cannot parse config %s: %w", path, err)
	}

	return cfg, nil
}

// buildRules создаёт набор правил на основе конфигурации.
func buildRules(cfg *Config) ([]Rule, error) {
	var rules []Rule

	if cfg.Rules.Lowercase {
		rules = append(rules, &LowercaseRule{})
	}

	if cfg.Rules.English {
		rules = append(rules, &EnglishRule{})
	}

	if cfg.Rules.SpecialChars {
		rules = append(rules, &SpecialCharsRule{})
	}

	if cfg.Rules.Sensitive.Enabled {
		// Компилируем регулярные выражения из конфига.
		var patterns []*regexp.Regexp
		for _, p := range cfg.Rules.Sensitive.Patterns {
			re, err := regexp.Compile(p)
			if err != nil {
				return nil, fmt.Errorf("logcheck: invalid sensitive pattern %q: %w", p, err)
			}
			patterns = append(patterns, re)
		}

		rules = append(rules, &SensitiveRule{
			Keywords: cfg.Rules.Sensitive.Keywords,
			Patterns: patterns,
		})
	}

	return rules, nil
}
