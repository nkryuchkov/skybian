module github.com/skycoin/skybian

// go 1.14 presents problems for fyne so 1.13 is used.
go 1.13

require (
	fyne.io/fyne v1.2.3
	github.com/google/go-github v17.0.0+incompatible
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/rakyll/statik v0.1.7
	github.com/sirupsen/logrus v1.6.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/skycoin/dmsg v0.0.0-20201009105839-3b0627b4ffa6
	github.com/skycoin/skycoin v0.26.0
	github.com/skycoin/skywire v0.2.3-0.20201009112751-1ab5d15d518c
	github.com/stretchr/testify v1.6.1
	nhooyr.io/websocket v1.8.4 // indirect
)
