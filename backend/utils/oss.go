package utils

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var (
	ossClient *cos.Client
)

func InitOssClient() {
	url, _ := url.Parse(GetOssConfig().Url)

	baseUrl := &cos.BaseURL{BucketURL: url}	
	ossClient = cos.NewClient(baseUrl, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID: GetOssConfig().SecretId,
			SecretKey: GetOssConfig().SecretKey,
		},
	})
}

func GetObjectUrl(key string) string {
	objectUrl := ossClient.Object.GetObjectURL(key)
	return objectUrl.String()
}

func GetPositiveCommentPic(siteIndex int) string {
	positiveCommentKey := fmt.Sprintf("positive_comment_word_cloud/word_cloud_%d.png", siteIndex)
	positiveCommentUrl := GetObjectUrl(positiveCommentKey)
	return positiveCommentUrl
}

func GetNegativeCommentPic(siteIndex int) string {
	negativeCommentKey := fmt.Sprintf("negative_comment_word_cloud/word_cloud_%d.png", siteIndex)
	negativeCommentUrl := GetObjectUrl(negativeCommentKey)
	return negativeCommentUrl
}

func GetSiteTouristTypePic(siteIndex int) string {
	touristTypeKey := fmt.Sprintf("site_pie_picture/attraction_%d_pie_chart.png", siteIndex)
	touristTypeUrl := GetObjectUrl(touristTypeKey)
	return touristTypeUrl
}

func GetTotalTouristTypePic() string {
	totalTouristTypeKey := "site_pie_picture/all_attractions_summary_pie_chart.png"	
	totalTouristTypeUrl := GetObjectUrl(totalTouristTypeKey)
	return totalTouristTypeUrl
}

