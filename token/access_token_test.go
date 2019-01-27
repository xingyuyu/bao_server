package token

import (
	"testing"
)

func Test_ParseAccessToekn(t *testing.T) {
	data := "{\"access_token\":\"18_x1GyKrhITuP6Y14pathwg_u_19sOMyb4tz5lp80ZbbClhCl5yvwXM2ONK_ERcmMWsRpK-OQr6bjHL_V_VG9qKPFJelAuYy2qyfb1KV_9XnoREi5-rw2LZ7Gbzel64B10158bzlqDXEvddQ13JNQfACAAEH\",\"expires_in\":7200}"
	parseAccessToken([]byte(data))
	//log.Println("access_token=", GetToken())
}

// func Test_GetAccessToekn(t *testing.T) {
// 	// data := "{\"access_token\":\"18_x1GyKrhITuP6Y14pathwg_u_19sOMyb4tz5lp80ZbbClhCl5yvwXM2ONK_ERcmMWsRpK-OQr6bjHL_V_VG9qKPFJelAuYy2qyfb1KV_9XnoREi5-rw2LZ7Gbzel64B10158bzlqDXEvddQ13JNQfACAAEH\",\"expires_in\":7200}"
// 	// parseAccessToken([]byte(data))
// 	getWeixinAccessToken()
// 	log.Println("access_token=", GetToken())
// }
