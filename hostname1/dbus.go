// Integration with the systemd hostnamed API.  See https://www.freedesktop.org/wiki/Software/systemd/hostnamed/
package hostname1

import (
	"os"
	"strconv"

	"github.com/godbus/dbus"
)

const (
	dbusInterface = "org.freedesktop.hostname1"
	dbusPath      = "/org/freedesktop/hostname1"
)

// Conn is a connection to systemds dbus endpoint.
type Conn struct {
	conn   *dbus.Conn
	object dbus.BusObject
}

// New() establishes a connection to the system bus and authenticates.
func New() (*Conn, error) {
	c := new(Conn)

	if err := c.initConnection(); err != nil {
		return nil, err
	}

	return c, nil
}

// Close closes the dbus connection
func (c *Conn) Close() {
	if c == nil {
		return
	}

	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Conn) initConnection() error {
	var err error
	c.conn, err = dbus.SystemBusPrivate()
	if err != nil {
		return err
	}

	// Only use EXTERNAL method, and hardcode the uid (not username)
	// to avoid a username lookup (which requires a dynamically linked
	// libc)
	methods := []dbus.Auth{dbus.AuthExternal(strconv.Itoa(os.Getuid()))}

	err = c.conn.Auth(methods)
	if err != nil {
		c.conn.Close()
		return err
	}

	err = c.conn.Hello()
	if err != nil {
		c.conn.Close()
		return err
	}

	c.object = c.conn.Object("org.freedesktop.hostname1", dbus.ObjectPath(dbusPath))

	return nil
}

// getProperties returns all of the dbus object properties
func (c *Conn) GetProperties() (map[string]interface{}, error) {
	var err error
	var props map[string]dbus.Variant

	err = c.object.Call("org.freedesktop.DBus.Properties.GetAll", 0, dbusInterface).Store(&props)
	if err != nil {
		return nil, err
	}

	out := make(map[string]interface{}, len(props))
	for k, v := range props {
		out[k] = v.Value()
	}

	return out, nil
}
