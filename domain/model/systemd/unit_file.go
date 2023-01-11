package systemd

import (
	"bytes"
	"reflect"
	"strings"
	"systemd-cd/domain/model/toml"
)

type (
	UnitFileService struct {
		Unit    UnitDirective
		Service ServiceDirective
		Install InstallDirective
	}

	UnitDirective struct {
		Description   string
		Documentation string
		After         []string
		Requires      []string
		Wants         []string
		Conflicts     []string
	}

	UnitType string

	ServiceDirective struct {
		Type            *UnitType
		EnvironmentFile *string
		ExecStart       string
		ExecStop        *string
		ExecReload      *string
		Restart         *string
		RemainAfterExit *string
	}

	InstallDirective struct {
		Alias           []string
		RequiredBy      []string
		WantedBy        []string
		Also            []string
		DefaultInstance *string
	}
)

const (
	UnitTypeSimple  UnitType = "simple"
	UnitTypeForking UnitType = "forking"
	UnitTypeOneShot UnitType = "oneshot"
	UnitTypeDbus    UnitType = "dbus"
	UnitTypeNotify  UnitType = "notify"
	UnitTypeIdle    UnitType = "idle"
)

func (c UnitFileService) Equals(d UnitFileService) bool {
	return reflect.DeepEqual(c, d)
}

type (
	unitFileServiceToml struct {
		Unit    unitDirectiveToml
		Service serviceDirectiveToml
		Install installDirectiveToml
	}

	unitDirectiveToml struct {
		Description   string  `toml:"Description"`
		Documentation string  `toml:"Documentation"`
		After         *string `toml:"After,omitempty"`
		Requires      *string `toml:"Requires,omitempty"`
		Wants         *string `toml:"Wants,omitempty"`
		Conflicts     *string `toml:"Conflicts,omitempty"`
	}

	serviceDirectiveToml struct {
		Type            *UnitType `toml:"Type,omitempty"`
		EnvironmentFile *string   `toml:"EnvironmentFile,omitempty"`
		ExecStart       string    `toml:"ExecStart"`
		ExecStop        *string   `toml:"ExecStop,omitempty"`
		ExecReload      *string   `toml:"ExecReload,omitempty"`
		Restart         *string   `toml:"Restart,omitempty"`
		RemainAfterExit *string   `toml:"RemainAfterExit,omitempty"`
	}

	installDirectiveToml struct {
		Alias           *string `toml:"Alias,omitempty"`
		RequiredBy      *string `toml:"RequiredBy,omitempty"`
		WantedBy        *string `toml:"WantedBy,omitempty"`
		Also            *string `toml:"Also,omitempty"`
		DefaultInstance *string `toml:"DefaultInstance,omitempty"`
	}
)

func MarshalUnitFile(u UnitFileService) ([]byte, error) {
	ut := unitFileServiceToml{
		Unit: unitDirectiveToml{
			Description:   u.Unit.Description,
			Documentation: u.Unit.Documentation,
			After:         spacedString(u.Unit.After),
			Requires:      spacedString(u.Unit.Requires),
			Wants:         spacedString(u.Unit.Wants),
			Conflicts:     spacedString(u.Unit.Conflicts),
		},
		Service: serviceDirectiveToml{
			Type:            u.Service.Type,
			EnvironmentFile: u.Service.EnvironmentFile,
			ExecStart:       u.Service.ExecStart,
			ExecStop:        u.Service.ExecStop,
			ExecReload:      u.Service.ExecReload,
			Restart:         u.Service.Restart,
			RemainAfterExit: u.Service.RemainAfterExit,
		},
		Install: installDirectiveToml{
			Alias:           spacedString(u.Install.Alias),
			RequiredBy:      spacedString(u.Install.RequiredBy),
			WantedBy:        spacedString(u.Install.WantedBy),
			Also:            spacedString(u.Install.Also),
			DefaultInstance: u.Install.DefaultInstance,
		},
	}

	// Encode to toml
	b := &bytes.Buffer{}
	indent := ""
	err := toml.Encode(b, ut, toml.EncodeOption{
		Indent: &indent,
	})
	if err != nil {
		return nil, err
	}
	b.WriteString("\n")

	// Convert to UnitFile format
	s := strings.ReplaceAll(b.String(), " = \"", "=")
	s = strings.ReplaceAll(s, "\"\n", "\n")

	return []byte(s), nil
}

func UnmarshalUnitFile(b *bytes.Buffer) (u UnitFileService, err error) {
	// Convert to toml format
	b2 := &bytes.Buffer{}
	for _, l := range strings.Split(b.String(), "\n") {
		sp := strings.Split(l, "=")
		if len(sp) < 2 {
			b2.WriteString(strings.Join([]string{sp[0], "\n"}, ""))
			continue
		}
		s := []string{sp[0], " = \""}
		s = append(s, sp[1:]...)
		s = append(s, "\"\n")
		b2.WriteString(strings.Join(s, ""))
	}

	// Decode toml
	ut := &unitFileServiceToml{}
	err = toml.Decode(b2, ut)
	if err != nil {
		return
	}

	u = UnitFileService{
		Unit: UnitDirective{
			Description:   ut.Unit.Description,
			Documentation: ut.Unit.Documentation,
			After:         slice(ut.Unit.After),
			Requires:      slice(ut.Unit.Requires),
			Wants:         slice(ut.Unit.Wants),
			Conflicts:     slice(ut.Unit.Conflicts),
		},
		Service: ServiceDirective{
			Type:            ut.Service.Type,
			EnvironmentFile: ut.Service.EnvironmentFile,
			ExecStart:       ut.Service.ExecStart,
			ExecStop:        ut.Service.ExecStop,
			ExecReload:      ut.Service.ExecReload,
			Restart:         ut.Service.Restart,
			RemainAfterExit: ut.Service.RemainAfterExit,
		},
		Install: InstallDirective{
			Alias:           slice(ut.Install.Alias),
			RequiredBy:      slice(ut.Install.RequiredBy),
			WantedBy:        slice(ut.Install.WantedBy),
			Also:            slice(ut.Install.Also),
			DefaultInstance: ut.Install.DefaultInstance,
		},
	}

	return
}

func slice(s *string) []string {
	if s == nil {
		return nil
	}
	return strings.Split(*s, " ")
}

func spacedString(s []string) *string {
	s2 := strings.Join(s, " ")
	return &s2
}
