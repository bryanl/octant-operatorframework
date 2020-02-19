package oof_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	"github.com/bryanl/octant-operatorframework/pkg/oof"
	"github.com/bryanl/octant-operatorframework/pkg/oof/fake"
)

func TestCatalogSourcePrinter_PrintObject(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	ctx := context.Background()

	request := fake.NewMockPrintRequest(controller)

	catalogSource := loadObject(t, "catalogsource.json")
	request.EXPECT().Object().Return(catalogSource).AnyTimes()

	dashboardClient := fake.NewMockDashboard(controller)
	request.EXPECT().DashboardClient().Return(dashboardClient).AnyTimes()

	request.EXPECT().Context().Return(ctx).AnyTimes()

	registryClient := fake.NewMockRegistryClient(controller)
	registryClient.EXPECT().ListPackages(gomock.Any()).Return([]string{"package"}, nil)

	pf := fake.NewMockPortForward(controller)
	pf.EXPECT().ToService(gomock.Any(), "default", "test-source", uint16(50051)).Return(
		"localhost:51515", func() {}, nil)

	setupCSP := func(csp oof.CatalogSourcePrinter) oof.CatalogSourcePrinter {
		csp.RegistryClientFactory = func(string) (oof.RegistryClient, error) {
			return registryClient, nil
		}

		csp.PortForwardFactory = func(dashboard service.Dashboard) (forward oof.PortForward, err error) {
			return pf, nil
		}

		return csp
	}

	csp := oof.NewCatalogSourcePrinter(setupCSP)

	got, err := csp.PrintObject(request)
	require.NoError(t, err)

	expected := plugin.PrintResponse{
		Config: []component.SummarySection{
			{
				Header:  "Image",
				Content: component.NewText("bryanl/opm-test-index:0.3.0"),
			},
			{
				Header:  "Source Type",
				Content: component.NewText("grpc"),
			},
		},
		Status: []component.SummarySection{
			{
				Header:  "Connection State",
				Content: component.NewText("READY"),
			},
		},
		Items: []component.FlexLayoutItem{
			{
				Width: component.WidthFull,
				View: component.NewTableWithRows("Packages", "", oof.CatalogSourcePackageCols, []component.TableRow{
					{
						"Name": component.NewText("package"),
					},
				}),
			},
		},
	}

	requireJSONEq(t, expected, got)
}
