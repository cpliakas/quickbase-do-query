package quickbase

const EnvVarPrefix = "QUICKBASE"

const (
	DefaultConfigFile = "$HOME/.config/quickbase/config"
	DefaultTicketFile = "$HOME/.config/quickbase/ticket"
)

const (
	FieldModeVirtual = "virtual"
	FieldModeLookup  = "lookup"
)

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
