package utils

// Truncate 把字符串 s 从头最多保留 n 个字符，超出部分替换为 "..."。
// 按 Unicode 字符（rune）计数，中文、emoji 都不会被截成半个导致乱码。
// n <= 0 且 s 非空时，返回 "..."。
func Truncate(s string, n int) string {
	if n <= 0 {
		if s == "" {
			return ""
		}
		return "..."
	}

	// len(s) 是字节数，字节数 <= n 时字符数必然 <= n，直接返回
	if len(s) <= n {
		return s
	}

	// range string 遍历的是 rune，i 是每个 rune 的起始字节下标
	count := 0
	for i := range s {
		if count == n {
			return s[:i] + "..."
		}
		count++
	}
	return s
}
