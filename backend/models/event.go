package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/zyghq/zyg/utils"
)

// EventSeverity Represents the event severity
type EventSeverity string

// Predefined event severity.
const (
	SeverityInfo     EventSeverity = "info"
	SeverityCritical EventSeverity = "critical"
	SeverityError    EventSeverity = "error"
	SeverityMuted    EventSeverity = "muted"
	SeveritySuccess  EventSeverity = "success"
	SeverityWarning  EventSeverity = "warning"
)

// IsValid checks if the Severity value is one of the predefined constants:
// SeverityInfo, SeverityCritical, SeverityError, SeverityMuted, SeveritySuccess, or SeverityWarning.
// Returns true if the value is valid, false otherwise.
func (ev EventSeverity) IsValid() bool {
	switch ev {
	case SeverityInfo, SeverityCritical, SeverityError, SeverityMuted, SeveritySuccess, SeverityWarning:
		return true
	default:
		return false
	}
}

func (ev EventSeverity) String() string {
	return string(ev)
}

// TextSize Represents text size for the text component.
type TextSize string

const (
	TextSizeL  TextSize = "L"
	TextSizeM  TextSize = "M"
	TextSizeS  TextSize = "S"
	TextSizeXS TextSize = "XS"
)

// IsValid checks if the TextSize value is one of the predefined constants:
// TextSizeL ("L"), TextSizeM ("M"), TextSizeS ("S"), or TextSizeXS ("XS").
// Returns true if the value is valid, false otherwise.
func (ts TextSize) IsValid() bool {
	switch ts {
	case TextSizeL, TextSizeM, TextSizeS, TextSizeXS:
		return true
	default:
		return false
	}
}

// TextColor Represents text color for the text component.
type TextColor string

const (
	TextNormal  TextColor = "NORMAL"
	TextMuted   TextColor = "MUTED"
	TextWarning TextColor = "WARNING"
	TextError   TextColor = "ERROR"
	TextSuccess TextColor = "SUCCESS"
)

func (tc TextColor) IsValid() bool {
	switch tc {
	case TextNormal, TextMuted, TextWarning, TextError, TextSuccess:
		return true
	default:
		return false
	}
}

type SpacerSize string

const (
	SpacerSizeL  SpacerSize = "L"
	SpacerSizeM  SpacerSize = "M"
	SpacerSizeS  SpacerSize = "S"
	SpacerSizeXS SpacerSize = "XS"
)

// IsValid checks if the SpacerSize value is one of the predefined constants:
// SpacerSizeL ("L"), SpacerSizeM ("M"), SpacerSizeS ("S"), or SpacerSizeXS ("XS").
// Returns true if the value is valid, false otherwise.
func (ss SpacerSize) IsValid() bool {
	switch ss {
	case SpacerSizeL, SpacerSizeM, SpacerSizeS, SpacerSizeXS:
		return true
	default:
		return false
	}
}

type DividerSize string

const (
	DividerSizeL  DividerSize = "L"
	DividerSizeM  DividerSize = "M"
	DividerSizeS  DividerSize = "S"
	DividerSizeXS DividerSize = "XS"
)

// IsValid checks if the DividerSize value is one of the predefined constants:
// DividerSizeL ("L"), DividerSizeM ("M"), DividerSizeS ("S"), or DividerSizeXS ("XS").
// Returns true if the value is valid, false otherwise.
func (ds DividerSize) IsValid() bool {
	switch ds {
	case DividerSizeL, DividerSizeM, DividerSizeS, DividerSizeXS:
		return true
	default:
		return false
	}
}

type BadgeColor string

const (
	BadgeColorRed    BadgeColor = "RED"
	BadgeColorGreen  BadgeColor = "GREEN"
	BadgeColorBlue   BadgeColor = "BLUE"
	BadgeColorGray   BadgeColor = "GRAY"
	BadgeColorYellow BadgeColor = "YELLOW"
)

