package oof

import (
	"fmt"
	"path"

	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type SubscriptionPrinter struct{}

var _ ObjectPrinter = &SubscriptionPrinter{}

func NewSubscriptionPrinter() *SubscriptionPrinter {
	s := SubscriptionPrinter{}

	return &s
}

func (s *SubscriptionPrinter) PrintObject(request PrintRequest) (plugin.PrintResponse, error) {
	if request == nil {
		return emptyResponse, fmt.Errorf("unable to print a nil object")
	}

	object, err := toUnstructured(request.Object())
	if err != nil {
		return emptyResponse, fmt.Errorf("invalid object: %w", err)
	}

	var subscription v1alpha1.Subscription
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, &subscription); err != nil {
		return emptyResponse, err
	}

	response := plugin.PrintResponse{}

	configuration, err := s.printConfiguration(subscription)
	if err != nil {
		return emptyResponse, fmt.Errorf("print configuration %w", err)
	}
	response.Config = configuration

	status, err := s.printStatus(subscription)
	if err != nil {
		return emptyResponse, fmt.Errorf("print status %w", err)
	}
	response.Status = status

	conditionsTable, err := s.printConditions(subscription)
	if err != nil {
		return emptyResponse, err
	}

	response.Items = append(response.Items, component.FlexLayoutItem{
		Width: component.WidthFull,
		View:  conditionsTable,
	})

	return response, nil
}

var (
	subscriptionPhases = map[v1alpha1.SubscriptionState]component.TextStatus{
		"UpgradeAvailable": component.TextStatusWarning,
		"UpgradePending":   component.TextStatusWarning,
		"AtLatestKnown":    component.TextStatusOK,
	}
)

func (s *SubscriptionPrinter) printConfiguration(subscription v1alpha1.Subscription) ([]component.SummarySection, error) {
	var sections []component.SummarySection

	// print source
	sourceName := subscription.Spec.CatalogSource
	sourceNamespace := subscription.Spec.CatalogSourceNamespace
	sourceLink := component.NewLink("", sourceName, path.Join("/overview/namespace", sourceNamespace,
		"custom-resources", "catalogsources.operators.coreos.com/v1alpha1", sourceName))
	sections = append(sections, component.SummarySection{
		Header:  "Source",
		Content: sourceLink,
	})

	return sections, nil
}

func (s *SubscriptionPrinter) printStatus(subscription v1alpha1.Subscription) ([]component.SummarySection, error) {
	var sections []component.SummarySection

	state := subscription.Status.State

	status, ok := subscriptionPhases[state]
	if !ok {
		status = component.TextStatusError
	}

	statusView := component.NewText(string(state))
	statusView.SetStatus(status)

	sections = append(sections, component.SummarySection{
		Header:  "Status",
		Content: statusView,
	})

	return sections, nil
}

var (
	SubscriptionConditionColumns = component.NewTableCols("Type", "Last Transition Time", "Reason", "Message")
)

func (s *SubscriptionPrinter) printConditions(subscription v1alpha1.Subscription) (*component.Table, error) {
	table := component.NewTable("Conditions", "", SubscriptionConditionColumns)

	for _, condition := range subscription.Status.Conditions {
		t := component.NewText(string(condition.Type))
		switch condition.Status {
		case v1.ConditionTrue:
			t.SetStatus(component.TextStatusError)
		case v1.ConditionFalse:
			t.SetStatus(component.TextStatusOK)
		default:
			t.SetStatus(component.TextStatusWarning)
		}

		row := component.TableRow{
			"Type":                 t,
			"Last Transition Time": component.NewTimestamp(condition.LastTransitionTime.Time),
			"Reason":               component.NewText(condition.Reason),
			"Message":              component.NewText(condition.Message),
		}

		table.Add(row)
	}

	return table, nil
}
