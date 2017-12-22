package pubsub

type Channel string

func (c *Channel) String() string {
	return string(*c)
}
