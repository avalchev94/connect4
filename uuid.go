package tarantula

import "github.com/google/uuid"

type UUID string

func NewUUID() (UUID, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return UUID(uuid.String()), nil
}
