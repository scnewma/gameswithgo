package main

import (
	"bufio"
	"fmt"
	"os"
)

type storyNode struct {
	text    string
	yesPath *storyNode
	noPath  *storyNode
}

func (node *storyNode) play() {
	fmt.Println(node.text)

	if node.yesPath == nil && node.noPath == nil {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()
		answer := scanner.Text()
		if answer == "yes" {
			node.yesPath.play()
			break
		} else if answer == "no" {
			node.noPath.play()
			break
		} else {
			fmt.Println("Please answer yes or no.")
		}
	}
}

func (node *storyNode) print() {
	node.doPrint(0)
}

func (node *storyNode) doPrint(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print("  ")
	}
	fmt.Println(node.text)

	if node.yesPath != nil {
		node.yesPath.doPrint(depth + 1)
	}

	if node.noPath != nil {
		node.noPath.doPrint(depth + 1)
	}
}

func main() {
	root := storyNode{"You are at the entrance to a dark cave. Do you want to go in the cave?", nil, nil}
	winning := storyNode{"You have won!", nil, nil}
	losing := storyNode{"You have lost!", nil, nil}
	root.yesPath = &losing
	root.noPath = &winning

	root.print()
}
