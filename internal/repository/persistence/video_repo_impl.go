package persistence

import (
	"eao/internal/model"
	"eao/internal/repository"
	"path"
	"strings"
)

type VideoRepositoryImpl struct{}

func NewVideoRepository() repository.VideoRepository {
	return &VideoRepositoryImpl{}
}

func (r *VideoRepositoryImpl) GetVideoList() ([]model.VideoConfig, error) {
	list := []model.VideoConfig{
		{FileName: "24u7qivyunz.mp4", VideoName: "爸爸带着女儿买烧鸡"},
		{FileName: "20260628/char_clip4.mp4", VideoName: "单韵母"},
		{FileName: "20260628/char.mp4", VideoName: "一年级拼音"},
		{FileName: "a9d922da.mp4", VideoName: "太有才了这哥们"},
		{FileName: "a9d922dc.mp4", VideoName: "不要学电影里炒菜，都是假的"},
		{FileName: "last-night.mp4", VideoName: "英雄联盟CG-最后的曙光"},
		{FileName: "a9d922e1.mp4", VideoName: "舔着舔着突然发现不对劲"},
		{FileName: "vjsab8nx29p.mp4", VideoName: "今晚王者荣耀双排"},

		{FileName: "20260613/c4f3cea3de5cb1.mp4", VideoName: "猴子的无语瞬间"},
		{FileName: "20260613/c4f3cea3de5cb2.mp4", VideoName: "猴子吃魔鬼辣椒"},
		{FileName: "20260613/c4f3cea3de5cb3.mp4", VideoName: "猴子欺负游客"},
		{FileName: "20260613/c4f3cea3de5cb4.mp4", VideoName: "来吞"},
		{FileName: "20260613/c4f3cea3de5cb5.mp4", VideoName: "薯片桶整蛊猴子"},
		{FileName: "20260613/c4f3cea3de5cb6.mp4", VideoName: "动物们的无语瞬间"},

		{FileName: "20260530/20260530-a07fff81c4571.mp4", VideoName: "人类跳高精华大赏"},
		{FileName: "20260530/20260530-206cc334489d5.mp4", VideoName: "KTV小遛一首无地自容"},
		{FileName: "20260530/20260530-7d648e5917ce88.mp4", VideoName: "喆学家同事在KTV，让所有人都无地自容"},
		{FileName: "20260530/20260530-879cc2564a4a78.mp4", VideoName: "陶喆现在有些梗我都不懂了再现“哦哦啊诶”经典"},
		{FileName: "20260530/20260530-273fe3e6115568.mp4", VideoName: "陶喆婚礼太高兴怒飙生涯最高音G6！《今天你要嫁给我》"},
		{FileName: "20260530/20260530-2370e4f863a74.mp4", VideoName: "林志炫贴脸开大模仿陶喆名场面《胡彦斌你让我哭》"},
		{FileName: "20260530/20260530-ceafb5c93d229.mp4", VideoName: "来KTV挑战一下无地自容"},
		{FileName: "20260530/20260530-ceafb5c93d230.mp4", VideoName: "来KTV复刻无地自容"},
		{FileName: "20260530/20260530-ceafb5c93d231.mp4", VideoName: "南宁的朋友"},
		{FileName: "20260530/20260530-ceafb5c93d232.mp4", VideoName: "大爷在湖南景区喝水，猴子上前欲抢包"},
		{FileName: "20260530/20260530-ceafb5c93d233.mp4", VideoName: "无地自容防空警报"},
		{FileName: "20260530/20260530-ceafb5c93d234.mp4", VideoName: "无地自容小包包了包"},

		{FileName: "a9d9233b.mp4", VideoName: "351国道荆州松滋段"},
		{FileName: "a9d9233f.mp4", VideoName: "美洲豹vs凯门鳄"},
		{FileName: "a9d92409.mp4", VideoName: "无骨草鱼杀法 草鱼去骨去刺整个过程。"},
		{FileName: "acsgldq4ygzl.mp4", VideoName: "龙虾锁屏电脑"},
		{FileName: "d4eqp6p7vmm.mp4", VideoName: "臭豆腐超过印度美食"},
		{FileName: "i8ecwboeoci.mp4", VideoName: "这猝不及防的推背感"},
		{FileName: "j7kaeb6gmk.mp4", VideoName: "自己种的蒜苗"},
		{FileName: "1628846678209mp4.mp4", VideoName: "印度街头版汉堡这也太独特了"},
		{FileName: "1654937296227mp4.mp4", VideoName: "印度街头炒面，一斤面硬生生炒成了六两！"},
		{FileName: "2.mp4", VideoName: "lol觉醒CG动画，艾欧里亚vs偌克萨斯"},
		{FileName: "2013cg.mp4", VideoName: "2013年《英雄联盟》CG宣传片「扭曲的命运」AI修复补帧画质增强版"},
		{FileName: "2hzcohudabj.mp4", VideoName: "好久没玩了"},
		{FileName: "SuratFamousEggRecipe_TikhariOmelette_PintuBhai.mp4", VideoName: "你永远不知道印度的鸡蛋到哪一步才可以吃了"},
		{FileName: "nubg65mfcb.mp4", VideoName: "这里是公安县藕池口，为荆江南岸四口之一"},
		{FileName: "studio_video_1702118321696mp4.mp4", VideoName: "艾欧尼亚与诺克萨斯最为惨烈的一场战役"},
		{FileName: "v28dddg7yvr.mp4", VideoName: "黄豆应该是最全面的食物了，让我们来看看都可以做那些食物出来"},
		{FileName: "ios/IMG_5287.mp4", VideoName: "ios/IMG_5287.mp4"},
		{FileName: "ios/IMG_5288.mp4", VideoName: "ios/IMG_5288.mp4"},
		{FileName: "ios/IMG_5316.mp4", VideoName: "ios/IMG_5316.mp4"},
		{FileName: "ios/IMG_5324.mp4", VideoName: "ios/IMG_5324.mp4"},
		{FileName: "ios/IMG_5343.mp4", VideoName: "ios/IMG_5343.mp4"},
		{FileName: "ios/IMG_5347.mp4", VideoName: "ios/IMG_5347.mp4"},
		{FileName: "ios/IMG_5349.mp4", VideoName: "ios/IMG_5349.mp4"},
		{FileName: "ios/IMG_5354.mp4", VideoName: "ios/IMG_5354.mp4"},
		{FileName: "ios/IMG_5355.mp4", VideoName: "ios/IMG_5355.mp4"},
		{FileName: "ios/IMG_5359.mp4", VideoName: "ios/IMG_5359.mp4"},
		{FileName: "ios/IMG_5394.mp4", VideoName: "ios/IMG_5394.mp4"},
		{FileName: "ios/IMG_5397.mp4", VideoName: "ios/IMG_5397.mp4"},
		{FileName: "ios/IMG_5398.mp4", VideoName: "ios/IMG_5398.mp4"},
		{FileName: "ios/IMG_5399.mp4", VideoName: "ios/IMG_5399.mp4"},
		{FileName: "ios/IMG_5400.mp4", VideoName: "ios/IMG_5400.mp4"},
	}

	for i := range list {
		list[i].Cover = videoCoverName(list[i].FileName)
	}

	return list, nil
}

func videoCoverName(fileName string) string {
	ext := path.Ext(fileName)
	if ext == "" {
		return fileName + ".png"
	}
	return strings.TrimSuffix(fileName, ext) + ".png"
}
