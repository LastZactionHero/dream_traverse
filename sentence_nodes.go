package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// Node for sentence
type Node struct {
	word    string
	next    *Node
	related []*Node
}

//var rootNode *Node
var nodeCache [1000]*Node

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	sentences := readBook()
	commonWords := readCommonWordlist()

	for _, sentence := range sentences {
		loadSentence(sentence, commonWords)
	}

	dream(nodeCache, commonWords)
}

func loadSentence(sentence string, commonWords []string) {
	words := strings.Split(sentence, " ")

	matchPunct := regexp.MustCompile("[.;,!?'\"]")
	for index, word := range words {
		words[index] = matchPunct.ReplaceAllString(word, "")
	}

	rootNode := Node{word: "BEGIN"}
	buildNodes(words, &rootNode)
	attachToRelated(&rootNode, nodeCache)
	cacheNodes(&nodeCache, &rootNode, commonWords)
}

func dream(cache [1000]*Node, commonWords []string) {
	for i := 0; i < 1000; i++ {
		node := cache[rand.Intn(1000)]
		for node != nil {
			dreamText, relatedNode := dreamTraverse(node)
			fmt.Print(dreamText)
			loadSentence(dreamText, commonWords)
			node = relatedNode
		}
	}
	fmt.Println("")
}

func dreamTraverse(rootNode *Node) (string, *Node) {
	var words []string
	node := rootNode
	var related *Node

	for node != nil {
		relatedCount := len(node.related)
		if relatedCount > 0 {
			node = node.related[rand.Intn(relatedCount)]
		} else {
			node = node.next
		}
		if node != nil {
			words = append(words, node.word)
		}
	}
	words = append(words, ".")
	return strings.Join(words, " "), related
}

func printCache(cache [1000]*Node) {
	for _, n := range cache {
		if n != nil {
			fmt.Println(n.word)
		} else {
			fmt.Println("nil")
		}
	}
}

func attachToRelated(rootNode *Node, cache [1000]*Node) {
	node := rootNode
	for node != nil {
		for _, cacheNode := range cache {
			if cacheNode == nil {
				continue
			}
			if cacheNode.word == node.word {
				cacheNode.related = append(cacheNode.related, node)
			}
		}
		node = node.next
	}
}

func readBook() []string {
	dat, _ := ioutil.ReadFile("./sawyer.txt")
	content := string(dat)
	content = strings.Replace(content, "\n", " ", -1)
	matchPunct := regexp.MustCompile("[.!?'\"]")
	return matchPunct.Split(content, -1)
}

func readCommonWordlist() []string {
	dat, _ := ioutil.ReadFile("./top_1000.txt")
	content := string(dat)
	return strings.Split(content, "\n")
}

func buildNodes(words []string, rootNode *Node) {
	lastNode := rootNode

	for _, word := range words {
		newNode := Node{word: word}
		if rootNode == nil {
			rootNode = &newNode
		} else {
			lastNode.next = &newNode
		}

		lastNode = &newNode
	}
}

func printNodes(rootNode *Node) {
	node := rootNode
	for node != nil {
		fmt.Println(node.word)
		node = node.next
	}
}

func nodesToString(rootNode *Node) string {
	var words []string
	node := rootNode
	for node != nil {
		words = append(words, node.word)
		node = node.next
	}
	return strings.Join(words, " ")
}

func cacheNodes(nodeCache *[1000]*Node, rootNode *Node, commonWords []string) {
	var importantNodes []*Node
	node := rootNode
	for node != nil {
		if !isInWordList(node.word, commonWords) {
			importantNodes = append(importantNodes, node)
		}
		node = node.next
	}

	importantLength := len(importantNodes)
	if importantLength == 0 {
		return
	}
	howManyToUse := rand.Intn(importantLength)
	for i := 0; i < howManyToUse; i++ {
		positionToReplace := rand.Intn(1000)
		nodeToReplaceWith := rand.Intn(importantLength)
		if !isCached(importantNodes[nodeToReplaceWith], nodeCache) {
			nodeCache[positionToReplace] = importantNodes[nodeToReplaceWith]
		}
	}
}

func isCached(node *Node, nodeCache *[1000]*Node) bool {
	for _, cache := range nodeCache {
		if cache != nil && cache.word == node.word {
			return true
		}
	}
	return false
}

func isInWordList(word string, wordList []string) bool {
	for _, wordListItem := range wordList {
		if word == wordListItem {
			return true
		}
	}
	return false
}
