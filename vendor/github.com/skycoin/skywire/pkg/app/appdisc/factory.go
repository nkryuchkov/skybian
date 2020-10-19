package appdisc

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/skycoin/dmsg/cipher"
	"github.com/skycoin/skycoin/src/util/logging"

	"github.com/skycoin/skywire/pkg/app/appcommon"
	"github.com/skycoin/skywire/pkg/servicedisc"
	"github.com/skycoin/skywire/pkg/skyenv"
)

// Factory creates appdisc.Updater instances based on the app name.
type Factory struct {
	Log            logrus.FieldLogger
	PK             cipher.PubKey
	SK             cipher.SecKey
	UpdateInterval time.Duration
	ProxyDisc      string // Address of proxy-discovery
}

func (f *Factory) setDefaults() {
	if f.Log == nil {
		f.Log = logging.MustGetLogger("appdisc")
	}
	if f.UpdateInterval == 0 {
		f.UpdateInterval = skyenv.AppDiscUpdateInterval
	}
	if f.ProxyDisc == "" {
		f.ProxyDisc = skyenv.DefaultServiceDiscAddr
	}
}

// VisorUpdater obtains a visor updater.
func (f *Factory) VisorUpdater(port uint16) Updater {
	// Always return empty updater if keys are not set.
	if f.setDefaults(); f.PK.Null() || f.SK.Null() {
		return &emptyUpdater{}
	}

	conf := servicedisc.Config{
		Type:     servicedisc.ServiceTypeVisor,
		PK:       f.PK,
		SK:       f.SK,
		Port:     port,
		DiscAddr: f.ProxyDisc,
	}

	return &serviceUpdater{
		client:   servicedisc.NewClient(f.Log, conf),
		interval: f.UpdateInterval,
	}
}

// AppUpdater obtains an app updater based on the app name and configuration.
func (f *Factory) AppUpdater(conf appcommon.ProcConfig) (Updater, bool) {
	// Always return empty updater if keys are not set.
	if f.setDefaults(); f.PK.Null() || f.SK.Null() {
		return &emptyUpdater{}, false
	}

	log := f.Log.WithField("appName", conf.AppName)

	// Do not update in proxy discovery if passcode-protected.
	if containsFlag(conf.ProcArgs, "passcode") {
		return &emptyUpdater{}, false
	}

	getServiceDiscConf := func(conf appcommon.ProcConfig, sType string) servicedisc.Config {
		return servicedisc.Config{
			Type:     sType,
			PK:       f.PK,
			SK:       f.SK,
			Port:     uint16(conf.RoutingPort),
			DiscAddr: f.ProxyDisc,
		}
	}

	switch conf.AppName {
	case skyenv.SkysocksName:
		return &serviceUpdater{
			client:   servicedisc.NewClient(log, getServiceDiscConf(conf, servicedisc.ServiceTypeProxy)),
			interval: f.UpdateInterval,
		}, true
	case skyenv.VPNServerName:
		return &serviceUpdater{
			client:   servicedisc.NewClient(log, getServiceDiscConf(conf, servicedisc.ServiceTypeVPN)),
			interval: f.UpdateInterval,
		}, true
	default:
		return &emptyUpdater{}, false
	}
}

func containsFlag(args []string, flag string) bool {
	for _, arg := range args {
		if argEqualsFlag(arg, flag) {
			return true
		}
	}
	return false
}

func argEqualsFlag(arg, flag string) bool {
	arg = strings.TrimSpace(arg)

	// strip prefixed '-'s.
	for {
		if len(arg) < 1 {
			return false
		}
		if arg[0] == '-' {
			arg = arg[1:]
			continue
		}
		break
	}

	// strip anything after (inclusive) of '='.
	arg = strings.Split(arg, "=")[0]

	return arg == flag
}