// IsValid checks if the BadgeColor value is one of the predefined constants:
// BadgeColorRed ("RED"), BadgeColorGreen ("GREEN"), BadgeColorBlue ("BLUE"),
// BadgeColorGray ("GRAY"), or BadgeColorYellow ("YELLOW").
// Returns true if the value is valid, false otherwise.
func (bc BadgeColor) IsValid() bool {
	switch bc {
	case BadgeColorRed, BadgeColorGreen, BadgeColorBlue, BadgeColorGray, BadgeColorYellow:
		return true
	default:
		return false
	}
}

// ComponentText represents the component text in event components.
type ComponentText struct {
	Text      string    `json:"text"`
	TextSize  TextSize  `json:"textSize"`
	TextColor TextColor `json:"textColor"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for ComponentText.
// It unmarshals the JSON data and validates the textSize field.
// If textSize is empty, it defaults to 'S'. Otherwise, checks if the value
// is one of the valid predefined sizes (L, M, S, XS).
// Returns an error if JSON unmarshal fails or if textSize is invalid.
func (ct *ComponentText) UnmarshalJSON(data []byte) error {
	type Aux struct {
		Text      string `json:"text"`
		TextSize  string `json:"textSize"`
		TextColor string `json:"textColor"`
	}

	var aux Aux
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	ct.Text = aux.Text
	// Set default value to 'S' if textSize is empty
	if aux.TextSize == "" {
		ct.TextSize = TextSizeS
	} else {
		ts := TextSize(aux.TextSize)
		if !ts.IsValid() {
			return fmt.Errorf("invalid textSize value: %s", aux.TextSize)
		}
		ct.TextSize = ts
	}
	// Set default value to 'NORMAL' if textColor is empty
	if aux.TextColor == "" {
		ct.TextColor = TextNormal
	} else {
		tc := TextColor(aux.TextColor)
		if !tc.IsValid() {
			return fmt.Errorf("invalid textColor value: %s", aux.TextColor)
		}
		ct.TextColor = tc
	}
	return nil
}

// ComponentSpacer represents the component spacer in event components.
type ComponentSpacer struct {
	SpacerSize SpacerSize `json:"spacerSize"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for ComponentSpacer.
// It unmarshals the JSON data and validates the spacerSize field.
// If spacerSize is empty, it defaults to 'S'. Otherwise, checks if the value
// is one of the valid predefined sizes (L, M, S, XS).
// Returns an error if JSON unmarshal fails or if spacerSize is invalid.
func (cs *ComponentSpacer) UnmarshalJSON(data []byte) error {
	type Aux struct {
		SpacerSize string `json:"spacerSize"`
	}

	var aux Aux
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Set default value to 'S' if spacerSize is empty
	if aux.SpacerSize == "" {
		cs.SpacerSize = SpacerSizeS
	} else {
		ss := SpacerSize(aux.SpacerSize)
		if !ss.IsValid() {
			return fmt.Errorf("invalid spacerSize value: %s", aux.SpacerSize)
		}
		cs.SpacerSize = ss
	}
	return nil
}

// ComponentLinkButton represents the component link button in event components.
type ComponentLinkButton struct {
	LinkButtonLabel string `json:"linkButtonLabel"`
	LinkButtonUrl   string `json:"linkButtonUrl"`
}

// ComponentDivider represents the component divider in event components.
type ComponentDivider struct {
	DividerSize DividerSize `json:"dividerSize"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for ComponentDivider.
// It unmarshals the JSON data and validates the dividerSize field.
// If dividerSize is empty, it defaults to 'S'. Otherwise, checks if the value
// is one of the valid predefined sizes (L, M, S, XS).
// Returns an error if JSON unmarshal fails or if dividerSize is invalid.
func (cd *ComponentDivider) UnmarshalJSON(data []byte) error {
	type Aux struct {
		DividerSize string `json:"dividerSize"`
	}

	var aux Aux
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Set default value to 'S' if dividerSize is empty
	if aux.DividerSize == "" {
		cd.DividerSize = DividerSizeS
	} else {
		ds := DividerSize(aux.DividerSize)
		if !ds.IsValid() {
			return fmt.Errorf("invalid dividerSize value: %s", aux.DividerSize)
		}
		cd.DividerSize = ds
	}
	return nil
}

// ComponentCopyButton represents the copy button components in event components.
type ComponentCopyButton struct {
	CopyButtonToolTipLabel string `json:"copyButtonToolTipLabel"`
	CopyButtonValue        string `json:"copyButtonValue"`
}

// ComponentBadge represents the badge component in event components.
type ComponentBadge struct {
	BadgeColor BadgeColor `json:"badgeColor"`
	BadgeLabel string     `json:"badgeLabel"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for ComponentBadge.
// It unmarshals the JSON data and validates the badgeColor field.
// If badgeColor is empty, it defaults to 'GRAY'. Otherwise, checks if the value
// is one of the valid predefined colors (RED, GREEN, BLUE, GRAY, YELLOW).
// Returns an error if JSON unmarshal fails or if badgeColor is invalid.
func (cb *ComponentBadge) UnmarshalJSON(data []byte) error {
	type Aux struct {
		BadgeColor string `json:"badgeColor"`
		BadgeLabel string `json:"badgeLabel"`
	}

	var aux Aux
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	cb.BadgeLabel = aux.BadgeLabel
	if aux.BadgeColor == "" {
		cb.BadgeColor = BadgeColorGray
	} else {
		bc := BadgeColor(aux.BadgeColor)
		if !bc.IsValid() {
			return fmt.Errorf("invalid badgeColor value: %s", aux.BadgeColor)
		}
		cb.BadgeColor = bc
	}
	return nil
}

// ComponentRow represents row components that can contain main content and aside content.
// Each field is an array of EventComponents, allowing for multiple components to be arranged
// horizontally in either the main or aside sections of the row.
type ComponentRow struct {
	RowMainContent  []EventComponent `json:"rowMainContent"`  // Components in the main content area of the row
	RowAsideContent []EventComponent `json:"rowAsideContent"` // Components in the aside content area of the row
}

// ValidateComponentText checks if the ComponentText has valid text and textSize fields.
// It takes a ComponentText struct and its index in a component array as parameters.
// Returns an error if the text field is empty or if the textSize is not one of the
// predefined values (L, M, S, XS). The index is used in error messages to identify
// which component failed validation.
func ValidateComponentText(ct *ComponentText, index int) error {
	if ct == nil {
		return nil
	}
	if ct.Text == "" {
		return fmt.Errorf("component %d: text field cannot be empty", index)
	}
	if !ct.TextSize.IsValid() {
		return fmt.Errorf(
			"component %d: invalid textSize value %s (expected one of L, M, S, XS)",
			index, ct.TextSize,
		)
	}
	return nil
}

// ValidateComponentSpacer checks if the ComponentSpacer has a valid spacerSize field.
// It takes a ComponentSpacer struct and its index in a component array as parameters.
// Returns an error if the spacerSize is not one of the predefined values (L, M, S, XS).
// The index is used in error messages to identify which component failed validation.
func ValidateComponentSpacer(cs *ComponentSpacer, index int) error {
	if cs == nil {
		return nil
	}
	if !cs.SpacerSize.IsValid() {
		return fmt.Errorf(
			"component %d: invalid spacerSize value %s (expected one of L, M, S, XS)",
			index, cs.SpacerSize,
		)
	}
	return nil
}

// ValidateComponentLinkButton checks if the ComponentLinkButton has valid label and URL fields.
// It takes a ComponentLinkButton struct and its index in a component array as parameters.
// Returns an error if either linkButtonLabel or linkButtonUrl fields are empty.
// The index is used in error messages to identify which component failed validation.
// Returns nil if the input pointer is nil or if validation passes.
func ValidateComponentLinkButton(lb *ComponentLinkButton, index int) error {
	if lb == nil {
		return nil
	}
	if lb.LinkButtonLabel == "" {
		return fmt.Errorf("component %d: linkButtonLabel cannot be empty", index)
	}
	if lb.LinkButtonUrl == "" {
		return fmt.Errorf("component %d: linkButtonUrl cannot be empty", index)
	}
	return nil
}

// ValidateComponentDivider checks if the ComponentDivider has a valid dividerSize field.
// It takes a ComponentDivider struct and its index in a component array as parameters.
// Returns an error if the dividerSize is not one of the predefined values (L, M, S, XS).
// The index is used in error messages to identify which component failed validation.
func ValidateComponentDivider(cd *ComponentDivider, index int) error {
	if cd == nil {
		return nil
	}
	if !cd.DividerSize.IsValid() {
		return fmt.Errorf(
			"component %d: invalid dividerSize value %s (expected one of L, M, S, XS)",
			index, cd.DividerSize,
		)
	}
	return nil
}

// ValidateComponentCopyButton checks if the ComponentCopyButton has valid tooltip label and value fields.
// It takes a ComponentCopyButton struct and its index in a component array as parameters.
// Returns an error if copyButtonValue field is empty. If the copyButtonToolTipLabel is empty,
// it defaults to "copy". The index is used in error messages to identify which component failed validation.
// Returns nil if the input pointer is nil or if validation passes.
func ValidateComponentCopyButton(cp *ComponentCopyButton, index int) error {
	if cp == nil {
		return nil
	}
	if cp.CopyButtonToolTipLabel == "" {
		cp.CopyButtonToolTipLabel = "copy"
	}
	if cp.CopyButtonValue == "" {
		return fmt.Errorf("component %d: copyButtonValue cannot be empty", index)
	}
	return nil
}

// ValidateComponentBadge checks if the ComponentBadge has valid badgeColor and badgeLabel fields.
// It takes a ComponentBadge struct and its index in a component array as parameters.
// Returns an error if the badgeColor is not one of the predefined colors
// (RED, GREEN, BLUE, GRAY, YELLOW) or if the badgeLabel is empty.
// The index is used in error messages to identify which component failed validation.
// Returns nil if the input pointer is nil or if validation passes.
func ValidateComponentBadge(cb *ComponentBadge, index int) error {
	if cb == nil {
		return nil
	}
	if !cb.BadgeColor.IsValid() {
		return fmt.Errorf(
			"component %d: invalid badgeColor value %s (expected one of RED, GREEN, BLUE, GRAY, YELLOW)",
			index, cb.BadgeColor,
		)
	}
	if cb.BadgeLabel == "" {
		return fmt.Errorf("component %d: badgeLabel cannot be empty", index)
	}
	return nil
}

// ValidateComponentRow validates a ComponentRow struct and its contents.
// It takes a ComponentRow pointer and an index indicating its position in a component array.
// The function checks that both RowMainContent and RowAsideContent arrays contain valid
// EventComponents by calling ValidateComponent on each item.
// Returns nil if the row is nil or valid, otherwise returns validation error with row:column index.
func ValidateComponentRow(cr *ComponentRow, index int) error {
	if cr == nil {
		return nil
	}
	if len(cr.RowMainContent) == 0 {
		return nil
	}
	if len(cr.RowAsideContent) == 0 {
		return nil
	}
	for i, c := range cr.RowMainContent {
		err := ValidateComponent(&c, i)
		if err != nil {
			return fmt.Errorf("%d:%d: %w", index, i, err)
		}
	}
	for i, c := range cr.RowAsideContent {
		err := ValidateComponent(&c, i)
		if err != nil {
			return fmt.Errorf("%d:%d: %w", index, i, err)
		}
	}
	return nil
}

// EventComponent Represents the collection of components in the event.
type EventComponent struct {
	ComponentText       *ComponentText       `json:"componentText,omitempty"`
	ComponentSpacer     *ComponentSpacer     `json:"componentSpacer,omitempty"`
	ComponentLinkButton *ComponentLinkButton `json:"componentLinkButton,omitempty"`
	ComponentDivider    *ComponentDivider    `json:"componentDivider,omitempty"`
	ComponentCopyButton *ComponentCopyButton `json:"componentCopyButton,omitempty"`
	ComponentBadge      *ComponentBadge      `json:"componentBadge,omitempty"`
	ComponentRow        *ComponentRow        `json:"componentRow,omitempty"`
}

// ValidateComponent performs validation on a EventComponent struct and its constituent components.
// It takes a EventComponent pointer and an index indicating its position in a component array.
// The function validates each subcomponent (Text, Spacer, LinkButton, Divider,
// CopyButton, Badge) if they are present by calling their respective validation functions.
// Returns nil if the component is valid or nil, otherwise returns the first validation
// error encountered.
// If comp is nil, returns an error indicating an invalid nil component.
func ValidateComponent(comp *EventComponent, index int) error {
	if comp == nil {
		return fmt.Errorf("component %d: invalid nil component", index)
	}

	if comp.ComponentText != nil {
		if err := ValidateComponentText(comp.ComponentText, index); err != nil {
			return err
		}
	}
	if comp.ComponentSpacer != nil {
		if err := ValidateComponentSpacer(comp.ComponentSpacer, index); err != nil {
			return err
		}
	}
	if comp.ComponentLinkButton != nil {
		if err := ValidateComponentLinkButton(comp.ComponentLinkButton, index); err != nil {
			return err
		}
	}
	if comp.ComponentDivider != nil {
		if err := ValidateComponentDivider(comp.ComponentDivider, index); err != nil {
			return err
		}
	}
	if comp.ComponentCopyButton != nil {
		if err := ValidateComponentCopyButton(comp.ComponentCopyButton, index); err != nil {
			return err
		}
	}
	if comp.ComponentBadge != nil {
		if err := ValidateComponentBadge(comp.ComponentBadge, index); err != nil {
			return err
		}
	}
	if comp.ComponentRow != nil {
		if err := ValidateComponentRow(comp.ComponentRow, index); err != nil {
			return err
		}
	}
	return nil
}

// Event Represents the customer or thread event.
type Event struct {
	EventID    string
	Customer   CustomerActor
	Title      string
	Severity   EventSeverity
	Timestamp  time.Time
	Components []EventComponent
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type EventOptions func(e *Event)

func (e *Event) GenId() string {
	return "ev" + xid.New().String()
}

// NewEvent creates a new Event with the given title and optional configurations.
// It generates a unique event ID, sets creation and update timestamps to current UTC time,
// and applies any provided EventOptions. The event is validated before being returned.
//
// Parameters:
//   - title: The title for the event. If empty, will be set to "Untitled" during validation
//   - opts: Optional variadic EventOptions that configure the event properties
//
// Returns:
//   - *Event: The created and validated event
//   - error: Error if validation fails, nil otherwise
func NewEvent(title string, opts ...EventOptions) (*Event, error) {
	eventId := (&Event{}).GenId()
	now := time.Now().UTC()
	event := &Event{
		EventID:   eventId,
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, opt := range opts {
		opt(event)
	}

	if err := ValidateEvent(event); err != nil {
		return nil, err // Return nil instead of empty Event on error
	}
	return event, nil
}

func SetEventCustomer(customer CustomerActor) EventOptions {
	return func(e *Event) {
		e.Customer = customer
	}
}

func SetEventSeverity(severity string) EventOptions {
	return func(e *Event) {
		e.Severity = EventSeverity(severity)
	}
}

func SetEventTimestampFromStr(timestamp string) EventOptions {
	return func(e *Event) {
		e.Timestamp = utils.FromRFC3339OrNow(timestamp)
	}
}

func WithEventComponents(components []EventComponent) EventOptions {
	return func(e *Event) {
		e.Components = components
	}
}

// ValidateEvent validates an Event struct and its components.
// It performs the following validations:
// - Checks if the event is not nil
// - Sets default title to "Untitled" if empty
// - Sets default severity to "muted" if empty, otherwise validates severity value
// - Validates all components in the event
//
// Parameters:
//   - event: Pointer to the Event struct to validate
//
// Returns:
//   - error: nil if validation passes, otherwise returns error describing the validation failure
func ValidateEvent(event *Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	if event.Title == "" {
		event.Title = "Untitled"
	}

	if event.Severity == "" {
		event.Severity = SeverityMuted
	} else {
		if !event.Severity.IsValid() {
			return fmt.Errorf("invalid severity value %s", event.Severity)
		}
	}

	validComponents := make([]EventComponent, 0, len(event.Components))
	for i, comp := range event.Components {
		if err := ValidateComponent(&comp, i); err != nil {
			return err
		}
		validComponents = append(validComponents, comp)
	}
	event.Components = validComponents

	return nil
}
