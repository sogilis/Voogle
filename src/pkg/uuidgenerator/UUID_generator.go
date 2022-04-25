package uuidgenerator

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type IUUIDGenerator interface {
	GenerateUuid() (string, error)
	IsValidUUID(string) bool
}

var _ IUUIDGenerator = &uuidGenerator{}

type uuidGenerator struct{}

func NewUuidGenerator() IUUIDGenerator {
	return &uuidGenerator{}
}

func (g *uuidGenerator) GenerateUuid() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Error("Error while creating uuid : ", err)
		return "", err
	}
	return uuid.String(), nil
}

func (g *uuidGenerator) IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
