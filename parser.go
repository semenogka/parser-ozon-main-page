package main

import (
	"fmt"
	"log"
	"os"
	"time"

	s "github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var (
	hrefs   []string
	names   []string
	prices  []string
	avgs    []string
	reviews []string
	products []Product
)

type Product struct{
	Name string
	Link string
	Avg string
	Price string
	Reviews string
}

func main() {
	service, err := s.NewChromeDriverService("./chromedriver", 4444)
	if err != nil {
		log.Fatal("не удалось запустить драйвер", err)
	}
	defer service.Stop()

	caps := s.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chrome.Capabilities{
		Args: []string{
			"window-size=1920x1080",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--disable-blink-features=AutomationControlled",
			"--disable-infobars",
		},
	})

	driver, _ := s.NewRemote(caps, "")
	if err != nil {
		log.Fatal("не удалось подключится к драйверу", err)
	}
	defer driver.Quit()

	err = driver.Get("https://www.ozon.ru/")
	if err != nil {
		log.Fatal("не удалось открыть страницу", err)
	}

	for i := 0; i < 100; i++ {
		_, err := driver.ExecuteScript("window.scrollBy(0, 700);", nil)
		if err != nil {
			log.Println("Ошибка при скроллинге:", err)
			return
		}
		time.Sleep(200 * time.Millisecond)
	}

	container, err := driver.FindElement(s.ByCSSSelector, "div.container")
	if err != nil {
		log.Println("Не удалось найти контейнер:", err)
		return
	}

	paginator, err := container.FindElement(s.ByXPATH, "./div[last()]")
	if err != nil {
		log.Println("Не удалось найти пагинатор:", err)
		return
	}

	beforeProducts, err := paginator.FindElement(s.ByXPATH, "./div[1]")
	if err != nil {
		log.Println("Не удалось найти див, который перед продуктами:", err)
		return
	}

	allProducts, err := beforeProducts.FindElement(s.ByXPATH, "./div[last()]")
	if err != nil {
		log.Println("Не удалось найти товары:", err)
		return
	}

	allProducts1, err := allProducts.FindElement(s.ByXPATH, "./div[1]")
	if err != nil {
		log.Println("Не удалось найти товары:", err)
		return
	}

	allProducts2, err := allProducts1.FindElement(s.ByXPATH, "./div[1]")
	if err != nil {
		log.Println("Не удалось найти товары:", err)
		return
	}

	allProducts3, err := allProducts2.FindElements(s.ByXPATH, "./div[1]")
	if err != nil {
		log.Println("Не удалось найти товары:", err)
		return
	}

	for _, elem := range allProducts3 {
		for i := 0; i < 29; i++ {
			beforeTenProducts, err := elem.FindElement(s.ByXPATH, fmt.Sprintf("./div[%d]", i+1))
			if err != nil {
				log.Println("Не удалось найти товары:", err)
				return
			}

			tenProducts, err := beforeTenProducts.FindElement(s.ByXPATH, "./div[1]")
			if err != nil {
				log.Println("Не удалось найти товары:", err)
				return
			}
			for i := 0; i < 10; i++ {
				var spanIndex int = 0
				var productJSON Product

				product, err := tenProducts.FindElement(s.ByXPATH, fmt.Sprintf("./div[%d]", i+1))
				if err != nil {
					log.Println("Не удалось найти товар:", err)
					return
				}

				link, err := product.FindElement(s.ByCSSSelector, "a")
				if err != nil {
					log.Println("Не удалось найти ссылку на товар:", err)
					continue
				}
				
				href, err := link.GetAttribute("href")
				if err != nil {
					log.Println("ссылка не найдена: ", err)
					continue
				}

				hrefs = append(hrefs, href)
				productJSON.Link = href

				spans, err := product.FindElements(s.ByCSSSelector, "span")
				if err != nil {
					log.Println("Не удалось найти спаны товара:", err)
					continue
				}

				if len(spans) == 8 {
					for _, span := range spans {
						if spanIndex == 0 {
							price, _ := span.Text()
							prices = append(prices, price)
							productJSON.Price = price
						}
						if spanIndex == 3 {
							name, _ := span.Text()
							names = append(names, name)
							productJSON.Name = name
						}

						if spanIndex == 4 {
							avg, _ := span.Text()
							avgs = append(avgs, avg)
							productJSON.Avg = avg
						}

						if spanIndex == 7 {
							review, _ := span.Text()
							reviews = append(reviews, review)
							productJSON.Reviews = review
						}

						spanIndex++
					}
				} else {
					for _, span := range spans {
						if spanIndex == 0 {
							price, _ := span.Text()
							prices = append(prices, price)
							productJSON.Price = price
						}
						if spanIndex == 2 {
							name, _ := span.Text()
							names = append(names, name)
							productJSON.Name = name
						}

						if spanIndex == 4 {
							avg, _ := span.Text()
							avgs = append(avgs, avg)
							productJSON.Avg = avg
						}

						if spanIndex == 6 {
							review, _ := span.Text()
							reviews = append(reviews, review)
							productJSON.Reviews = review
						}

						spanIndex++
					}
				}
				products = append(products, productJSON)
			}
		}
	}

	json, _ := os.Create("products.json")
	for _, product := range products {
		line := fmt.Sprintf("{\"name\": \"%s\", \"price\": \"%s\", \"avg\": \"%s\", \"reviews\": \"%s\", \"link\": \"%s\"}\n", product.Name, product.Price, product.Avg, product.Reviews, product.Link)
		_, err := json.WriteString(line)
		if err != nil {
			log.Println("Ошибка при записи в файл:", err)
			return
		}
		log.Println(line)
	}

	log.Println("Данные успешно добавлены в файл games.json")
}
