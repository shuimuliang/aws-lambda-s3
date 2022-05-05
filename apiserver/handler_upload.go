package apiserver

import (
	"ascendex.io/act-aws-lambda-s3/library/ecode"
    "github.com/bitly/go-simplejson"
	"encoding/json"
	"io/ioutil"
	"net/http"
    "go.uber.org/zap"
)

//--------------------- handler method ---------------------

func (a *App) handlerSNSS3Callback(w http.ResponseWriter, r *http.Request) {
	var body []byte
	var err error
	var isConfirm bool

	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		err = ecode.RequestErr
		goto failed
	}

	isConfirm, err = a.ConfirmSubscription(body)

	if !isConfirm {
		// err = a.DispatchS3Callback(body)
	}
	body, _ = json.Marshal(&Response{Code: ecode.OK.Code(), Error: "", Data: ""})
	_, _ = w.Write(body)
	return
failed:
	a.log.Warn("handlerSNSS3Callback\t" + err.Error())
	a.WriteErrorResponse(w, err)
}

//--------------------- internal method ---------------------

func (a *App) ConfirmSubscription(body []byte) (isConfirm bool, err error) {
	a.log.Info("the payload is:", zap.String("uuid", string(body)))

	res, err := simplejson.NewJson(body)
	if err != nil {
		return false, err
	}
	message_type, err := res.Get("Type").String()
	if err != nil {
		return false, err
	}
	if message_type == "SubscriptionConfirmation" {
		subscribeURL, err := res.Get("SubscribeURL").String()
		if err != nil {
			return false, err
		}

		resp, err := http.Get(subscribeURL)
		if err != nil {
			return false, err
		}

		defer resp.Body.Close()

		return true, nil
	}
	return false, nil
}

/*
func (a *App) makeLiveImageCalbackDone(object_key string) (err error) {
	var upstatus LiveRoomCoverUpStatus

	if err = a.DB.Where("s3image_key=?", object_key).First(&upstatus).Error; err != nil {
		return errors.Wrap(err, "db sgtok_s3_live_room_cover_upstatus err")
	}
	if err = a.DB.Model(&LiveRoomCoverUpStatus{}).Where("s3image_key = ?", object_key).Update("delflag", S3UploadFlagReady).Error; err != nil {
		return errors.Wrap(err, "db sgtok_s3_live_room_cover_upstatus error")
	}

	// 更新user表对应media_id
	if err = a.DB.Model(&LiveRoom{}).Where("uuid = ?", upstatus.Uuid).Update("cover_media_id", upstatus.MediaID).Error; err != nil {
		return errors.Wrap(err, "db sgtok_live_room error")
	}

	return nil
}
*/
