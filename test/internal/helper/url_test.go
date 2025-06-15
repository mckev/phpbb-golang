package helper

import (
	"testing"

	"phpbb-golang/internal/helper"
)

func TestUrlWithSID(t *testing.T) {
	{
		actual := helper.UrlWithSID("./forums", "123456")
		expected := "./forums?sid=123456"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./forums?f=1", "123456")
		expected := "./forums?f=1&sid=123456"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./posts?p=101#p101", "123456")
		expected := "./posts?p=101&sid=123456#p101"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./topics?f=1&start=100", "123456")
		expected := "./topics?f=1&sid=123456&start=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./topics?f=1&sid=old&start=100", "123456")
		expected := "./topics?f=1&sid=123456&start=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestUrlWithSID_NoSID(t *testing.T) {
	{
		actual := helper.UrlWithSID("./forums", helper.NO_SID)
		expected := "./forums"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./forums?f=1", helper.NO_SID)
		expected := "./forums?f=1"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./posts?p=101#p101", helper.NO_SID)
		expected := "./posts?p=101#p101"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./topics?f=1&start=100", helper.NO_SID)
		expected := "./topics?f=1&start=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./topics?f=1&sid=old&start=100", helper.NO_SID)
		expected := "./topics?f=1&start=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestUrlWithSID_Xss(t *testing.T) {
	{
		actual := helper.UrlWithSID("./top<script>ics?f=<script>1&start<script>=100", "123456")
		expected := "./top%3Cscript%3Eics?f=%3Cscript%3E1&sid=123456&start%3Cscript%3E=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./top%3Cscript%3Eics?f=%3Cscript%3E1&start%3Cscript%3E=100", "123456")
		expected := "./top%3Cscript%3Eics?f=%3Cscript%3E1&sid=123456&start%3Cscript%3E=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./top%3Cscript%3Eics?sid=old&f=%3Cscript%3E1&start%3Cscript%3E=100", helper.NO_SID)
		expected := "./top%3Cscript%3Eics?f=%3Cscript%3E1&start%3Cscript%3E=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}
