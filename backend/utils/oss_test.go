package utils

import "testing"

func TestGetObjectUrl(t *testing.T) {
	LoadConfig("../config.json")

	InitOssClient()

	// 已存在对象
	firstObjectKey := "positive_comment_word_cloud/word_cloud_1.png"
	secondObjectKey := "negative_comment_word_cloud/word_cloud_1.png"
	thirdObjectKey := "site_pie_picture/all_attractions_summary_pie_chart.png"	

	// 不存在对象
	noexisitsObjectKey := "site_pie_picture/xxx.png"

	firstObjectUrl := GetObjectUrl(firstObjectKey)
	secondObjectUrl := GetObjectUrl(secondObjectKey)
	thirdObjectUrl := GetObjectUrl(thirdObjectKey)

	t.Log("first object key: ", firstObjectUrl)
	t.Log("second object key: ", secondObjectUrl)
	t.Log("third object key: ", thirdObjectUrl)

	noexisitsObjectUrl := GetObjectUrl(noexisitsObjectKey)

	t.Log("no exisit key: ", noexisitsObjectUrl)

}

