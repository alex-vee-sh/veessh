package connectors

import (
	"context"
	"os/exec"
	"strconv"

	"github.com/alex-vee-sh/veessh/internal/config"
	"github.com/alex-vee-sh/veessh/internal/util"
)

type sshConnector struct{}

func (s *sshConnector) Name() string { return "ssh" }

func (s *sshConnector) Exec(ctx context.Context, p config.Profile, _ string) error {
	args := []string{}
	if p.Port > 0 {
		args = append(args, "-p", strconv.Itoa(p.Port))
	}
	if p.Username != "" {
		args = append(args, "-l", p.Username)
	}
	if p.IdentityFile != "" {
		args = append(args, "-i", p.IdentityFile)
	}
	if p.ProxyJump != "" {
		args = append(args, "-J", p.ProxyJump)
	}
	for _, lf := range p.LocalForwards {
		if lf != "" {
			args = append(args, "-L", lf)
		}
	}
	for _, rf := range p.RemoteForwards {
		if rf != "" {
			args = append(args, "-R", rf)
		}
	}
	for _, df := range p.DynamicForwards {
		if df != "" {
			args = append(args, "-D", df)
		}
	}
	if len(p.ExtraArgs) > 0 {
		args = append(args, p.ExtraArgs...)
	}
	args = append(args, p.Host)
	cmd := exec.CommandContext(ctx, "ssh", args...)
	return util.RunAttached(cmd)
}

func init() {
	Register(config.ProtocolSSH, &sshConnector{})
}
