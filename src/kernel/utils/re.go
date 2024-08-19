package utils

import "regexp"

func Match(reStr string, target string) bool {
	if reStr[0] != 94 { // '^'
		reStr = string(append([]byte{94}, []byte(reStr)...))
	}

	if reStr[len(reStr) - 1] != 36 { // '$'
		reStr = string(append([]byte(reStr), []byte{36}...))
	}

	re, _ := regexp.Compile(reStr)

	return re.MatchString(target)
}

func FindOneSubmatch(reStr string, target string) []string {
	re, _ := regexp.Compile(reStr)

	return re.FindStringSubmatch(target)
}

func FindAllSubmatch(reStr string, target string) [][]string {
	re, _ := regexp.Compile(reStr)

	// res := make([][]string, 0)
	
	// for _, match := range re.FindAllStringSubmatch(target, -1) {
	// 	res = append(res, match[1:])
	// }

	// return res

	return re.FindAllStringSubmatch(target, -1)
}
