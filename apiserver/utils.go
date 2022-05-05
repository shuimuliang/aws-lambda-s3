package apiserver

import (
    "io"
    "encoding/json"
    "github.com/bitly/go-simplejson"
)

func GetSNSCallBackBucket(json_byte []byte) (string, string) {
	// 从SNS回调通知中提取bucket_name和object_key
	res, err := simplejson.NewJson(json_byte)
	if err != nil {
		panic(err)
	}
	message := res.Get("Message")
	message_bytes, err := message.Bytes()
	mdic, err := simplejson.NewJson(message_bytes)
	record := mdic.Get("Records").GetIndex(0)
	s3_bucket_name, _ := record.Get("s3").Get("bucket").Get("name").String()
	s3_object_key, _ := record.Get("s3").Get("object").Get("key").String()

	return s3_bucket_name, s3_object_key
}

func DecodeResponse(r io.Reader, response *Response) (err error) {
    d := json.NewDecoder(r)
    d.UseNumber()
    err = d.Decode(response)
    return err
}
