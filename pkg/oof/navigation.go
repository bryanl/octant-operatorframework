package oof

import (
	"fmt"

	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func HandleNavigation(request *service.NavigationRequest) (navigation.Navigation, error) {
	return navigation.Navigation{
		Title: "Operator Framework",
		Path:  request.GeneratePath(),
		Children: []navigation.Navigation{
			{
				Title:    "Nested Once",
				Path:     request.GeneratePath("nested-once"),
				IconName: "folder",
				Children: []navigation.Navigation{
					{
						Title:    "Nested Twice",
						Path:     request.GeneratePath("nested-once", "nested-twice"),
						IconName: "folder",
					},
				},
			},
		},
		IconName: "cloud",
	}, nil
}

func InitRoutes(router *service.Router) {
	gen := func(name, accessor, requestPath string) component.Component {
		cardBody := component.NewText(fmt.Sprintf("hello from plugin: path %s", requestPath))
		card := component.NewCard(component.TitleFromString(fmt.Sprintf("My Card - %s", name)))
		card.SetBody(cardBody)
		cardList := component.NewCardList(name)
		cardList.AddCard(*card)
		cardList.SetAccessor(accessor)

		return cardList
	}

	router.HandleFunc("*", func(request service.Request) (component.ContentResponse, error) {
		// For each page, generate two tabs with a some content.
		component1 := gen("Tab 1", "tab1", request.Path())
		component2 := gen("Tab 2", "tab2", request.Path())

		contentResponse := component.NewContentResponse(component.TitleFromString("Example"))
		contentResponse.Add(component1, component2)

		return *contentResponse, nil
	})
}
