package mysql

import (
	"fmt"
	"os/exec"

	"github.com/juju/gocharm/hook"
	"gopkg.in/juju/charm.v5"
	//"github.com/vtolstov/gocharm/charmbits/ospackage"
)

type mysqlctxt struct {
	ctxt *hook.Context
}

func (n *mysqlctxt) setContext(ctxt *hook.Context) error {
	n.ctxt = ctxt
	return nil
}

func (n *mysqlctxt) hook() error {
	var err error

	n.ctxt.Logf("hook %s is running", n.ctxt.HookName)

	switch n.ctxt.HookName {
	case "install":
		err = n.install()
	case "start":
		err = n.start()
	case "stop":
		err = n.stop()
	case "upgrade-charm":
		err = n.upgrade()
	case "config-changed":
		err = n.config()
	default:
		err = fmt.Errorf("not implemented")
	}
	return err
}

func (n *mysqlctxt) install() error {
	args := []string{"install", "--no-install-suggests", "--no-install-recommends", "-y"}
	pkgs := []string{"mysql-server"}
	args = append(args, pkgs...)
	cmd := exec.Command("apt-get", args...)
	return cmd.Run()
}

func (n *mysqlctxt) start() error {
	err := exec.Command("service", "mysql", "status").Run()
	if err != nil {
		err = exec.Command("service", "mysql", "start").Run()
	}
	return err
}

func (n *mysqlctxt) stop() error {
	err := exec.Command("service", "mysql", "status").Run()
	if err == nil {
		err = exec.Command("service", "mysql", "stop").Run()
	}
	return err
}

func (n *mysqlctxt) config() error {
	return nil
}

func (n *mysqlctxt) upgrade() error {
	exec.Command("apt-get", "update").Run()
	return n.install()
}

func RegisterHooks(r *hook.Registry) {
	var n mysqlctxt
	r.RegisterContext(n.setContext, nil)
	// Standard hooks
	for _, item := range []string{"install", "config-changed", "start", "upgrade-charm", "stop"} {
		r.RegisterHook(item, n.hook)
	}
	// Relation hooks
	for _, item := range []string{"relation-joined", "relation-changed", "relation-departed", "relation-broken"} {
		r.RegisterHook("db-"+item, n.hook)
	}
	//	r.RegisterHook("*", r.changed)
	r.RegisterConfig("tuning-level", charm.Option{
		Type:        "string",
		Description: "tuning level [default, safe, fast ]",
		Default:     "default",
	})

	r.RegisterConfig("max-connections", charm.Option{
		Type:        "int",
		Description: "max server connections",
		Default:     -1,
	})

	r.RegisterConfig("flavor", charm.Option{
		Type:        "string",
		Description: "mysql flavor to use [ default, mariadb, percona ]",
		Default:     "default",
	})

	r.RegisterConfig("memory", charm.Option{
		Type:        "string",
		Description: "memory usage",
		Default:     "70%",
	})

	r.RegisterRelation(charm.Relation{Name: "db", Interface: "mysql", Role: charm.RoleProvider, Scope: charm.ScopeGlobal})

}
