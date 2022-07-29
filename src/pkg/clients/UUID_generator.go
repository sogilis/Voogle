package clients

import (
	"errors"
	"strconv"
	"time"

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

	uuidString := uuid.String()
	if len(uuidString) == 0 {
		return uuidString, errors.New("invalid uuid")
	}

	// Use minute since EPOCH as hexa in uuid (COMB uuid)
	unixMinuteHexa := strconv.FormatInt(time.Now().Unix()/int64(60), 16)
	uuidString = unixMinuteHexa[len(unixMinuteHexa)-4:] + uuidString[4:]
	return uuidString, nil
}

func (g *uuidGenerator) IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
