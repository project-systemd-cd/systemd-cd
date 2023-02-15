package pipeline

import "github.com/google/uuid"

func UUID() string {
	id := uuid.NewString()
	for id == "" {
		id = uuid.NewString()
	}
	return id
}
