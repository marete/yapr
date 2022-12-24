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

	// st has > 1 space before the status byte.
	st := `1 (s y s t e m d)  S 0 1 1 0 -1 4194560 144924 4641887 153 2311970 773 627 247606 108670 20 0 1 0 25 172802048 3451 18446744073709551615 94136074661888 94136075575885 140723048252272 0 0 0 671173123 4096 1260 1 0 0 17 7 0 0 0 0 0 94136075972240 94136076292380 94136089485312 140723048259355 140723048259373 140723048259373 140723048259565 0`
	s, err = yapr.ParseStatString(st)
	if err != nil {
		panic(err)
	}
	fmt.Println()

	fmt.Printf("%+v\n", s)

	<-time.After(10 * time.Minute)
}
