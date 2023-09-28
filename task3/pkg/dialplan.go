package ari

var exten = map[string]string{
	"101": example(),
	"102": "PJSIP/SOFTPHONE_B",
	"103": "PJSIP/testo",
}

func example() string {
	return "PJSIP/SOFTPHONE_A"
}
