package ui

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"

	"github.com/alex-vee-sh/veessh/internal/config"
)

// PickProfile shows an interactive menu to choose a profile. Optional filters can be applied.
func PickProfile(cfg config.Config, protocolFilter string, groupFilter string) (config.Profile, error) {
	profiles := cfg.ListProfiles()
	var filtered []config.Profile
	for _, p := range profiles {
		if protocolFilter != "" && strings.ToLower(string(p.Protocol)) != strings.ToLower(protocolFilter) {
			continue
		}
		if groupFilter != "" && strings.ToLower(p.Group) != strings.ToLower(groupFilter) {
			continue
		}
		filtered = append(filtered, p)
	}
	if len(filtered) == 0 {
		return config.Profile{}, fmt.Errorf("no profiles found")
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].Group == filtered[j].Group {
			return filtered[i].Name < filtered[j].Name
		}
		return filtered[i].Group < filtered[j].Group
	})

	labels := make([]string, 0, len(filtered))
	for _, p := range filtered {
		group := p.Group
		if group == "" {
			group = "default"
		}
		userHost := p.Host
		if p.Username != "" {
			userHost = p.Username + "@" + userHost
		}
		desc := p.Description
		if desc != "" {
			desc = " - " + desc
		}
		labels = append(labels, fmt.Sprintf("%s/%s  (%s)  %s:%d%s", group, p.Name, p.Protocol, userHost, effectivePort(p), desc))
	}

	var selected string
	prompt := &survey.Select{
		Message:  "Select profile:",
		Options:  labels,
		PageSize: 15,
		VimMode:  false,
		Filter: func(filter string, opt string, idx int) bool {
			f := strings.ToLower(filter)
			o := strings.ToLower(opt)
			return strings.Contains(o, f)
		},
	}
	if err := survey.AskOne(prompt, &selected); err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			return config.Profile{}, context.Canceled
		}
		return config.Profile{}, err
	}
	// find selected index
	var idxSel int = -1
	for i, l := range labels {
		if l == selected {
			idxSel = i
			break
		}
	}
	if idxSel < 0 {
		return config.Profile{}, fmt.Errorf("selection not found")
	}
	return filtered[idxSel], nil
}

func effectivePort(p config.Profile) int {
	if p.Port > 0 {
		return p.Port
	}
	switch p.Protocol {
	case config.ProtocolSSH, config.ProtocolSFTP:
		return 22
	case config.ProtocolTelnet:
		return 23
	default:
		return 0
	}
}
