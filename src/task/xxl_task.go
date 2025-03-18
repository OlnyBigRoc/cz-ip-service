package task

import (
	"cz-ip-service/src/service"
)

type XXLTask struct {
	Controller *Controller `Path:"task"`
}
type Controller struct {
	SearchService *service.SearchService `AW:"true"`
}
