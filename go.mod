module github.com/kassisol/twic

go 1.25.7

require (
	github.com/glebarez/sqlite v1.11.0
	github.com/juliengk/go-cert v0.0.0-20180306183847-781fb30cc8dc
	github.com/kassisol/tsa v0.0.0-20190325122521-9633304510c6
	github.com/spf13/cobra v1.10.2
	golang.org/x/term v0.42.0
	gorm.io/gorm v1.31.1
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/juliengk/go-utils v0.0.0-20170323144949-ea868d6f9306 // indirect
	github.com/juliengk/stack v0.0.0-20170807090609-ade92a7733a9 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
)

replace (
	github.com/howeyc/gopass => ./local/github.com/howeyc/gopass
	github.com/juliengk/go-cert => ./local/github.com/juliengk/go-cert
	github.com/juliengk/go-utils => ./local/github.com/juliengk/go-utils
	github.com/juliengk/stack => ./local/github.com/juliengk/stack
	github.com/kassisol/tsa => ./local/github.com/kassisol/tsa
)
