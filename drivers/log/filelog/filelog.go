package filelog

import (
	"fmt"

	filelog_transporter "github.com/eolinker/eosc/log/filelog"
	log_transport "github.com/eolinker/goku/log-transport"

	"github.com/eolinker/eosc"
	transporterManager "github.com/eolinker/eosc/log/transporter-manager"
)

type filelog struct {
	id                 string
	name               string
	config             *filelog_transporter.Config
	formatterName      string
	transporterReset   log_transport.TransporterReset
	transporterManager transporterManager.ITransporterManager
}

func (f *filelog) Id() string {
	return f.id
}

func (f *filelog) Start() error {

	formatter, err := CreateFormatter(f.formatterName)
	if err != nil {
		return err
	}

	transporterReset := filelog_transporter.CreateTransporter(f.config.Level)
	err = transporterReset.Reset(f.config, formatter)
	if err != nil {
		return err
	}

	f.transporterReset = transporterReset
	return f.transporterManager.Set(f.id, transporterReset)
}

func (f *filelog) Reset(conf interface{}, workers map[eosc.RequireId]interface{}) error {
	config, ok := conf.(*DriverConfig)
	if !ok {
		return fmt.Errorf("need %s,now %s", eosc.TypeNameOf((*DriverConfig)(nil)), eosc.TypeNameOf(conf))
	}

	c, err := toConfig(config)
	if err != nil {
		return err
	}
	f.config = c
	f.formatterName = config.FormatterName

	formatter, err := CreateFormatter(f.formatterName)
	if err != nil {
		return err
	}

	err = f.transporterReset.Reset(c, formatter)
	if err != nil {
		return err
	}

	return nil
}

func (f *filelog) Stop() error {
	err := f.transporterReset.Close()
	if err != nil {
		return err
	}

	return f.transporterManager.Del(f.id)
}

func (f *filelog) CheckSkill(skill string) bool {
	return false
}