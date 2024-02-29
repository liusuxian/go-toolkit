/*
 * @Author: liusuxian 382185882@qq.com
 * @Date: 2024-02-28 15:51:23
 * @LastEditors: liusuxian 382185882@qq.com
 * @LastEditTime: 2024-02-29 15:12:43
 * @Description:
 *
 * Copyright (c) 2024 by liusuxian email: 382185882@qq.com, All Rights Reserved.
 */
package redbookweb

// PublishNoteRequest 发布笔记请求数据
type PublishNoteRequest struct {
	Common    PublishNoteCommon    `json:"common"`
	ImageInfo PublishNoteImageInfo `json:"image_info"`
	VideoInfo any                  `json:"video_info"`
}

// PublishNoteResponse 发布笔记响应数据
type PublishNoteResponse struct {
	BusinessBindResults []any           `json:"business_bind_results"`
	Result              int             `json:"result"`
	Success             bool            `json:"success"`
	Msg                 string          `json:"msg"`
	Data                PublishNoteData `json:"data"`
	ShareLink           string          `json:"share_link"`
}

// PublishNoteCommon
type PublishNoteCommon struct {
	Type          string                       `json:"type"`
	Title         string                       `json:"title"`
	NoteID        string                       `json:"note_id"`
	Desc          string                       `json:"desc"`
	Source        string                       `json:"source"`
	BusinessBinds string                       `json:"business_binds"`
	Ats           []any                        `json:"ats"`
	HashTag       []PublishNoteCommonHashTag   `json:"hash_tag"`
	PostLoc       PublishNoteCommonPostLoc     `json:"post_loc"`
	PrivacyInfo   PublishNoteCommonPrivacyInfo `json:"privacy_info"`
}

// PublishNoteCommonHashTag
type PublishNoteCommonHashTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
	Type string `json:"type"`
}

// PublishNoteCommonPostLoc
type PublishNoteCommonPostLoc struct {
}

// PublishNoteCommonPrivacyInfo
type PublishNoteCommonPrivacyInfo struct {
	OpType int `json:"op_type"`
	Type   int `json:"type"`
}

// PublishNoteImageInfo
type PublishNoteImageInfo struct {
	Images []PublishNoteImages `json:"images"`
}

// PublishNoteImageInfo
type PublishNoteImages struct {
	FileID        string                       `json:"file_id"`
	Width         int                          `json:"width"`
	Height        int                          `json:"height"`
	Metadata      PublishNoteImageInfoMetadata `json:"metadata"`
	Stickers      PublishNoteImageInfoStickers `json:"stickers"`
	ExtraInfoJSON string                       `json:"extra_info_json"`
}

// PublishNoteImageInfoMetadata
type PublishNoteImageInfoMetadata struct {
	Source int `json:"source"`
}

// PublishNoteImageInfoStickers
type PublishNoteImageInfoStickers struct {
	Version  int   `json:"version"`
	Floating []any `json:"floating"`
}

// PublishNoteData
type PublishNoteData struct {
	ID    string `json:"id"`
	Score int    `json:"score"`
}
