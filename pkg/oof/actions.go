package oof

import (
	"fmt"

	"github.com/vmware-tanzu/octant/pkg/plugin/service"
)

func HandleActions(request *service.ActionRequest) error {
	if request == nil {
		return fmt.Errorf("request is nil")
	}

	actionName, err := request.Payload.String("action")
	if err != nil {
		return err
	}

	// ctx := context.Background()

	switch actionName {
	case ActionSubscribeToPackage:

	}

	return nil
}
