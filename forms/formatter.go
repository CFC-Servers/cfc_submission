package forms

import (
	"fmt"
	"strings"
)

type Formatter interface {
	GetFormattedContent(submission Submission) FormattedContent
}

type DefaultFormatter struct {
	Color      int
	FieldOrder []string
	FieldNames map[string]string
}

func (formatter DefaultFormatter) Fields(newFields ...string) DefaultFormatter {
	formatter.FieldOrder = append(formatter.FieldOrder, newFields...)
	return formatter
}

func (formatter DefaultFormatter) SetFieldName(field, name string) DefaultFormatter {
	if formatter.FieldNames == nil {
		formatter.FieldNames = make(map[string]string)
	}
	formatter.FieldNames[field] = name

	return formatter
}

func (formatter DefaultFormatter) GetFormattedContent(submission Submission) FormattedContent {
	realm := submission.Fields.GetString("realm")
	realm = getPrettyRealm(realm)
	// TODO description and title shouldnt be hardcoded in the default formatter
	description := fmt.Sprintf(
		"**%v**\n\n%v",
		submission.Fields.GetString("title"),
		submission.Fields.GetString("description"),
	)

	content := FormattedContent{
		Color:       formatter.Color,
		Image:       submission.Fields.GetString("image"),
		Title:       fmt.Sprintf("[%v] %v", realm, submission.FormName),
		Description: description,
		Fields:      nil,
	}

	for _, fieldName := range formatter.FieldOrder {

		v := submission.Fields.GetString(fieldName)
		prettyName, ok := formatter.FieldNames[fieldName]

		if !ok {
			prettyName = strings.Title(fieldName)
		}

		if v != "" {
			content.Fields = append(content.Fields, FormattedContentField{
				Name:  prettyName,
				Value: v,
			})
		}
	}

	return content
}

type FormattedContent struct {
	Image       string
	Title       string
	Description string
	Color       int
	Fields      []FormattedContentField
}

type FormattedContentField struct {
	Name  string
	Value string
}

func getPrettyRealm(realm string) string {
	switch realm {
	case "cfc3":
		return "Build/Kill"
	case "cfcmc":
		return "Minecraft"
	case "discord":
		return "Discord"
	case "cfcttt":
		return "TTT"
	case "cfcprophunt":
		return "Prop Hunt"
	default:
		return realm
	}
}
