package kvreceiptstore

import (
	"github.com/herb-go/herbdata/kvdb"
)

//StoreConfig key-value receipt store config
type Config struct {
	//kvdb config
	Database *kvdb.Config
	//Limit count limit,defalut value is notificationquque.DefaultStoreListLimit
	Limit int
	//RretentionDays data rretention in days
	RretentionDays int
}

//CreateStore create store
func (c *Config) CreateStore() (*Store, error) {
	s := New()
	s.DB = kvdb.New()
	err := c.Database.ApplyTo(s.DB)
	if err != nil {
		return nil, err
	}
	s.Limit = c.Limit
	s.DataRetentionDays = c.RretentionDays
	return s, nil
}
