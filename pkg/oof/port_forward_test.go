package oof_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/octant/pkg/plugin/api"
	"github.com/vmware-tanzu/octant/pkg/store"

	"github.com/bryanl/octant-operatorframework/pkg/oof"
	"github.com/bryanl/octant-operatorframework/pkg/oof/fake"
)

func TestDashboardPortForward_ToService(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashboardClient := fake.NewMockDashboard(controller)

	endpoints := loadObject(t, "endpoints.json")
	dashboardClient.EXPECT().Get(gomock.Any(), store.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Endpoints",
		Name:       "test-source",
	}).Return(endpoints, nil)

	pod := loadObject(t, "endpoints-pod.json")
	dashboardClient.EXPECT().Get(gomock.Any(), store.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Pod",
		Name:       "test-source-p89bv",
	}).Return(pod, nil)

	dashboardClient.EXPECT().PortForward(gomock.Any(), api.PortForwardRequest{
		Namespace: "default",
		PodName:   "test-source-p89bv",
		Port:      51515,
	}).Return(api.PortForwardResponse{
		ID:   "12345",
		Port: 50000,
	}, nil)
	dashboardClient.EXPECT().CancelPortForward(gomock.Any(), "12345")

	pf, err := oof.NewDashboardPortForward(dashboardClient)
	require.NoError(t, err)

	ctx := context.Background()

	got, cancel, err := pf.ToService(ctx, "default", "test-source", uint16(51515))
	require.NoError(t, err)

	defer cancel()

	require.Equal(t, "localhost:50000", got)
}
