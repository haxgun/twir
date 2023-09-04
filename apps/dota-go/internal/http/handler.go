package http

import (
	"io"
	"net/http"

	json "github.com/bytedance/sonic"
	"github.com/satont/twir/apps/dota-go/internal/types"
)

func (c *Http) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			c.logger.Errorln(err)
		}
	}(r.Body)

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		c.logger.Errorln(err)
		return
	}

	packet := types.Packet{}
	if err := json.Unmarshal(bytes, &packet); err != nil {
		c.logger.Errorln(err)
		return
	}

	if packet.Auth == nil {
		c.logger.Warnln("Http: request without auth, skipped")
		_, _ = w.Write([]byte("not ok"))
		return
	}

	err = c.processor.Process(r.Context(), &packet)
	if err != nil {
		c.logger.Errorln(err)
		_, _ = w.Write([]byte("not ok"))
		return
	}

	_, _ = w.Write([]byte("hello"))
}
