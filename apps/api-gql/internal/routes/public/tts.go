package public

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
	"github.com/twirapp/twir/libs/grpc/websockets"
)

type GenerateTTSFileOpts struct {
	Text         string                  `json:"text"`
	Voice        string                  `json:"voice"`
	VoiceService websockets.VoiceService `json:"voiceService"`
	Rate         int                     `json:"rate"`
	Pitch        int                     `json:"pitch"`
	Volume       int                     `json:"volume"`
}

func (p *Public) HandleChannelTTSGenerateFilePost(c *gin.Context) {
	var data GenerateTTSFileOpts
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if data.VoiceService == websockets.VoiceService_RHVOICE || data.VoiceService == -1 {
		reqUrl, err := url.Parse(fmt.Sprintf("http://%s/say", p.config.TTSServiceUrl))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		query := reqUrl.Query()

		query.Set("voice", data.Voice)
		query.Set("pitch", strconv.Itoa(data.Pitch))
		query.Set("volume", strconv.Itoa(data.Volume))
		query.Set("rate", strconv.Itoa(data.Rate))
		query.Set("text", data.Text)

		reqUrl.RawQuery = query.Encode()

		var b bytes.Buffer
		resp, err := req.SetContext(c.Request.Context()).SetOutput(&b).Get(reqUrl.String())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if !resp.IsSuccessState() {
			c.JSON(500, gin.H{"error": fmt.Sprintf("cannot use say %s", resp.String())})
			return
		}

		c.Data(200, "audio/mp3", b.Bytes())
	}
}
