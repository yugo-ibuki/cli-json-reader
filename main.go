package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

type Node struct {
	Path     string
	Value    interface{}
	Children []*Node
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: jsonbrowser <json-file>")
		os.Exit(1)
	}

	// JSONファイルを読み込む
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// JSONをパースする
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// ルートノードを作成
	root := &Node{
		Path:  "root",
		Value: jsonData,
	}
	buildTree(root)

	// インタラクティブな探索を開始
	exploreNode(root)
}

func buildTree(node *Node) {
	switch v := node.Value.(type) {
	case map[string]interface{}:
		for key, value := range v {
			child := &Node{
				Path:  fmt.Sprintf("%s.%s", node.Path, key),
				Value: value,
			}
			node.Children = append(node.Children, child)
			buildTree(child)
		}
	case []interface{}:
		for i, value := range v {
			child := &Node{
				Path:  fmt.Sprintf("%s[%d]", node.Path, i),
				Value: value,
			}
			node.Children = append(node.Children, child)
			buildTree(child)
		}
	}
}

func exploreNode(node *Node) {
	for {
		var items []string
		if len(node.Children) > 0 {
			for _, child := range node.Children {
				preview := getPreview(child.Value)
				items = append(items, fmt.Sprintf("%s: %s", strings.TrimPrefix(child.Path, node.Path+"."), preview))
			}
		} else {
			items = []string{"(leaf node) " + fmt.Sprint(node.Value)}
		}
		items = append(items, "⬅️ Back")

		prompt := promptui.Select{
			Label: fmt.Sprintf("Current path: %s", node.Path),
			Items: items,
			Size:  15,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if idx == len(items)-1 {
			return
		}

		if len(node.Children) > 0 {
			exploreNode(node.Children[idx])
		}
	}
}

func getPreview(value interface{}) string {
	switch v := value.(type) {
	case map[string]interface{}:
		return fmt.Sprintf("Object (%d keys)", len(v))
	case []interface{}:
		return fmt.Sprintf("Array (%d items)", len(v))
	default:
		preview := fmt.Sprint(v)
		if len(preview) > 50 {
			preview = preview[:47] + "..."
		}
		return preview
	}
}
