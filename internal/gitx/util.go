package gitx

import "strings"

func GetRepoName(path string) string {
	// 提取仓库名
	parts := strings.Split(path, "/")
	repoNameWithExt := parts[len(parts)-1] // 获取最后一个部分，例如 tikuAdapter.git

	// 去除 ".git" 扩展名
	repoName := strings.TrimSuffix(repoNameWithExt, ".git")

	// 转换为小写的驼峰命名格式 tiku-adapter
	return toLowerCamelCase(repoName)
}

func toLowerCamelCase(input string) string {
	// 将驼峰命名格式转换为小写并在单词间加上连字符 "-"
	result := ""
	for i, r := range input {
		if i == 0 || (i > 0 && input[i-1] == ' ') {
			result += string(r)
		} else if r >= 'A' && r <= 'Z' {
			result += "-" + strings.ToLower(string(r))
		} else {
			result += string(r)
		}
	}
	return strings.ReplaceAll(result, " ", "-")
}
