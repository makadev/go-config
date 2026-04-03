package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"

	"go.yaml.in/yaml/v3"
)

type DumpOptions struct {
	// Format specifies the output format: "json", "yaml", "text", "table"
	Format string

	// Content specifies what to include: "config", "env", "metadata", "all"
	Content string

	// MaskSecrets determines if the secrets values should be masked
	MaskSecrets bool

	// MaskWith determines the string used to mask secret values
	MaskWith string
}

type DumpEntry struct {
	ConfigKey  string      `json:"config_key,omitempty"`
	ConfigName string      `json:"config_name,omitempty"`
	EnvVar     string      `json:"env_var,omitempty"`
	FieldPath  string      `json:"field_path,omitempty"`
	Value      interface{} `json:"value"`
	IsSecret   bool        `json:"is_secret,omitempty"`
	IsMasked   bool        `json:"is_masked,omitempty"`
}

func NewDumpOptions() *DumpOptions {
	return &DumpOptions{
		Format:      "table",
		Content:     "config",
		MaskSecrets: true,
		MaskWith:    "***",
	}
}

func (c *Config[T]) Dump() (string, error) {
	return c.DumpWithOptions(&DumpOptions{
		Format:      "table",
		Content:     "config",
		MaskSecrets: true,
		MaskWith:    "***",
	})
}

func (c *Config[T]) DumpWithOptions(options *DumpOptions) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if options == nil {
		return "", fmt.Errorf("dump options cannot be nil")
	}

	entries, err := c.collectDumpEntries(options)
	if err != nil {
		return "", err
	}

	return c.formatDumpEntries(entries, options)
}

func (c *Config[T]) DumpEnv() (string, error) {
	return c.DumpWithOptions(&DumpOptions{
		Format:      "text",
		Content:     "env",
		MaskSecrets: true,
		MaskWith:    "***",
	})
}

func (c *Config[T]) collectDumpEntries(options *DumpOptions) ([]DumpEntry, error) {
	var entries []DumpEntry

	for _, fieldInfo := range c.Metadata.FieldPathMap {
		// Apply content filter
		if !c.shouldIncludeField(fieldInfo, options) {
			continue
		}

		// Get field value
		value, err := c.getFieldValue(fieldInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to get value for field %s: %w", fieldInfo.FieldPath, err)
		}

		// Apply secret handling
		finalValue, isMasked := c.handleSecretValue(value, fieldInfo.Secret, options)

		entry := DumpEntry{
			Value:    finalValue,
			IsSecret: fieldInfo.Secret,
			IsMasked: isMasked,
		}

		// Add metadata based on content type
		switch options.Content {
		case "config":
			entry.ConfigKey = fieldInfo.ConfigKey
		case "env":
			if fieldInfo.EnvVar != "" {
				entry.EnvVar = fieldInfo.EnvVar
			}
		case "metadata", "all":
			entry.ConfigKey = fieldInfo.ConfigKey
			entry.ConfigName = fieldInfo.ConfigName
			entry.EnvVar = fieldInfo.EnvVar
			if options.Content == "all" {
				entry.FieldPath = fieldInfo.FieldPath
			}
		}

		entries = append(entries, entry)
	}

	// Sort entries for consistent output
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].ConfigKey != "" && entries[j].ConfigKey != "" {
			return entries[i].ConfigKey < entries[j].ConfigKey
		}
		if entries[i].EnvVar != "" && entries[j].EnvVar != "" {
			return entries[i].EnvVar < entries[j].EnvVar
		}
		return entries[i].FieldPath < entries[j].FieldPath
	})

	return entries, nil
}

func (c *Config[T]) shouldIncludeField(fieldInfo *FieldInfo, options *DumpOptions) bool {
	switch options.Content {
	case "env":
		return fieldInfo.EnvVar != ""
	case "config", "metadata", "all":
		return true
	default:
		return false
	}
}

func (c *Config[T]) getFieldValue(fieldInfo *FieldInfo) (interface{}, error) {
	field, err := c.getFieldByPath(fieldInfo.FieldPath, false)
	if err != nil {
		return nil, err
	}

	if !field.IsValid() {
		return nil, nil
	}

	return field.Interface(), nil
}

func (c *Config[T]) handleSecretValue(value interface{}, isSecret bool, options *DumpOptions) (interface{}, bool) {
	if !isSecret {
		return value, false
	}

	if options.MaskSecrets {
		return options.MaskWith, true
	}

	return value, false
}

func (c *Config[T]) formatDumpEntries(entries []DumpEntry, options *DumpOptions) (string, error) {
	switch strings.ToLower(options.Format) {
	case "json":
		return c.formatJSON(entries, options)
	case "yaml":
		return c.formatYAML(entries, options)
	case "text":
		return c.formatText(entries, options)
	case "table":
		return c.formatTable(entries, options)
	default:
		return "", fmt.Errorf("unsupported format: %s", options.Format)
	}
}

