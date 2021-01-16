package embeddeddraftbox

import "github.com/herb-go/notification/notificationdelivery/notificationqueue"
import "github.com/herb-go/notification-drivers/store/embeddedstore"

//DirectiveName registered direcvive name
const DirectiveName = "embeddeddraftbox"

type Config struct {
	embeddedstore.Config
}

//AppylToPublisher create draft box and add to publisher
func (c *Config) AppylToPublisher(p *notificationqueue.Publisher) error {
	store, err := c.CreateStore()
	if err != nil {
		return err
	}
	p.Draftbox = store
	return nil
}

//Factory draftbox directive factory
var Factory = func(loader func(v interface{}) error) (notificationqueue.Directive, error) {
	c := &Config{}
	err := loader(&c.Config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func init() {
	notificationqueue.Register(DirectiveName, Factory)
}
