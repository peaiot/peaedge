package models

type DeviceMessage struct {
	Mn   string         `json:"mn"`
	Sign string         `json:"sign"`
	Data map[string]any `json:"data"`
}
