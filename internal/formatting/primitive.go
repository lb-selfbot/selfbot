package formatting

func BoolYesNo(b bool) string {
	if b {
		return "Yes"
	}

	return "No"
}

func BoolOnOff(b bool) string {
	if b {
		return "On"
	}

	return "Off"
}
