package yapr

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

var tt = []struct {
	is               string
	expectError      bool
	expectedCommName string
}{
	{
		`1 (systemd) S 0 1 1 0 -1 4194560 160830 5416764 153 2727262 859 696 263875 113557 20 0 1 0 25 172802048 3451 18446744073709551615 94136074661888 94136075575885 140723048252272 0 0 0 671173123 4096 1260 1 0 0 17 4 0 0 0 0 0 94136075972240 94136076292380 94136089485312 140723048259355 140723048259373 140723048259373 140723048259565 0
`,
		false,
		"systemd",
	},
	{
		`1 () S 0 1 1 0 -1 4194560 160830 5416764 153 2727262 859 696 263875 113557 20 0 1 0 25 172802048 3451 18446744073709551615 94136074661888 94136075575885 140723048252272 0 0 0 671173123 4096 1260 1 0 0 17 4 0 0 0 0 0 94136075972240 94136076292380 94136089485312 140723048259355 140723048259373 140723048259373 140723048259565 0
`,
		true,
		"",
	},
	{
		`1 (systemd) S 0 1 1 0 -1 4194560 160830 5416764 153 2727262 859 696 263875 113557 20 0 1 0 25 172802048 3451 18446744073709551615 94136074661888 94136075575885 140723048252272 0 0 0 671173123 4096 1260 1 0 0 17 4 0 0 0 0 0 94136075972240 94136076292380 94136089485312 140723048259355 140723048259373 140723048259373 140723048259565 0`, // Like above but no new line.
		false,
		"systemd",
	},
	{
		"",
		true,
		"",
	},
	{
		`1 (s y s t e m d)  S 0 1 1 0 -1 4194560          144924 4641887 153 2311970 773 627 247606 108670 20 0 1 0 25 172802048 3451 18446744073709551615 94136074661888 94136075575885 140723048252272 0 0 0 671173123 4096 1260 1 0 0 17 7 0 0 0 0 0 94136075972240 94136076292380 94136089485312 140723048259355 140723048259373 140723048259373 140723048259565 0`,
		false,
		"s y s t e m d",
	},
	{
		`1 ( s y s t e m d )  S 0 1 1 0 -1 4194560          144924 4641887 153 2311970 773 627 247606 108670 20 0 1 0 25 172802048 3451 18446744073709551615 94136074661888 94136075575885 140723048252272 0 0 0 671173123 4096 1260 1 0 0 17 7 0 0 0 0 0 94136075972240 94136076292380 94136089485312 140723048259355 140723048259373 140723048259373 140723048259565 0`,
		false,
		" s y s t e m d ",
	},
	{
		`1 (🤦syst🤦emd🤦) S 0 1 1 0 -1 4194560 160830 5416764 153 2727262 859 696 263875 113557 20 0 1 0 25 172802048 3451 18446744073709551615 94136074661888 94136075575885 140723048252272 0 0 0 671173123 4096 1260 1 0 0 17 4 0 0 0 0 0 94136075972240 94136076292380 94136089485312 140723048259355 140723048259373 140723048259373 140723048259565 0
`,
		false,
		`🤦syst🤦emd🤦`,
	},
	{
		`1 (systemd)`,
		true,
		"",
	},
	{
		`1 systemd S 0 1 1 0 -1 4194560 160830 5416764 153 2727262 859 696 263875 113557 20 0 1 0 25 172802048 3451 18446744073709551615 94136074661888 94136075575885 140723048252272 0 0 0 671173123 4096 1260 1 0 0 17 4 0 0 0 0 0 94136075972240 94136076292380 94136089485312 140723048259355 140723048259373 140723048259373 140723048259565 0
`,
		true,
		"",
	},
	{
		`1 systemd) S 0 1 1 0 -1 4194560 160830 5416764 153 2727262 859 696 263875 113557 20 0 1 0 25 172802048 3451 18446744073709551615 94136074661888 94136075575885 140723048252272 0 0 0 671173123 4096 1260 1 0 0 17 4 0 0 0 0 0 94136075972240 94136076292380 94136089485312 140723048259355 140723048259373 140723048259373 140723048259565 0
`,
		true,
		"",
	},
	{
		`1 (systemd S 0 1 1 0 -1 4194560 160830 5416764 153 2727262 859 696 263875 113557 20 0 1 0 25 172802048 3451 18446744073709551615 94136074661888 94136075575885 140723048252272 0 0 0 671173123 4096 1260 1 0 0 17 4 0 0 0 0 0 94136075972240 94136076292380 94136089485312 140723048259355 140723048259373 140723048259373 140723048259565 0
`,
		true,
		"",
	},
	{
		`1 )systemd( S 0 1 1 0 -1 4194560 160830 5416764 153 2727262 859 696 263875 113557 20 0 1 0 25 172802048 3451 18446744073709551615 94136074661888 94136075575885 140723048252272 0 0 0 671173123 4096 1260 1 0 0 17 4 0 0 0 0 0 94136075972240 94136076292380 94136089485312 140723048259355 140723048259373 140723048259373 140723048259565 0
`,
		true,
		"",
	},
	{
		"1(0)R 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9 9",
		false,
		"0",
	},
}

func TestParseStatString(t *testing.T) {
	for _, tc := range tt {
		st, err := ParseStatString(tc.is)
		if (err == nil && tc.expectError) || (err != nil && !tc.expectError) {
			t.Errorf("%v: got error %v but expectError is %v", tc, err, tc.expectError)
		}

		if err != nil {
			if st.Comm != tc.expectedCommName {
				t.Errorf("%v: got Stat.Comm: %s, expected %s", tc, st.Comm, tc.expectedCommName)
			}
		}
	}
}

// TODO(marete): Fix the code repetition with a refactored string added up from multiple substrings.
var re = regexp.MustCompile(`^(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]){0,}[\-\+]{0,1}(?:[\d]){1,}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]){0,}\((?:[^\n\x00]){1,}\)(?:[\t\v\f\r\n ]){0,}[^)\n\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1,}[\-\+]{0,1}(?:[\d]){1,}){5}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1,}(?:[\d]){1,}){7}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1,}[\-\+]{0,1}(?:[\d]){1,}){6}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1,}(?:[\d]){1,}){2}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1,}[\-\+]{0,1}(?:[\d]){1,}){1}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1,}(?:[\d]){1,}){13}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1,}[\-\+]{0,1}(?:[\d]){1,}){2}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1,}(?:[\d]){1,}){4}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1,}[\-\+]{0,1}(?:[\d]){1,}){1}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1,}(?:[\d]){1,}){7}(?:[\t\v\f\r\x{0009}\x{000d}\x{0020}\x{0085}\x{00a0}\x{1680}\x{2000}\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000} ]{1,}[\-\+]{0,1}(?:[\d]){1,}){1}[^)]*$`)

func FuzzParseStatString(f *testing.F) {
	for _, tc := range tt {
		f.Add(tc.is)
	}

	f.Fuzz(func(t *testing.T, s string) {
		_, err := ParseStatString(s)

		if err == nil && !re.MatchString(s) {
			t.Errorf("expected non-nil error because input %s fails to match regexp %v",
				s, re)
		}

		if err != nil && re.MatchString(s) && !errors.Is(err, strconv.ErrRange) && strings.Index(err.Error(), "integer overflow on token") == -1 {
			t.Errorf("got error: %v, expected nil error because input %s matches regexp %v",
				err, s, re)
		}

	})
}
