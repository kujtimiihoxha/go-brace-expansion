package gobrex

import (
	"fmt"
	"github.com/kujtimiihoxha/go-balanced-match"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var (
	escSlash  = "\u0000SLASH" + fmt.Sprintf("%f", rand.Float32()) + "\u0000"
	escOpen   = "\u0000OPEN" + fmt.Sprintf("%f", rand.Float32()) + "\u0000"
	escClose  = "\u0000CLOSE" + fmt.Sprintf("%f", rand.Float32()) + "\u0000"
	escComma  = "\u0000COMMA" + fmt.Sprintf("%f", rand.Float32()) + "\u0000"
	escPeriod = "\u0000PERIOD" + fmt.Sprintf("%f", rand.Float32()) + "\u0000"
)

func parseCommaParts(str string) []string {
	if str == "" {
		return []string{""}
	}
	parts := []string{}
	m := goblma.Balanced("{", "}", str)
	if m == nil {
		return strings.Split(str, ",")
	}
	pre := m.Pre
	body := m.Body
	post := m.Post
	p := strings.Split(pre, ",")
	p[len(p)-1] += "{" + body + "}"
	postParts := parseCommaParts(post)
	if len(post) > 0 {
		p[len(p)-1] += postParts[0]
		postParts = postParts[1:]
		p = append(p, postParts...)
	}
	parts = append(parts, p...)
	return parts
}
func expand(str string, isTop bool) []string {
	expansions := []string{}
	m := goblma.Balanced("{", "}", str)
	reg := regexp.MustCompile(`\$$`)
	if m == nil || reg.Match([]byte(m.Pre)) {
		return []string{str}
	}
	isNumericSequence := regexp.MustCompile(`^-?\d+\.\.-?\d+(?:\.\.-?\d+)?$`).Match([]byte(m.Body))
	isAlphaSequence := regexp.MustCompile(`^[a-zA-Z]\.\.[a-zA-Z](?:\.\.-?\d+)?$`).Match([]byte(m.Body))
	isSequence := isNumericSequence || isAlphaSequence
	isOptions := regexp.MustCompile(`^(.*,)+(.+)?$`).Match([]byte(m.Body))
	if !isSequence && !isOptions {
		if regexp.MustCompile(`,.*\}`).Match([]byte(m.Post)) {
			str = m.Pre + "{" + m.Body + escClose + m.Post
			return expand(str, false)
		}
		return []string{str}
	}
	var n []string
	if isSequence {
		n = strings.Split(m.Body, `..`)
	} else {
		n = parseCommaParts(m.Body)
		if len(n) == 1 {
			//// x{{a,b}}y ==> x{a}y x{b}y
			n = mapArray(expand(n[0], false), embrace)
			if len(n) == 1 {
				post := []string{}
				if len(m.Post) > 0 {
					post = expand(m.Post, false)
				} else {
					post = []string{""}
				}
				return mapArray(post, func(s string) string {
					return m.Pre + n[0] + s
				})
			}
		}
	}
	pre := m.Pre
	post := []string{}
	if len(m.Post) > 0 {
		post = expand(m.Post, false)
	} else {
		post = []string{""}
	}
	N := []string{}
	if isSequence {
		x := numeric(n[0])
		y := numeric(n[1])
		width := max(len(n[0]), len(n[1]))
		var incr int64
		if len(n) == 3 {
			incr = int64(math.Abs(float64(numeric(n[2]))))
		} else {
			incr = 1
		}
		test := lte
		reverse := y < x
		if reverse {
			incr *= -1
			test = gte
		}
		pad := some(n, isPadded)
		N = []string{}
		for i := x; test(i, y); i += incr {
			var c string
			if isAlphaSequence {
				c = string(i)
				if c == "\\" {
					c = ""
				}
			} else {
				c = strconv.FormatInt(i, 10)
				if pad {
					var need = width - len(c)
					if need > 0 {
						var z = strings.Join(make([]string, need+1), "0")
						if i < 0 {
							c = "-" + z + c[1:]
						} else {
							c = z + c
						}

					}
				}
			}
			N = append(N, c)
		}
	} else {
		N = concatMap(n, func(el string) []string { return expand(el, false) })
	}

	for j := 0; j < len(N); j++ {
		for k := 0; k < len(post); k++ {
			expansion := pre + N[j] + post[k]
			if !isTop || isSequence || expansion != "" {
				expansions = append(expansions, expansion)
			}
		}
	}
	return expansions
}
func concatMap(xs []string, fn func(el string) []string) []string {
	res := []string{}
	for i := 0; i < len(xs); i++ {
		var x = fn(xs[i])
		res = append(res, x...)
	}
	return res
}
func some(arr []string, fn func(el string) bool) bool {
	for _, v := range arr {
		if fn(v) {
			return true
		}
	}
	return false
}
func Expand(str string) []string {
	if str == "" {
		return []string{}
	}

	return mapArray(expand(escapeBraces(str), true), unescapeBraces)
}
func mapArray(arr []string, call func(str string) string) []string {
	ret := []string{}
	for _, v := range arr {
		ret = append(ret, call(v))
	}
	return ret
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
func isPadded(el string) bool {
	return regexp.MustCompile(`^-?0\d`).Match([]byte(el))
}
func embrace(str string) string {
	return "{" + str + "}"
}
func lte(i int64, y int64) bool {
	return i <= y
}

func gte(i int64, y int64) bool {
	return i >= y
}

func numeric(str string) int64 {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return int64(str[0])
	}
	return v
}

func escapeBraces(str string) string {
	return sliceAndJoin(sliceAndJoin(sliceAndJoin(sliceAndJoin(sliceAndJoin(str, escSlash, "\\\\"), escOpen, "\\{"), escClose, "\\}"), escComma, "\\,"), escPeriod, "\\.")

}
func unescapeBraces(str string) string {
	return sliceAndJoin(sliceAndJoin(sliceAndJoin(sliceAndJoin(sliceAndJoin(str, "\\", escSlash), "{", escOpen), "}", escClose), ",", escComma), ".", escPeriod)

}
func sliceAndJoin(str string, join string, slice string) string {
	return strings.Join(strings.Split(str, slice), join)
}
