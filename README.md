# YAPR

`yapr` (Yet Another Proc Reader) is a set of Go routines to parse
various files in the Linux `/proc` filesystem. It was inspired by [an
observation by Dmitry Vyukov on
Twitter](https://twitter.com/dvyukov/status/1605507498420473857) that
almost all attempts to parse `/proc/[pid]/stat` are buggy.

Therefore the first routine `yapr` provides attempts to read
`/proc/[pid]/stat` robustly. But it almost certainly has bugs too, so
feel free to send me a PR 😊.