func (c *Config[T]) formatJSON(entries []DumpEntry, options *DumpOptions) (string, error) {
	switch options.Content {
	case "config":
		data := make(map[string]interface{})
		for _, entry := range entries {
			if entry.ConfigKey != "" {
				setNestedValue(data, entry.ConfigKey, entry.Value)
			}
		}
		bytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	case "env":
		data := make(map[string]interface{})
		for _, entry := range entries {
			if entry.EnvVar != "" {
				data[entry.EnvVar] = entry.Value
			}
		}
		bytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	default:
		bytes, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}

func (c *Config[T]) formatYAML(entries []DumpEntry, options *DumpOptions) (string, error) {
	switch options.Content {
	case "config":
		data := make(map[string]interface{})
		for _, entry := range entries {
			if entry.ConfigKey != "" {
				setNestedValue(data, entry.ConfigKey, entry.Value)
			}
		}
		bytes, err := yaml.Marshal(data)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	case "env":
		data := make(map[string]interface{})
		for _, entry := range entries {
			if entry.EnvVar != "" {
				data[entry.EnvVar] = entry.Value
			}
		}
		bytes, err := yaml.Marshal(data)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	default:
		bytes, err := yaml.Marshal(entries)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}

func setNestedValue(data map[string]interface{}, key string, value interface{}) {
	parts := strings.Split(key, ".")
	for i, part := range parts {
		if i == len(parts)-1 {
			data[part] = value
		} else {
			existing, ok := data[part]
			if !ok {
				data[part] = make(map[string]interface{})
			} else if _, isMap := existing.(map[string]interface{}); !isMap {
				// Non-map value already at this key (e.g. a struct-level entry was
				// processed before its leaf children); replace it with a nested map.
				data[part] = make(map[string]interface{})
			}
			data = data[part].(map[string]interface{})
		}
	}
}


func nonprimitiveToString(value interface{}) (interface{}, bool) {
	kind := reflect.TypeOf(value).Kind()
	if kind == reflect.Invalid {
		return nil, false
	}
	if kind == reflect.Ptr {
		kind = reflect.TypeOf(value).Elem().Kind()
	}
	switch kind {
	case reflect.Struct:
		// Skip struct fields
		return nil, false
	case reflect.Slice, reflect.Array:
		// Handle slice and array fields -> convert to string for printing
		len := reflect.ValueOf(value).Len()
		valLines := make([]string, len)
		for i := 0; i < len; i++ {
			item := reflect.ValueOf(value).Index(i)
			valLines[i] = fmt.Sprintf("%v", item.Interface())
		}
		value = strings.Join(valLines, ",")
	case reflect.Map:
		len := reflect.ValueOf(value).Len()
		valLines := make([]string, len)
		for i, key := range reflect.ValueOf(value).MapKeys() {
			valLines[i] = fmt.Sprintf("%s=%v", key, reflect.ValueOf(value).MapIndex(key).Interface())
		}
		value = strings.Join(valLines, ",")
	}
	return value, true
}

func (c *Config[T]) formatText(entries []DumpEntry, options *DumpOptions) (string, error) {
	var lines []string

	for _, entry := range entries {
		val, ok := nonprimitiveToString(entry.Value)
		if !ok {
			continue
		}

		switch options.Content {
		case "config":
			if entry.ConfigKey != "" {
				lines = append(lines, fmt.Sprintf("%s=%v", entry.ConfigKey, val))
			}
		case "env":
			if entry.EnvVar != "" {
				lines = append(lines, fmt.Sprintf("%s=%v", entry.EnvVar, val))
			}
		default:
			line := ""
			if entry.ConfigKey != "" {
				line += fmt.Sprintf("ConfigKey=%s ", entry.ConfigKey)
			}
			if entry.EnvVar != "" {
				line += fmt.Sprintf("EnvVar=%s ", entry.EnvVar)
			}
			if entry.FieldPath != "" {
				line += fmt.Sprintf("FieldPath=%s ", entry.FieldPath)
			}
			line += fmt.Sprintf("Value=%v", val)
			if entry.IsSecret {
				line += " (secret)"
			}
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n"), nil
}

func (c *Config[T]) formatTable(entries []DumpEntry, options *DumpOptions) (string, error) {
	if len(entries) == 0 {
		return "", nil
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)

	switch options.Content {
	case "config":
		fmt.Fprintln(w, "CONFIG_KEY\tVALUE\tSECRET")
		fmt.Fprintln(w, "----------\t-----\t------")
		for _, entry := range entries {
			val, ok := nonprimitiveToString(entry.Value)
			if !ok {
				continue
			}
			if entry.ConfigKey != "" {
				secret := ""
				if entry.IsSecret {
					secret = "yes"
				}
				fmt.Fprintf(w, "%s\t%v\t%s\n", entry.ConfigKey, val, secret)
			}
		}
	case "env":
		fmt.Fprintln(w, "ENV_VAR\tVALUE\tSECRET")
		fmt.Fprintln(w, "-------\t-----\t------")
		for _, entry := range entries {
			val, ok := nonprimitiveToString(entry.Value)
			if !ok {
				continue
			}
			if entry.EnvVar != "" {
				secret := ""
				if entry.IsSecret {
					secret = "yes"
				}
				fmt.Fprintf(w, "%s\t%v\t%s\n", entry.EnvVar, val, secret)
			}
		}
	default:
		if options.Content == "all" {
			fmt.Fprintln(w, "CONFIG_KEY\tENV_VAR\tFIELD_PATH\tVALUE\tSECRET")
			fmt.Fprintln(w, "----------\t-------\t----------\t-----\t------")
		} else {
			fmt.Fprintln(w, "CONFIG_KEY\tENV_VAR\tVALUE\tSECRET")
			fmt.Fprintln(w, "----------\t-------\t-----\t------")
		}

		for _, entry := range entries {
			val, ok := nonprimitiveToString(entry.Value)
			if !ok {
				continue
			}
			secret := ""
			if entry.IsSecret {
				secret = "yes"
			}

			if options.Content == "all" {
				fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n",
					entry.ConfigKey, entry.EnvVar, entry.FieldPath, val, secret)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%v\t%s\n",
					entry.ConfigKey, entry.EnvVar, val, secret)
			}
		}
	}

	w.Flush()
	return strings.TrimSuffix(buf.String(), "\n"), nil
}
