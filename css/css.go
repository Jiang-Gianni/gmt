package css

import (
	"regexp"
	"strings"
)

func init() {
	TailwindMap = extractKeyValuePairs(TailwindCSS)
}

func GetStyles(classes []string) []string {
	re := regexp.MustCompile(`animate-(.*?)`)
	var styles []string
	for _, class := range classes {
		style, ok := TailwindMap[class]
		isAnimate := re.MatchString(class)
		if ok {
			styles = append(styles, style)
			if isAnimate {
				styles = append(styles, TailwindKeyFramesMap[class])
			}
		}
	}
	return styles
}

func GetClasses(input string) []string {
	var resultClasses []string
	re := regexp.MustCompile(`class="(.*?)"`)
	alreadySeen := map[string]struct{}{}
	yes := struct{}{}
	matches := re.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		if len(match) == 2 {
			classText := match[1]
			classes := strings.Split(classText, " ")
			for _, class := range classes {
				if _, ok := alreadySeen[class]; !ok {
					alreadySeen[class] = yes
					resultClasses = append(resultClasses, "."+class)
				}
			}
		}
	}
	return resultClasses
}

func extractKeyValuePairs(input string) map[string]string {
	re := regexp.MustCompile(`(.*?)\{(.*?)\}`)
	matches := re.FindAllStringSubmatch(input, -1)
	keyValuePairs := make(map[string]string)
	for _, match := range matches {
		if len(match) == 3 {
			key := match[1]
			value := match[1] + "{" + match[2] + "}"
			keyValuePairs[key] = value
		}
	}
	return keyValuePairs
}
