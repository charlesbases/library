package encoder

import (
	"fmt"
	"testing"
)

func TestMD5(t *testing.T) {
	encoder := NewMD5()
	fmt.Println(encoder.Encode("中国"))
}

func TestHMAC(t *testing.T) {
	encoder := NewHMAC(WithSecretKey([]byte("dahsjdaydasd")))
	data := encoder.Encode("中国")
	fmt.Println(data)
}

var (
	publicKey Bytes = []byte(`
-----BEGIN RSA Public Key-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDe1I0NexZRl69yeVGoqw8SMqnn
sTmryRjFZme1iYhjP6XadXNGosWdCL9sa7CsuEXwVaw5CzrciuoygwTpgw+NhSXI
aPGs+OuXrckCoOVVD2d0YZlBNiSwVDcQJPhhmYJQ6A3HyxlsUC/MF+UEMEBiGH8/
bliK/0dHiE4k/gAR3QIDAQAB
-----END RSA Public Key-----
`)
	privateKey Bytes = []byte(`
-----BEGIN RSA Private Key-----
MIICXQIBAAKBgQDe1I0NexZRl69yeVGoqw8SMqnnsTmryRjFZme1iYhjP6XadXNG
osWdCL9sa7CsuEXwVaw5CzrciuoygwTpgw+NhSXIaPGs+OuXrckCoOVVD2d0YZlB
NiSwVDcQJPhhmYJQ6A3HyxlsUC/MF+UEMEBiGH8/bliK/0dHiE4k/gAR3QIDAQAB
AoGAadEiEq7LEIAp7wCxyJlDFO8+RCqjKnLa8pMI2Oqw/ACzsCRqU7bkhQgsbz5M
rhjsDY+Bs60jjKvjP418fa+haD5gPWY4qf+jGD4HLBZbaVdrRpmW9mCcVgF4up3x
Plvz3xRfKvNvON4gKG1t7NUprZEyn6WfAjSYmp+7Go+P6vkCQQDr2jTISRZW5EvV
v63EFYjOKReIavRkfrgpdk/yvTDSrJyXL/jFLXBnd+/ScnZm6tYZuMx0GGwbXsyS
Tge9D1srAkEA8d2Rb4z2+EH2gFbEaCSCLxOBtXO5pfsAZO940xgQeQjHBORUOZxS
7QDqnsrDlwm/04oVAVXrkLC7pSCvhHgjFwJAaJg0oD4Jci98kiaXYUZLjWIb1ZvZ
Flg8Q+b8PaI5bLSwHTxhDtC/8KL38FRivfGXUYDq6vGJv/mir597PxT4UQJBAKa3
KWaQ7jOVlEpGhL+cWrgEZCYlDNSiPVVV1Bz9u20SZcyzbnL/lBGVziOCdGuJ5tXz
miL/jI6Bo/Zgn1taTCUCQQDWDYVs4BtJ2ENHWNngcwY0LDnluflVQs8BAZSzM1by
Wq2OoFUAFdrTgVRCpImcE0thMRhQmE7K4X5YvgJa/nqD
-----END RSA Private Key-----
`)
)

func TestRSA(t *testing.T) {
	encoder, err := NewRSA(WithPrivateKey(publicKey, privateKey))
	if err != nil {
		panic(err)
	}
	data := encoder.Encode("中国")
	fmt.Println(data)

	fmt.Println(encoder.Decode(data))
}
