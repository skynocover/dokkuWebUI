package utils

import (
	"strings"
)

type Info struct {
	Key   string
	Value string
}

func Parse(s string) (infos []Info) {
	ss := strings.Split(s, "\n")
	for i := range ss {
		if i == 0 {
			continue
		}
		sp := strings.Split(ss[i], ":")
		if len(sp) > 1 {
			infos = append(infos, Info{Key: strings.TrimSpace(sp[0]), Value: strings.TrimSpace(strings.Replace(ss[i], sp[0]+":", "", 1))})
		} else if len(sp) > 0 && sp[0] != "" {
			infos = append(infos, Info{Key: strings.TrimSpace(sp[0]), Value: ""})
		}

	}
	return
}
