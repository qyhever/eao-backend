package persistence

import (
	"eao/internal/model"
	"eao/internal/repository"
)

type VideoRepositoryImpl struct{}

func NewVideoRepository() repository.VideoRepository {
	return &VideoRepositoryImpl{}
}

func (r *VideoRepositoryImpl) GetVideoList() ([]model.VideoConfig, error) {
	return []model.VideoConfig{
		{FileName: "24u7qivyunz.mp4", VideoName: "爸爸带着女儿买烧鸡"},
		{FileName: "a9d922da.mp4", VideoName: "太有才了这哥们"},
		{FileName: "a9d922dc.mp4", VideoName: "不要学电影里炒菜，都是假的"},
		{FileName: "last-night.mp4", VideoName: "英雄联盟CG-最后的曙光"},
		{FileName: "a9d922e1.mp4", VideoName: "舔着舔着突然发现不对劲"},
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
		{FileName: "vjsab8nx29p.mp4", VideoName: "今晚王者荣耀双排"},
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
	}, nil
}
