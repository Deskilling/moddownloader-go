package request

func CheckConnection() error {
	_, err := ModrinthWebRequest(ModrinthEndpoint["default"])
	if err != nil {
		return err
	}
	return nil
}
