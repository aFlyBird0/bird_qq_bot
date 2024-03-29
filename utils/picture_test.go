package utils

import (
	"image"
	"os"
	"testing"
)

const msg = `
最后更新于:2023-02-22 22:14:21

---考研分数段统计来啦！---

计算机学硕(共101个分数)
390 - 399: 1人(累计1)
380 - 389: 1人(累计2)
370 - 379: 4人(累计6)
360 - 369: 15人(累计21)
350 - 359: 9人(累计30)
340 - 349: 17人(累计47)
330 - 339: 15人(累计62)
320 - 329: 10人(累计72)
310 - 319: 9人(累计81)
300 - 309: 11人(累计92)
290 - 299: 6人(累计98)
280 - 289: 2人(累计100)
240 - 249: 1人(累计101)
计算机专硕(共131个分数)
390 - 399: 1人(累计1)
380 - 389: 3人(累计4)
370 - 379: 5人(累计9)
360 - 369: 5人(累计14)
350 - 359: 19人(累计33)
340 - 349: 12人(累计45)
330 - 339: 23人(累计68)
320 - 329: 21人(累计89)
310 - 319: 22人(累计111)
300 - 309: 10人(累计121)
290 - 299: 6人(累计127)
280 - 289: 4人(累计131)
软件工程学硕(共15个分数)
340 - 349: 3人(累计3)
320 - 329: 3人(累计6)
310 - 319: 5人(累计11)
300 - 309: 1人(累计12)
290 - 299: 2人(累计14)
280 - 289: 1人(累计15)
软件工程专硕(共48个分数)
370 - 379: 1人(累计1)
350 - 359: 3人(累计4)
340 - 349: 4人(累计8)
330 - 339: 5人(累计13)
320 - 329: 8人(累计21)
310 - 319: 4人(累计25)
300 - 309: 9人(累计34)
290 - 299: 11人(累计45)
280 - 289: 2人(累计47)
260 - 269: 1人(累计48)
中俄(共15个分数)
350 - 359: 1人(累计1)
340 - 349: 2人(累计3)
330 - 339: 1人(累计4)
320 - 329: 3人(累计7)
310 - 319: 2人(累计9)
300 - 309: 1人(累计10)
280 - 289: 2人(累计12)
270 - 279: 3人(累计15)
中日(共7个分数)
340 - 349: 1人(累计1)
310 - 319: 1人(累计2)
300 - 309: 2人(累计4)
290 - 299: 2人(累计6)
270 - 279: 1人(累计7)


以上结果通过群名片分析而得，存在一定误差，仅供参考。
祝大家复试顺利！

---过密分数段分布来啦！(某段达到10人及以上）---

计算机学硕 过密分数段分布
【360 - 369】(共15人)
369分: 1人  368分: 2人  367分: 4人  366分: 3人  365分: 1人
362分: 2人  360分: 2人
【340 - 349】(共17人)
349分: 3人  347分: 6人  345分: 3人  344分: 2人  343分: 1人
341分: 1人  340分: 1人
【330 - 339】(共15人)
339分: 1人  338分: 1人  337分: 1人  336分: 1人  335分: 3人
334分: 1人  333分: 2人  331分: 3人  330分: 2人
【320 - 329】(共10人)
329分: 2人  328分: 4人  327分: 1人  324分: 1人  322分: 2人
【300 - 309】(共11人)
308分: 1人  307分: 1人  303分: 2人  302分: 1人  301分: 3人
300分: 3人

计算机专硕 过密分数段分布
【350 - 359】(共19人)
358分: 1人  357分: 4人  356分: 1人  355分: 2人  354分: 1人
353分: 2人  352分: 1人  351分: 2人  350分: 5人
【340 - 349】(共12人)
349分: 2人  347分: 2人  346分: 1人  345分: 2人  344分: 1人
343分: 1人  341分: 1人  340分: 2人
【330 - 339】(共23人)
339分: 4人  338分: 1人  337分: 2人  336分: 1人  335分: 3人
334分: 2人  333分: 4人  332分: 2人  331分: 1人  330分: 3人
【320 - 329】(共21人)
329分: 2人  328分: 1人  327分: 2人  326分: 2人  325分: 5人
324分: 1人  323分: 1人  322分: 2人  321分: 2人  320分: 3人
【310 - 319】(共22人)
319分: 2人  318分: 1人  317分: 2人  316分: 4人  315分: 4人
314分: 3人  312分: 2人  311分: 1人  310分: 3人
【300 - 309】(共10人)
306分: 1人  305分: 1人  303分: 1人  301分: 1人  300分: 6人

软件工程专硕 过密分数段分布
【290 - 299】(共11人)
299分: 1人  298分: 1人  296分: 1人  295分: 2人  294分: 1人
293分: 2人  292分: 2人  290分: 1人`

func TestString2PicFile(t *testing.T) {
	filePath := "test.png"
	tailPicPath := "qq_group.png"
	f, err := os.Open(tailPicPath)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	tailPic, _, err := image.Decode(f)
	if err != nil {
		t.Error(err)
	}
	err = String2PicFileWithTailPicture(msg, "DingTalk JinBuTi.ttf", filePath, tailPic)
	if err != nil {
		t.Error(err)
	}
}
