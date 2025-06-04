package helper

import (
	"testing"

	"phpbb-golang/internal/helper"
)

func TestEmbedSessionId(t *testing.T) {
	{
		actual := helper.EmbedSessionId("./forums", "123456")
		expected := "./forums?sid=123456"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.EmbedSessionId("./forums?f=1", "123456")
		expected := "./forums?f=1&sid=123456"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.EmbedSessionId("./topics?f=1&start=100", "123456")
		expected := "./topics?f=1&sid=123456&start=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.EmbedSessionId("./topics?f=1&sid=old&start=100", "123456")
		expected := "./topics?f=1&sid=123456&start=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestEmbedSessionId_Xss(t *testing.T) {
	{
		actual := helper.EmbedSessionId("./top<script>ics?f=<script>1&start<script>=100", "123456")
		expected := "./top%3Cscript%3Eics?f=%3Cscript%3E1&sid=123456&start%3Cscript%3E=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.EmbedSessionId("./top%3Cscript%3Eics?f=%3Cscript%3E1&start%3Cscript%3E=100", "123456")
		expected := "./top%3Cscript%3Eics?f=%3Cscript%3E1&sid=123456&start%3Cscript%3E=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}
