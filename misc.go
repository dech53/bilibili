package bilibili

import (
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	regBv   = regexp.MustCompile(`(?i)bv([\dA-Za-z]{10})`)
	regLive = regexp.MustCompile(`^https://live.bilibili.com/(\d+)`)
)

// UnwrapShortUrl 解析短链接，传入一个完整的短链接。
//
// 第一个返回值如果是"bvid"，则第二个返回值是视频的bvid (string)。
// 第一个返回值如果是"live"，则第二个返回值是直播间id (int)。
func (c *Client) UnwrapShortUrl(shortUrl string) (string, any, error) {
	resp, err := c.resty.R().Get(shortUrl)
	if resp == nil {
		return "", nil, errors.WithStack(err)
	}
	if resp.StatusCode() != 302 {
		return "", nil, errors.Errorf("解析短链接失败，status code: %d", resp.StatusCode())
	}
	url := resp.Header().Get("Location")
	{
		ret := regBv.FindString(url)
		if len(ret) > 0 {
			return "bvid", ret, nil
		}
	}
	{
		ret := regLive.FindStringSubmatch(url)
		if len(ret) > 0 {
			rid, err := strconv.Atoi(ret[1])
			if err != nil {
				return "", nil, errors.WithStack(err)
			}
			return "live", rid, nil
		}
	}
	return "", nil, errors.New("无法解析链接：" + url)
}

// Now 获取当前时间戳
func (c *Client) Now() (time.Time, error) {
	const (
		method = resty.MethodGet
		url    = "https://api.bilibili.com/x/report/click/now"
	)
	type resultType struct {
		Now int64 `json:"now"`
	}
	result, err := execute[*resultType](c, method, url, nil)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(result.Now, 0), nil
}

// Av2Bv 将av号转换为bv号，返回格式为"BV1xxxxxxxxx"。
func Av2Bv(aid int) string {
	const (
		xorCode  = 0x1552356C4CDB
		maxAid   = 1 << 51
		alphabet = "FcwAPNKTMug3GV5Lj7EJnHpWsx4tb8haYeviqBz6rkCy12mUSDQX9RdoZf"
	)
	bvid := []byte("BV1000000000")
	tmp := (maxAid | aid) ^ xorCode
	for _, e := range []int{11, 10, 3, 8, 4, 6, 5, 7, 9} {
		bvid[e] = alphabet[tmp%len(alphabet)]
		tmp /= len(alphabet)
	}
	return string(bvid)
}

// Bv2Av 将bv号转换为av号，传入的bv号格式为"BV1xxxxxxxxx"，前面的"BV"不区分大小写。
func Bv2Av(bvid string) int {
	if len(bvid) != 12 {
		panic("bvid 格式错误: " + bvid)
	}
	const (
		xorCode  = 0x1552356C4CDB
		maskCode = 1<<51 - 1
		alphabet = "FcwAPNKTMug3GV5Lj7EJnHpWsx4tb8haYeviqBz6rkCy12mUSDQX9RdoZf"
	)
	tmp := 0
	for _, e := range []int{9, 7, 5, 6, 4, 8, 3, 10, 11} {
		idx := strings.IndexByte(alphabet, bvid[e])
		tmp = tmp*len(alphabet) + idx
	}
	return (tmp & maskCode) ^ xorCode
}
