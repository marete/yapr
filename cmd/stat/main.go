// The stat command contains some print debugging for the
// ParseStatBytes and similar routines. It used by copying the
// compiled `stat' binary to various strange filenames after which its
// output is compared to the output of `cat /proc/[pid]/stat' while
// the command sleeps.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/marete/yapr"
)

func main() {
	pids := strconv.FormatInt(int64(os.Getpid()), 10)

	rc, err := os.Open("/proc/" + pids + "/stat")
	if err != nil {
		panic(err)
	}

	defer rc.Close()

	b, err := io.ReadAll(rc)
	if err != nil {
		panic(err)
	}

	s, err := yapr.ParseStatBytes(b)
	if err != nil {
		fmt.Printf("%+v\n", s)
		panic(err)
	}

	fmt.Printf("%+v\n", s)

	fmt.Println()

	b, err = json.Marshal(s)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))

	<-time.After(10 * time.Minute)
}
