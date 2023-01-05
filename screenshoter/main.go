package screenshoter

import (
	"context"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/chromedp/chromedp"
)

func Screenshot(url string, filename string, roleNumber int) {
	// Стартуем хром
	ctx, cancel := chromedp.NewContext(context.Background())
	// Не забываем, что его надо закрыть
	// при выходе из main вызовется cancel и он передаст
	// хрому о закрытии
	defer cancel()

	// инициализируем пустой массив, куда будет сохранен скриншот
	var imageBuf []byte
	// и отправляем хрому задачи, которые он должен выполнить
	// у нас только одна - ScreenshotTasks, но можно закинуть сколько угодно
	if err := chromedp.Run(
		ctx,
		ScreenshotTasks(url, &imageBuf, roleNumber),
	); err != nil {
		log.Fatal(err)
	}

	// Задача выполнена, можно сохранить полученное изображение в файл
	if err := ioutil.WriteFile("./pics/"+filename, imageBuf, 0644); err != nil {
		log.Fatal(err)
	}

}

// ScreenshotTasks записывает в imageBuf скриншот страницы, расположенной на url
func ScreenshotTasks(url string, imageBuf *[]byte, roleNumber int) chromedp.Tasks {
	return chromedp.Tasks{
		// задача (таска) состоит из последовательности действий
		// сначала мы переходим по заданному url
		chromedp.Navigate(url),
		chromedp.ScrollIntoView(".roles", chromedp.ByQuery),
		chromedp.Click("#tabs-"+strconv.Itoa(roleNumber), chromedp.ByQuery),
		chromedp.Screenshot(".roles", imageBuf, chromedp.ByQuery),
	}
}
