package main

import (
	"encoding/base32"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

type DecodeCmd struct {
	Path string `arg:"positional"`
}

func (cmd *DecodeCmd) Run() {
	inFile := os.Stdin
	var err error
	if len(cmd.Path) > 0 && cmd.Path != "-" {
		inFile, err = os.Open(cmd.Path)
	}
	if err != nil {
		log.Fatalln("Error opening file:", err)
	}

	zbarData, err := ioutil.ReadAll(inFile)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}

	parseFromUrl(string(zbarData)[8 : len(zbarData)-1])

}

func parseFromUrl(qrCodeQuery string) {

	qrCodeUrl, err := url.Parse(qrCodeQuery)
	if err != nil {
		log.Println("read value: ", err)
	}
	qrCodeValues := qrCodeUrl.Query()
	if err != nil {
		log.Println("read value: ", qrCodeValues)
		log.Fatalln("Can not parse QR data as opt url:", err)
	}

	otpDataEncoded := qrCodeValues["data"][0]
	in, err := base64.StdEncoding.DecodeString(otpDataEncoded)
	if err != nil {
		log.Println("read value: ", otpDataEncoded)
		log.Println("failed to decode as base64 string:", err)
	}

	payload := &MigrationPayload{}
	if err := proto.Unmarshal(in, payload); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	out := []OtpParametersOut{}
	for _, v := range payload.OtpParameters {
		out = append(out, NewOtpParametersOut(v))
	}
	outBytes, err := yaml.Marshal(out)
	if err != nil {
		log.Fatalln("Failed to marshall into yaml:", err)
	}
	os.Stdout.Write(outBytes)
}

func NewOtpParametersOut(o *OtpParameters) OtpParametersOut {

	return OtpParametersOut{
		Secret:    Base32Bytes(o.Secret),
		Name:      o.Name,
		Issuer:    o.Issuer,
		Algorithm: o.Algorithm,
		Digits:    o.Digits,
		Type:      o.Type,
		Counter:   o.Counter,
	}
}

type Base32Bytes []byte

func (src Base32Bytes) MarshalJSON() ([]byte, error) {
	buf := make([]byte, base32.StdEncoding.EncodedLen(len(src)))
	base32.StdEncoding.Encode(buf, src)
	return buf, nil
}

func (src Base32Bytes) MarshalYAML() (interface{}, error) {
	return base32.StdEncoding.EncodeToString(src), nil
}

type OtpParametersOut struct {
	Secret    Base32Bytes `protobuf:"bytes,1,opt,name=secret,proto3,oneof" json:"secret,omitempty"`
	Name      *string     `protobuf:"bytes,2,opt,name=name,proto3,oneof" json:"name,omitempty"`
	Issuer    *string     `protobuf:"bytes,3,opt,name=issuer,proto3,oneof" json:"issuer,omitempty"`
	Algorithm *Algorithm  `protobuf:"varint,4,opt,name=algorithm,proto3,enum=Algorithm,oneof" json:"algorithm,omitempty"`
	Digits    *DigitCount `protobuf:"varint,5,opt,name=digits,proto3,enum=DigitCount,oneof" json:"digits,omitempty"`
	Type      *OtpType    `protobuf:"varint,6,opt,name=type,proto3,enum=OtpType,oneof" json:"type,omitempty"`
	Counter   *int64      `protobuf:"varint,7,opt,name=counter,proto3,oneof" json:"counter,omitempty"`
}
