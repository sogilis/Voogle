package uuidgenerator

var _ IUUIDGenerator = &uuidGeneratorDummy{}

type uuidGeneratorDummy struct {
	generateUuid func() (string, error)
	isValidUUID  func(string) bool
}

func NewUuidGeneratorDummy(generateUuid func() (string, error), isValidUUID func(string) bool) IUUIDGenerator {
	return &uuidGeneratorDummy{generateUuid, isValidUUID}
}

func (g *uuidGeneratorDummy) GenerateUuid() (string, error) {
	if g.generateUuid != nil {
		return g.generateUuid()
	}
	return "", nil
}

func (g *uuidGeneratorDummy) IsValidUUID(u string) bool {
	if g.isValidUUID != nil {
		return g.isValidUUID(u)
	}
	return true
}
