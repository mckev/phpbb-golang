package helper

import (
	"phpbb-golang/internal/helper"

	"testing"
)

func TestFormatAttributeValue(t *testing.T) {
	{
		actual := helper.FormatAttributeValue("User123")
		expected := "User123"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.FormatAttributeValue("User 123")
		expected := `"User 123"`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.FormatAttributeValue("1234567890")
		expected := "1234567890"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestFormatAttributeValue_XssCanHappen(t *testing.T) {
	{
		// Dangerous example: <input type="text" value="\"><script>alert('XSS')</script>">
		actual := helper.FormatAttributeValue(`"><script>alert('XSS')</script>`)
		expected := `"\"><script>alert('XSS')</script>"`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		// Dangerous example: <div title="\" onmouseover=alert('XSS')>Hover me</div>"
		actual := helper.FormatAttributeValue(`" onmouseover=alert('XSS')>Hover me</div>`)
		expected := `"\" onmouseover=alert('XSS')>Hover me</div>"`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		// Dangerous example: <div title="\"><img src=x onerror=alert('XSS')>"
		actual := helper.FormatAttributeValue(`"><img src=x onerror=alert('XSS')>`)
		expected := `"\"><img src=x onerror=alert('XSS')>"`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}
