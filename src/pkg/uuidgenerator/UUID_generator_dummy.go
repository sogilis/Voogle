package uuidgenerator

var _ IUUIDGenerator = &uuidGeneratorDummy{}

type uuidGeneratorDummy struct {
	generateUuid func() (string, error)
}

func NewUuidGeneratorDummy(generateUuid func() (string, error)) IUUIDGenerator {
	return &uuidGeneratorDummy{generateUuid}
}

func (g *uuidGeneratorDummy) GenerateUuid() (string, error) {
	if g.generateUuid != nil {
		return g.generateUuid()
	}
	return "", nil
}
