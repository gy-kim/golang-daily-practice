package advantages

type CarV2 struct {
	engine Engine
}

func (c *CarV2) Drive() error {
	// use the engine
	c.engine.Start()
	c.engine.IncreasePower()

	return nil
}

func (c *CarV2) Stop() error {
	// use the engine
	c.engine.DecreasePower()
	c.engine.Stop()

	return nil
}
