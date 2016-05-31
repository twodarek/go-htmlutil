/*
Package htmlutil provides various utility functions for working with HTML nodes.
*/
package htmlutil

import (
	"bytes"
	"errors"

	"golang.org/x/net/html"
)

// ErrNodeNotFound is returned by any of the functions for getting HTML nodes.
var ErrNodeNotFound = errors.New("no nodes found")

// GetAllHtmlNodes is a convenience function for GetHtmlNodes() that returns all
// HTML nodes.
func GetAllHtmlNodes(n *html.Node, tag string, attr string, attrValue string) ([]*html.Node, error) {
	return GetHtmlNodes(n, tag, attr, attrValue, -1)
}

// GetFirstHtmlNode is a convenience function for GetHtmlNodes() that returns
// the first matching node.
func GetFirstHtmlNode(n *html.Node, tag string, attr string, attrValue string) (*html.Node, error) {
	htmlNodes, err := GetHtmlNodes(n, tag, attr, attrValue, 1)

	if len(htmlNodes) > 0 {
		return htmlNodes[0], err
	}

	return &html.Node{}, err
}

// GetHtmlNodes returns the HTML nodes found within the provided node given a
// tag, attribute, and attribute value up to the provided count.
//
// The attribute and the attribute value are optional. If they are empty, they
// will not be used as search criteria.
//
// If the count is -1, all nodes will be returned.
func GetHtmlNodes(n *html.Node, tag string, attr string, attrValue string, count int) ([]*html.Node, error) {
	var err error
	var foundNodes []*html.Node

	var f func(*html.Node)
	f = func(n *html.Node) {
		// Find the element with the provided tag
		if n.Type == html.ElementNode && n.Data == tag {
			if attr == "" {
				foundNodes = append(foundNodes, n)
			} else {
				for _, a := range n.Attr {
					if attrValue == "" {
						if a.Key == attr {
							foundNodes = append(foundNodes, n)
							if count != -1 && len(foundNodes) >= count {
								break
							}
						}
					} else {
						if a.Key == attr && a.Val == attrValue {
							foundNodes = append(foundNodes, n)
							if count != -1 && len(foundNodes) >= count {
								break
							}
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			// Stop parsing if we've reached the desired count
			if count == -1 || len(foundNodes) < count {
				f(c)
			}
		}
	}
	f(n)

	if len(foundNodes) == 0 {
		err = ErrNodeNotFound
	}

	return foundNodes, err
}

// HtmlNodeToString converts an HTML node to a string for easier printing.
func HtmlNodeToString(n *html.Node) (string, error) {
	var buf bytes.Buffer

	if err := html.Render(&buf, n); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RemoveAllHtmlNodes is a convenience function for RemoveHtmlNodes() that
// removes all matching HTML nodes.
func RemoveAllHtmlNodes(n *html.Node, tag string, attr string, attrValue string) *html.Node {
	return RemoveHtmlNodes(n, tag, attr, attrValue, -1)
}

// RemoveFirstHtmlNode is a convenience function for RemoveHtmlNodes() that
// removes the first matching node.
func RemoveFirstHtmlNode(n *html.Node, tag string, attr string, attrValue string) *html.Node {
	return RemoveHtmlNodes(n, tag, attr, attrValue, 1)
}

// RemoveHtmlNodes removes the HTML nodes found within the provided node given a
// tag, attribute, and attribute value up to the provided count and returns the
// provided node with the nodes removed.
//
// The attribute and the attribute value are optional. If they are empty, they
// will not be used as search criteria.
//
// If the count is -1, all nodes meeting the criteria will be removed.
func RemoveHtmlNodes(n *html.Node, tag string, attr string, attrValue string, count int) *html.Node {
	nodesToDelete, err := GetHtmlNodes(n, tag, attr, attrValue, count)
	// If the nodes to delete aren't found, there's nothing to do
	if err == ErrNodeNotFound {
		return n
	}

	// Delete nodes in reverse order (so the children get deleted first)
	for i := len(nodesToDelete) - 1; i >= 0; i-- {
		nodesToDelete[i].Parent.RemoveChild(nodesToDelete[i])
	}

	return n
}
