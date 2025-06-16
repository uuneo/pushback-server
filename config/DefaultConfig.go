package config

var DefaultConfig = &Config{
	System: System{
		User:               "",
		Password:           "",
		Address:            "0.0.0.0:8080",
		Name:               "pushback",
		Debug:              false,
		Dsn:                "",
		MaxApnsClientCount: 1,
	},
	Apple: Apple{
		ApnsPrivateKey: `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgvjopbchDpzJNojnc
o7ErdZQFZM7Qxho6m61gqZuGVRigCgYIKoZIzj0DAQehRANCAAQ8ReU0fBNg+sA+
ZdDf3w+8FRQxFBKSD/Opt7n3tmtnmnl9Vrtw/nUXX4ldasxA2gErXR4YbEL9Z+uJ
REJP/5bp
-----END PRIVATE KEY-----`,
		Topic:   "me.uuneo.Meoworld",
		KeyID:   "BNY5GUGV38",
		TeamID:  "FUWV6U942Q",
		Develop: true,
		AdminId: "",
	},
}
