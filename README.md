# Go bash style brace expansion

This is a go implementation of [juliangruber/brace-expansion](https://github.com/juliangruber/brace-expansion)

#### Usage
```go
package main

import (
	"github.com/kujtimiihoxha/go-brace-expansion"
	"fmt"
)

func main() {
	v:=gobrex.Expand("file-{a,b,c}.jpg")
	fmt.Println(v)	// => ['file-a.jpg', 'file-b.jpg', 'file-c.jpg']
	
	v =gobrex.Expand("-v{,,}")
	fmt.Println(v)	// => ['-v', '-v', '-v']
	
	v =gobrex.Expand("file{0..2}.jpg")
	fmt.Println(v)	// => ['file0.jpg', 'file1.jpg', 'file2.jpg']
	
	v =gobrex.Expand("file{2..0}.jpg")
	fmt.Println(v)	// => ['file2.jpg', 'file1.jpg', 'file0.jpg']
	
	v =gobrex.Expand("file{0..4..2}.jpg")
	fmt.Println(v)// => ['file0.jpg', 'file2.jpg', 'file4.jpg']
	
	v =gobrex.Expand("file-{a..e..2}.jpg")
	fmt.Println(v) // => ['file-a.jpg', 'file-c.jpg', 'file-e.jpg']
	
	v =gobrex.Expand("file{00..10..5}.jpg")
	fmt.Println(v) // => ['file00.jpg', 'file05.jpg', 'file10.jpg']
	
	v =gobrex.Expand("{{A..C},{a..c}}")
	fmt.Println(v) // => ['A', 'B', 'C', 'a', 'b', 'c']
	
	v =gobrex.Expand("ppp{,config,oe{,conf}}")
	fmt.Println(v) // => ['ppp', 'pppconfig', 'pppoe', 'pppoeconf']
}
```