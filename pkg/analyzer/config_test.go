package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	if !cfg.Rules.Lowercase {
		t.Error("Lowercase должен быть включён по умолчанию")
	}
	if !cfg.Rules.English {
		t.Error("English должен быть включён по умолчанию")
	}
	if !cfg.Rules.SpecialChars {
		t.Error("SpecialChars должен быть включён по умолчанию")
	}
	if !cfg.Rules.Sensitive.Enabled {
		t.Error("Sensitive должен быть включён по умолчанию")
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	content := `{
		"rules": {
			"lowercase": true,
			"english": false,
			"special_chars": true,
			"sensitive": {
				"enabled": true,
				"keywords": ["ssn", "credit_card"],
				"patterns": ["(?i)bearer\\s+"]
			}
		}
	}`

	dir := t.TempDir()
	path := filepath.Join(dir, ".logcheck.json")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := loadConfigFromFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !cfg.Rules.Lowercase {
		t.Error("Lowercase должен быть true")
	}
	if cfg.Rules.English {
		t.Error("English должен быть false")
	}
	if !cfg.Rules.SpecialChars {
		t.Error("SpecialChars должен быть true")
	}
	if !cfg.Rules.Sensitive.Enabled {
		t.Error("Sensitive должен быть включён")
	}
	if len(cfg.Rules.Sensitive.Keywords) != 2 {
		t.Errorf("ожидалось 2 ключевых слова, получено %d", len(cfg.Rules.Sensitive.Keywords))
	}
	if len(cfg.Rules.Sensitive.Patterns) != 1 {
		t.Errorf("ожидался 1 паттерн, получено %d", len(cfg.Rules.Sensitive.Patterns))
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{not json}`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := loadConfigFromFile(path)
	if err == nil {
		t.Error("ожидалась ошибка при невалидном JSON")
	}
}

func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := loadConfigFromFile("/nonexistent/path.json")
	if err == nil {
		t.Error("ожидалась ошибка при отсутствующем файле")
	}
}

func TestBuildRulesAllEnabled(t *testing.T) {
	cfg := defaultConfig()
	rules, err := buildRules(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 4 {
		t.Errorf("ожидалось 4 правила, получено %d", len(rules))
	}
}

func TestBuildRulesPartialDisable(t *testing.T) {
	cfg := defaultConfig()
	cfg.Rules.English = false
	cfg.Rules.SpecialChars = false

	rules, err := buildRules(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 2 {
		t.Errorf("ожидалось 2 правила, получено %d", len(rules))
	}
}

func TestBuildRulesInvalidPattern(t *testing.T) {
	cfg := defaultConfig()
	cfg.Rules.Sensitive.Patterns = []string{"[invalid"}

	_, err := buildRules(cfg)
	if err == nil {
		t.Error("ожидалась ошибка при невалидном регулярном выражении")
	}
}

func TestBuildRulesWithCustomPatterns(t *testing.T) {
	cfg := defaultConfig()
	cfg.Rules.Sensitive.Patterns = []string{"(?i)bearer\\s+", "(?i)auth[_-]?token"}

	rules, err := buildRules(cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Находим SensitiveRule среди правил.
	var sensitive *SensitiveRule
	for _, r := range rules {
		if sr, ok := r.(*SensitiveRule); ok {
			sensitive = sr
			break
		}
	}
	if sensitive == nil {
		t.Fatal("SensitiveRule не найден в списке правил")
	}
	if len(sensitive.Patterns) != 2 {
		t.Errorf("ожидалось 2 паттерна, получено %d", len(sensitive.Patterns))
	}
}
