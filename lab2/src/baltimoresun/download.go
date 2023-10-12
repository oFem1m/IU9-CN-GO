package main

import (
	"github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
	"io"
	"net/http"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

type Item struct {
	Time, Title, SourceLink string
}

func search(node *html.Node) []*Item {
	var items []*Item

	if isElem(node, "article") && (getAttr(node, "class") == "container-fluid row flex_row " ||
		getAttr(node, "class") == "container-fluid row flex_row padding-sm-bottom") {
		item := &Item{}

		for a := node.FirstChild; a != nil; a = a.NextSibling {
			if isElem(a, "a") {
				item.SourceLink = "https://www.baltimoresun.com" + getAttr(a, "href")
				for c := a.FirstChild; c != nil; c = c.NextSibling {
					if isDiv(c, "headline-wrapper flex flex_col col-sm-xl-7") {
						for cChild := c.FirstChild; cChild != nil; cChild = cChild.NextSibling {
							if isElem(cChild, "div") &&
								(getAttr(cChild, "class") == "promo-headline font_20_custom font_mobile_custom "+
									"font_normal story-headline-link color_primary decoration_none  ") ||
								getAttr(cChild, "class") == "promo-headline font_18_custom font_mobile_custom "+
									"font_normal story-headline-link color_primary decoration_none  " {
								for tNode := cChild.FirstChild; tNode != nil; tNode = tNode.NextSibling {
									if isElem(tNode, "h2") || isElem(tNode, "h1") {
										item.Title = tNode.FirstChild.Data
									}
								}
							} else if isElem(cChild, "div") && getAttr(cChild, "class") == "isDisabled" {
								for cChild2 := cChild.FirstChild; cChild2 != nil; cChild2 = cChild2.NextSibling {
									if isElem(cChild2, "div") && getAttr(cChild2, "class") == "padding-xs-top" {
										for tNode := cChild2.FirstChild; tNode != nil; tNode = tNode.NextSibling {
											if isElem(tNode, "time") {
												item.Time = getAttr(tNode, "datetime")
											}
										}
									}
								}
							}
						}
					}
				}
			}

		}
		items = append(items, item)
		log.Info("items:", "items", item)
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		childItems := search(c)
		if childItems != nil {
			items = append(items, childItems...)
		}
	}

	return items
}

func downloadNews() []*Item {
	log.Info("sending request to www.baltimoresun.com")
	if response, err := http.Get("https://www.baltimoresun.com/latest/"); err != nil {
		log.Error("request to baltimoresun.com failed", "error", err)
	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(response.Body)
		status := response.StatusCode
		log.Info("got response from www.baltimoresun.com", "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from www.baltimoresun.com", "error", err)
			} else {
				log.Info("HTML from www.baltimoresun.com parsed successfully")
				return search(doc)
			}
		}
	}
	return nil
}
