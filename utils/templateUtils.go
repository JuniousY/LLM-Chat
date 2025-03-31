package utils

import (
	"regexp"
	"strings"
)

func ReplacePlaceholders(template string, params map[string]string) string {
	// 编译正则表达式，匹配 {{#param_name#}} 格式的占位符
	re := regexp.MustCompile(`{{#([^#]+)#}}`)

	// 查找并替换所有匹配的占位符
	result := re.ReplaceAllStringFunc(template, func(match string) string {
		// 提取占位符中的参数名
		paramName := strings.TrimPrefix(strings.TrimSuffix(match, "#}}"), "{{#")
		// 检查参数名是否存在于 params 中
		if value, exists := params[paramName]; exists {
			return value
		}
		// 如果参数名不存在，返回原占位符
		return match
	})

	return result
}
