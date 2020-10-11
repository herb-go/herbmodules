package overseers

import "github.com/herb-go/worker"

//ExtractorFactoryOverseerConfig overseer config struct
type ExtractorFactoryOverseerConfig struct {
}

//ApplyTo apply config to overseer
func (c *ExtractorFactoryOverseerConfig) ApplyTo(o *worker.PlainOverseer) error {
	o.WithIntroduction("httpinfo extractor factory overseer workers")
	return nil
}

//NewExtractorFactoryOverseerConfig create new reader overseer config
func NewExtractorFactoryOverseerConfig() *ExtractorFactoryOverseerConfig {
	return &ExtractorFactoryOverseerConfig{}
}

//ExtractorOverseerConfig overseer config struct
type ExtractorOverseerConfig struct {
}

//ApplyTo apply config to overseer
func (c *ExtractorOverseerConfig) ApplyTo(o *worker.PlainOverseer) error {
	o.WithIntroduction("httpinfo extractor overseer workers")
	return nil
}

//NewExtractorOverseerConfig create new reader overseer config
func NewExtractorOverseerConfig() *ExtractorOverseerConfig {
	return &ExtractorOverseerConfig{}
}

//FormatterFactoryOverseerConfig overseer config struct
type FormatterFactoryOverseerConfig struct {
}

//ApplyTo apply config to overseer
func (c *FormatterFactoryOverseerConfig) ApplyTo(o *worker.PlainOverseer) error {
	o.WithIntroduction("httpinfo formatter factory overseer workers")
	return nil
}

//NewFormatterFactoryOverseerConfig create new reader overseer config
func NewFormatterFactoryOverseerConfig() *FormatterFactoryOverseerConfig {
	return &FormatterFactoryOverseerConfig{}
}

//FormatterOverseerConfig overseer config struct
type FormatterOverseerConfig struct {
}

//ApplyTo apply config to overseer
func (c *FormatterOverseerConfig) ApplyTo(o *worker.PlainOverseer) error {
	o.WithIntroduction("httpinfo formatter overseer workers")
	return nil
}

//NewFormatterOverseerConfig create new reader overseer config
func NewFormatterOverseerConfig() *FormatterOverseerConfig {
	return &FormatterOverseerConfig{}
}
