package quickbase

// EnvVarPrefix is the environment variable prefix for configuration options.
const EnvVarPrefix = "QUICKBASE"

type ctxKey int

// CtxKey contexts contain context keys.
const (
	CtxKeyAction ctxKey = iota
	CtxKeyRealmHost
)

// Default* constants contain configuration defaults.
const (
	DefaultConfigFile = "$HOME/.config/quickbase/config"
	DefaultTicketFile = "$HOME/.config/quickbase/ticket"
)

// FieldMode* constants contain valid Quick Base field mode settings.
const (
	FieldModeVirtual = "virtual"
	FieldModeLookup  = "lookup"
)

// FieldType* constants contain valid Quick Base field types.
const (
	FieldTypeCheckbox        = "checkbox"
	FieldTypeDate            = "date"
	FieldTypeDuration        = "duration"
	FieldTypeEmailAddress    = "email"
	FieldTypeFileAttachment  = "file"
	FieldTypeListUser        = "multiuserid"
	FieldTypeMultiSelectText = "multitext"
	FieldTypeNumeric         = "float"
	FieldTypeNumericCurrency = "currency"
	FieldTypeNumericPercent  = "percent"
	FieldTypeNumericRating   = "rating"
	FieldTypePhoneNumber     = "phone"
	FieldTypeReportLink      = "dblink"
	FieldTypeText            = "text"
	FieldTypeTimeOfDay       = "timeofday"
	FieldTypeURL             = "url"
	FieldTypeUser            = "userid"
)

// FieldTypes return all valid Quick Base field types.
func FieldTypes() []string {
	return []string{
		FieldTypeCheckbox,
		FieldTypeDate,
		FieldTypeDuration,
		FieldTypeEmailAddress,
		FieldTypeFileAttachment,
		FieldTypeListUser,
		FieldTypeMultiSelectText,
		FieldTypeNumeric,
		FieldTypeNumericCurrency,
		FieldTypeNumericPercent,
		FieldTypeNumericRating,
		FieldTypePhoneNumber,
		FieldTypeReportLink,
		FieldTypeText,
		FieldTypeTimeOfDay,
		FieldTypeURL,
		FieldTypeUser,
	}
}
