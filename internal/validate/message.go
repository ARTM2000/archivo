package validate

import (
	"fmt"
)

func GetValidatorErrorMessage(tag string, field string, param string, info ...string) string {
	switch tag {
	case "required":
		return fmt.Sprintf(MessageRequired, field)
	case "omitempty":
		return ""
	case "len":
		return fmt.Sprintf(MessageLen, field, param)
	case "min":
		return fmt.Sprintf(MessageMin, field, param)
	case "max":
		return fmt.Sprintf(MessageMax, field, param)
	case "eq":
		return fmt.Sprintf(MessageEqual, field, param)
	case "ne":
		return fmt.Sprintf(MessageNotEqual, field, param)
	case "oneof":
		return fmt.Sprintf(MessageOneOf, field, param)
	case "lt":
		return fmt.Sprintf(MessageLessThan, field, param)
	case "lte":
		return fmt.Sprintf(MessageLessThanEqual, field, param)
	case "gt":
		return fmt.Sprintf(MessageGreaterThan, field, param)
	case "gte":
		return fmt.Sprintf(MessageGreaterThanEqual, field, param)
	case "eqfield":
		return fmt.Sprintf(MessageEqualField, field, param)
	case "nefield":
		return fmt.Sprintf(MessageNotEqualField, field, param)
	case "gtfield":
		return fmt.Sprintf(MessageGreaterThanField, field, param)
	case "gtefield":
		return fmt.Sprintf(MessageGreaterThanEqualField, field, param)
	case "ltfield":
		return fmt.Sprintf(MessageLessThanField, field, param)
	case "ltefield":
		return fmt.Sprintf(MessageLessThanEqualField, field, param)
	case "alpha":
		return fmt.Sprintf(MessageAlpha, field)
	case "alphanum":
		return fmt.Sprintf(MessageAlphaNum, field)
	case "alphaunicode":
		return fmt.Sprintf(MessageAlphaUnicode, field)
	case "alphanumunicode":
		return fmt.Sprintf(MessageAlphaNumUnicode, field)
	case "numeric":
		return fmt.Sprintf(MessageNumeric, field)
	case "number":
		return fmt.Sprintf(MessageNumber, field)
	case "hexadecimal":
		return fmt.Sprintf(MessageHexadecimal, field)
	case "email":
		return fmt.Sprintf(MessageEmail, field)
	case "url":
		return fmt.Sprintf(MessageURL, field)
	case "uri":
		return fmt.Sprintf(MessageURI, field)
	case "base64":
		return fmt.Sprintf(MessageBase64, field)
	case "contains":
		return fmt.Sprintf(MessageContains, field, param)
	case "containsany":
		return fmt.Sprintf(MessageContainsAny, field, param)
	case "excludes":
		return fmt.Sprintf(MessageExcludes, field, param)
	case "excludesall":
		return fmt.Sprintf(MessageExcludesAll, field, param)
	case "excludesrune":
		return fmt.Sprintf(MessageExcludesRune, field, param)
	case "uuid":
		return fmt.Sprintf(MessageUUID, field)
	case "uuid3":
		return fmt.Sprintf(MessageUUID3, field)
	case "uuid4":
		return fmt.Sprintf(MessageUUID4, field)
	case "uuid5":
		return fmt.Sprintf(MessageUUID5, field)
	case "datauri":
		return fmt.Sprintf(MessageDataURI, field)
	case "ipv4":
		return fmt.Sprintf(MessageIPv4, field)
	case "ip":
		return fmt.Sprintf(MessageIP, field)
	case "boolean":
		return fmt.Sprintf(MessageBoolean, field)
	case "cron":
		return fmt.Sprintf(MessageCron, field)

	// custom tag
	case "groupinvalid":
		return fmt.Sprintf(MessageGroupInvalid, field, info[0])
	case "password":
		return fmt.Sprintf(MessagePassword, field)

	default:
		fmt.Printf("not defined tag for validation in messages: %s\n", tag)
		return fmt.Sprintf("Validation error for %s", field)
	}
}
