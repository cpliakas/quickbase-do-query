package quickbase

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// DoQueryInput models an API_DoQuery request.
type DoQueryInput struct {
	XMLName          xml.Name             `xml:"qdbapi"`
	UserToken        string               `xml:"usertoken,omitempty"`
	Ticket           string               `xml:"ticket,omitempty"`
	AppToken         string               `xml:"apptoken,omitempty"`
	Query            string               `xml:"query,omitempty"`
	QueryID          int                  `xml:"qid,omitempty"`
	QueryName        string               `xml:"qname,omitempty"`
	IncludeRecordIDs BoolToInt            `xml:"includeRids,omitempty"`
	ReturnPercentage BoolToInt            `xml:"returnpercentage,omitempty"`
	UseFIDs          BoolToInt            `xml:"useFids,omitempty"`
	Format           string               `xml:"fmt,omitempty"`
	FieldList        DoQueryInput_Fields  `xml:"clist,omitempty"`
	SortList         DoQueryInput_Fields  `xml:"slist,omitempty"`
	Options          DoQueryInput_Options `xml:"options,omitempty"`
}

// Fields sets the fields that are returned.
func (i *DoQueryInput) Fields(fids ...int) *DoQueryInput {
	i.FieldList = fids
	return i
}

// SortBy sets the fields to be sorted by.
func (i *DoQueryInput) SortBy(fids ...int) *DoQueryInput {
	i.SortList = fids
	return i
}

// SortOrder sets the "sortorder" option.
func (i *DoQueryInput) SortOrder(order ...string) *DoQueryInput {
	i.Options.SortOrderList = order
	return i
}

// Limit sets the "num" option.
func (i *DoQueryInput) Limit(n int) *DoQueryInput {
	i.Options.Limit = n
	return i
}

// Offset sets the "skp" option.
func (i *DoQueryInput) Offset(n int) *DoQueryInput {
	i.Options.Limit = n
	return i
}

// OnlyNew sets the "onlynew" option.
func (i *DoQueryInput) OnlyNew() *DoQueryInput {
	i.Options.OnlyNew = true
	return i
}

// Unsorted sets the "nosort" option.
func (i *DoQueryInput) Unsorted() *DoQueryInput {
	i.Options.Unsorted = true
	return i
}

// DoQueryInput_Fields models field lists in API_DoQuery requests.
type DoQueryInput_Fields []int

// MarshalXML converts a list of fields to a "." delimited string.
func (f DoQueryInput_Fields) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(sliceToString(f), start)
}

// DoQueryInput_Options models the "options" element in API_DoQuery requests.
type DoQueryInput_Options struct {
	SortOrderList []string
	Limit         int
	Offset        int
	OnlyNew       bool
	Unsorted      bool
}

// MarshalXML implements Marshaler.MarshalXML and formats the value of the
// "options" element.
func (o DoQueryInput_Options) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	opts := []string{}

	if o.Offset > 0 {
		opts = append(opts, "skp-"+strconv.Itoa(o.Offset))
	}
	if o.Limit > 0 {
		opts = append(opts, "num-"+strconv.Itoa(o.Limit))
	}
	if len(o.SortOrderList) > 0 {
		opts = append(opts, "sortorder-"+strings.Join(o.SortOrderList, ""))
	}
	if o.OnlyNew == true {
		opts = append(opts, "onlynew")
	}
	if o.Unsorted == true {
		opts = append(opts, "nosort")
	}

	return e.EncodeElement(strings.Join(opts, "."), start)
}

// DoQueryOutput models an API_DoQuery response.
type DoQueryOutput struct {
	XMLName   xml.Name               `xml:"qdbapi"`
	ErrorCode int                    `xml:"errcode"`
	ErrorText string                 `xml:"errtext"`
	Fields    []DoQueryOutput_Field  `xml:"table>fields>field"`
	Records   []DoQueryOutput_Record `xml:"table>records>record"`
}

// DoQueryOutput_Field models the "table>fields" element. in an API_DoQuery
// response.
type DoQueryOutput_Field struct {
	FieldID  int    `xml:"id,attr"`
	Type     string `xml:"field_type,attr"`
	BaseType string `xml:"base_type,attr"`
	Mode     string `xml:"mode,attr"`
	Label    string `xml:"label"`
}

// DoQueryOutput_Record models the "table>records>record" element in an
// API_DoQuery response.
type DoQueryOutput_Record struct {
	RecordID int                          `xml:"rid,attr"`
	UpdateID int                          `xml:"update_id"`
	Fields   []DoQueryOutput_Record_Field `xml:"f"`
}

// DoQueryOutput_Record_Field models the "table>records>record>field" element
// in an API_DoQuery response.
type DoQueryOutput_Record_Field struct {
	FieldID int    `xml:"id,attr"`
	Value   string `xml:",chardata"`
}
