package wordpress

import (
	"fmt"
	"os/exec"

	"github.com/juju/gocharm/hook"
	"gopkg.in/juju/charm.v5"
	//"github.com/vtolstov/gocharm/charmbits/ospackage"
)

type wordpressctxt struct {
	ctxt *hook.Context
}

func (n *wordpressctxt) setContext(ctxt *hook.Context) error {
	n.ctxt = ctxt
	return nil
}

func (n *wordpressctxt) hook() error {
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

func (n *wordpressctxt) install() error {
	return nil
}

func (n *wordpressctxt) start() error {
	return nil
}

func (n *wordpressctxt) stop() error {
	return nil
}

func (n *wordpressctxt) config() error {
	return nil
}

func (n *wordpressctxt) upgrade() error {
	exec.Command("apt-get", "update").Run()
	return n.install()
}

func RegisterHooks(r *hook.Registry) {
	var n wordpressctxt
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

	r.RegisterConfig("engine", charm.Option{
		Type:        "string",
		Description: "web server engine [ nginx, apache ]",
		Default:     "apache",
	})

	r.RegisterConfig("domain", charm.Option{
		Type:        "string",
		Description: "domain name",
		Default:     "localhost",
	})

	r.RegisterConfig("memory", charm.Option{
		Type:        "string",
		Description: "memory usage",
		Default:     "70%",
	})

	r.RegisterRelation(charm.Relation{Name: "db", Interface: "mysql", Role: charm.RoleRequirer, Scope: charm.ScopeGlobal})
	r.RegisterRelation(charm.Relation{Name: "cache", Interface: "memcache", Role: charm.RoleRequirer, Scope: charm.ScopeGlobal})
	r.RegisterRelation(charm.Relation{Name: "mail", Interface: "smtp", Role: charm.RoleRequirer, Scope: charm.ScopeGlobal})
	r.RegisterRelation(charm.Relation{Name: "website", Interface: "http", Role: charm.RoleRequirer, Scope: charm.ScopeGlobal})
	r.RegisterRelation(charm.Relation{Name: "loadbalancer", Interface: "reverseproxy", Role: charm.RolePeer, Scope: charm.ScopeGlobal})

}
