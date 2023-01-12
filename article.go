package bilibili

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
)

type Article struct {
	Id          int      `json:"id"`           // 专栏cvid
	Title       string   `json:"title"`        // 文章标题
	State       int      `json:"state"`        // 固定值0，作用尚不明确
	PublishTime int      `json:"publish_time"` // 发布时间戳（秒）
	Words       int      `json:"words"`        // 文章字数
	ImageUrls   []string `json:"image_urls"`   // 文章封面
	Category    struct { // 文章标签
		Id       int    `json:"id"`
		ParentId int    `json:"parent_id"`
		Name     string `json:"name"`
	} `json:"category"`
	Categories []struct { // 文章标签列表
		Id       int    `json:"id"`
		ParentId int    `json:"parent_id"`
		Name     string `json:"name"`
	} `json:"categories"`
	Summary string   `json:"summary"` // 文章摘要
	Stats   struct { // 文章状态数信息
		View     int `json:"view"`     // 阅读数
		Favorite int `json:"favorite"` // 收藏数
		Like     int `json:"like"`     // 点赞数
		Dislike  int `json:"dislike"`  // 点踩数
		Reply    int `json:"reply"`    // 评论数
		Share    int `json:"share"`    // 分享数
		Coin     int `json:"coin"`     // 投币数
		Dynamic  int `json:"dynamic"`  // 动态转发数
	} `json:"stats"`
	LikeState int `json:"like_state"` // 是否点赞
}

type ArticlesInfo struct {
	List struct { // 文集概览
		Id            int    `json:"id"`             // 文集rlid
		Mid           int    `json:"mid"`            // 文集作者mid
		Name          string `json:"name"`           // 文集名称
		ImageUrl      string `json:"image_url"`      // 文集封面图片url
		UpdateTime    int    `json:"update_time"`    // 文集更新时间戳
		Ctime         int    `json:"ctime"`          // 文集创建时间戳
		PublishTime   int    `json:"publish_time"`   // 文集发布时间戳
		Summary       string `json:"summary"`        // 文集简介
		Words         int    `json:"words"`          // 文集字数
		Read          int    `json:"read"`           // 文集阅读量
		ArticlesCount int    `json:"articles_count"` // 1或3，作用尚不明确
		State         int    `json:"state"`          // 空，作用尚不明确
		Reason        string `json:"reason"`         // 空，作用尚不明确
		ApplyTime     string `json:"apply_time"`     // 空，作用尚不明确
		CheckTime     string `json:"check_time"`     // 空，作用尚不明确
	} `json:"list"`
	Articles []Article `json:"articles"` // 文集内的文章列表
	Author   struct {  // 文集作者信息
		Mid            int            `json:"mid"`  // 作者mid
		Name           string         `json:"name"` // 作者昵称
		Face           string         `json:"face"` // 作者头像url
		Pendant        Pendant        `json:"pendant"`
		OfficialVerify OfficialVerify `json:"official_verify"` // 作者认证信息
		Nameplate      NamePlate      `json:"nameplate"`
		Vip            Vip            `json:"vip"`
	} `json:"author"`
	Last      Article `json:"last"`      // 作用尚不明确
	Attention bool    `json:"attention"` // 是否关注文集作者
}

// GetArticlesInfo 获取文集基本信息
func GetArticlesInfo(id int) (*ArticlesInfo, error) {
	return std.GetArticlesInfo(id)
}
func (c *Client) GetArticlesInfo(id int) (*ArticlesInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("id", strconv.Itoa(id)).Get("https://api.bilibili.com/x/article/list/web/articles")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取文集基本信息")
	if err != nil {
		return nil, err
	}
	var ret *ArticlesInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}

type ArticleViewInfo struct {
	Like      int      `json:"like"`      // 是否点赞，0：未点赞，1：已点赞
	Attention bool     `json:"attention"` // 是否关注文章作者
	Favorite  bool     `json:"favorite"`  // 是否收藏
	Coin      int      `json:"coin"`      // 为文章投币数
	Stats     struct { // 状态数信息
		View     int `json:"view"`     // 阅读数
		Favorite int `json:"favorite"` // 收藏数
		Like     int `json:"like"`     // 点赞数
		Dislike  int `json:"dislike"`  // 点踩数
		Reply    int `json:"reply"`    // 评论数
		Share    int `json:"share"`    // 分享数
		Coin     int `json:"coin"`     // 投币数
		Dynamic  int `json:"dynamic"`  // 动态转发数
	} `json:"stats"`
	Title           string     `json:"title"`             // 文章标题
	BannerUrl       string     `json:"banner_url"`        // 文章头图url
	Mid             int        `json:"mid"`               // 文章作者mid
	AuthorName      string     `json:"author_name"`       // 文章作者昵称
	IsAuthor        bool       `json:"is_author"`         // 固定值true，作用尚不明确
	ImageUrls       []string   `json:"image_urls"`        // 动态封面图片url
	OriginImageUrls []string   `json:"origin_image_urls"` // 文章封面图片url
	Shareable       bool       `json:"shareable"`         // 固定值true，作用尚不明确
	ShowLaterWatch  bool       `json:"show_later_watch"`  // 固定值true，作用尚不明确
	ShowSmallWindow bool       `json:"show_small_window"` // 固定值true，作用尚不明确
	InList          bool       `json:"in_list"`           // 是否收于文集
	Pre             int        `json:"pre"`               // 上一篇文章cvid
	Next            int        `json:"next"`              // 下一篇文章cvid
	ShareChannels   []struct { // 分享方式列表
		Name         string `json:"name"`          // 分享名称：QQ，QQ空间，微信，朋友圈，微博
		Picture      string `json:"picture"`       // 分享图片url
		ShareChannel string `json:"share_channel"` // 分享代号：QQ，QZONE，WEIXIN，WEIXIN_MOMENT，SINA
	} `json:"share_channels"`
}

// GetArticleViewInfo 获取专栏文章基本信息
func GetArticleViewInfo(id int) (*ArticleViewInfo, error) {
	return std.GetArticleViewInfo(id)
}
func (c *Client) GetArticleViewInfo(id int) (*ArticleViewInfo, error) {
	resp, err := c.resty().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("id", strconv.Itoa(id)).Get("https://api.bilibili.com/x/article/viewinfo")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := getRespData(resp, "获取专栏文章基本信息")
	if err != nil {
		return nil, err
	}
	var ret *ArticleViewInfo
	err = json.Unmarshal(data, &ret)
	return ret, errors.WithStack(err)
}
