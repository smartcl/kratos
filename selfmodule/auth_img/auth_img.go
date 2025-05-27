package auth_img

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

const (
	AuthImgCaptchaTypeAudio   = 1
	AuthImgCaptchaTypeString  = 2
	AuthImgCaptchaTypeMath    = 3
	AuthImgCaptchaTypeChinese = 4
	AuthImgCaptchaTypeDigit   = 5
)

var result = base64Captcha.DefaultMemStore

func CreateCode(captchaType int) (string, string, error) {
	var driver base64Captcha.Driver
	switch captchaType {
	case AuthImgCaptchaTypeAudio:
		driver = autoConfig()
	case AuthImgCaptchaTypeString:
		driver = stringConfig()
	case AuthImgCaptchaTypeMath:
		driver = mathConfig()
	case AuthImgCaptchaTypeChinese:
		driver = chineseConfig()
	case AuthImgCaptchaTypeDigit:
		driver = digitConfig()
	}
	if driver == nil {
		panic("生成验证码的类型没有配置，请在yaml文件中配置完再次重试启动项目")
	}
	// 创建验证码并传入创建的类型的配置，以及存储的对象
	c := base64Captcha.NewCaptcha(driver, result)
	id, b64s, err := c.Generate()
	return id, b64s, err
}

func VerifyCaptcha(id, VerifyValue string) bool {
	return result.Verify(id, VerifyValue, true)
}

func GetCodeAnswer(codeId string) string {
	// result 为步骤1 创建的图片验证码存储对象
	return result.Get(codeId, false)
}

// mathConfig 生成图形化算术验证码配置
func mathConfig() *base64Captcha.DriverMath {
	mathType := &base64Captcha.DriverMath{
		Height:          50,
		Width:           100,
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowHollowLine,
		BgColor: &color.RGBA{
			R: 40,
			G: 30,
			B: 89,
			A: 29,
		},
		Fonts: nil,
	}
	return mathType
}

// digitConfig 生成图形化数字验证码配置
func digitConfig() *base64Captcha.DriverDigit {
	digitType := &base64Captcha.DriverDigit{
		Height:   50,
		Width:    100,
		Length:   5,
		MaxSkew:  0.45,
		DotCount: 80,
	}
	return digitType
}

// stringConfig 生成图形化字符串验证码配置
func stringConfig() *base64Captcha.DriverString {
	stringType := &base64Captcha.DriverString{
		Height:          50,
		Width:           100,
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSlimeLine,
		Length:          5,
		Source:          "123456789qwertyuiopasdfghjklzxcvb",
		BgColor: &color.RGBA{
			R: 40,
			G: 30,
			B: 89,
			A: 29,
		},
		Fonts: nil,
	}
	return stringType
}

// chineseConfig 生成图形化汉字验证码配置
func chineseConfig() *base64Captcha.DriverChinese {
	chineseType := &base64Captcha.DriverChinese{
		Height:          50,
		Width:           100,
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowSlimeLine,
		Length:          2,
		Source: "设想,你在,处理,消费者,的音,频输,出音,频可,能无,论什,么都,没有,任何,输出," +
			"或者,它可,能是,单声道,立体声,或是,环绕立,体声的,不想要,的值",
		BgColor: &color.RGBA{
			R: 40,
			G: 30,
			B: 89,
			A: 29,
		},
		Fonts: nil,
	}
	return chineseType
}

// autoConfig 生成图形化数字音频验证码配置
func autoConfig() *base64Captcha.DriverAudio {
	autoType := &base64Captcha.DriverAudio{
		Length:   4,
		Language: "zh",
	}
	return autoType
}
