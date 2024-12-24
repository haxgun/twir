package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func RHVoiceTo(e entity.TTSRHVoice) gqlmodel.RHVoice {
	return gqlmodel.RHVoice{
		Country: e.Country,
		Gender:  e.Gender,
		Lang:    e.Lang,
		Name:    e.Name,
		Code:    e.Code,
	}
}
