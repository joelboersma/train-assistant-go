package main

import (
	"fmt"
	"os"
	"strings"
)

type Domino struct {
	val1 int
	val2 int
}

func (d1 *Domino) equals(d2 *Domino) bool {
	return (d1.val1 == d2.val1 && d1.val2 == d2.val2) || (d1.val1 == d2.val2 && d1.val2 == d2.val1)
}

type TrainTreeNode struct {
	val      int
	children []*TrainTreeNode
}

func (n TrainTreeNode) String() string {
	childrenStrings := []string{}
	for _, child := range n.children {
		childrenStrings = append(childrenStrings, child.String())
	}
	return fmt.Sprintf("%v -> [%v]", n.val, strings.Join(childrenStrings, ", "))
}

func (n *TrainTreeNode) populateChildren(dominoes []Domino) {
	for _, domino := range dominoes {
		if n.val == domino.val1 {
			// Make new node with val2
			newNode := TrainTreeNode{domino.val2, []*TrainTreeNode{}}
			n.children = append(n.children, &newNode)
			newDominoes := removeDomino(dominoes, &domino)
			newNode.populateChildren(newDominoes)
		} else if n.val == domino.val2 {
			// Make new node with val1
			newNode := TrainTreeNode{domino.val1, []*TrainTreeNode{}}
			n.children = append(n.children, &newNode)
			newDominoes := removeDomino(dominoes, &domino)
			newNode.populateChildren(newDominoes)
		}
	}
}

func (n *TrainTreeNode) longestBranches() [][]*TrainTreeNode {
	if len(n.children) == 0 {
		return [][]*TrainTreeNode{{n}}
	}

	// Find longest child branches
	longestChildBranchLength := 0
	allChildBranches := [][]*TrainTreeNode{}
	for _, child := range n.children {
		childBranches := child.longestBranches()
		allChildBranches = append(allChildBranches, childBranches...)
		for _, branch := range childBranches {
			if len(branch) > longestChildBranchLength {
				longestChildBranchLength = len(branch)
			}
		}
	}

	longestBranches := [][]*TrainTreeNode{}
	for _, branch := range allChildBranches {
		if len(branch) == longestChildBranchLength {
			newBranch := make([]*TrainTreeNode, len(branch)+1)
			newBranch[0] = n
			for i, node := range branch {
				newBranch[i+1] = node
			}
			longestBranches = append(longestBranches, newBranch)
		}
	}
	return longestBranches
}

type TrainTree struct {
	root TrainTreeNode
}

func (t TrainTree) String() string {
	return t.root.String()
}

func (t *TrainTree) longestTrains() [][]int {
	branches := t.root.longestBranches()
	trains := make([][]int, len(branches))
	for i, branch := range branches {
		vals := make([]int, len(branch))
		for j, node := range branch {
			vals[j] = node.val
		}
		trains[i] = vals
	}
	return trains
}

// Make a new []Domino that is identical to @dominos,
// except with the first instance of @target removed, if any.
func removeDomino(dominos []Domino, target *Domino) []Domino {
	result := []Domino{}
	found := false
	for _, domino := range dominos {
		if !found && domino.equals(target) {
			found = true
			continue
		}
		result = append(result, domino)
	}
	return result
}

func getTrainTrees(availablePlayValues []int, dominoes []Domino) []TrainTree {
	trainTrees := []TrainTree{}
	for _, val := range availablePlayValues {
		root := TrainTreeNode{val, []*TrainTreeNode{}}
		root.populateChildren(dominoes)
		trainTree := TrainTree{root}
		trainTrees = append(trainTrees, trainTree)
	}
	return trainTrees
}

func getLongestTrains(trainTrees []TrainTree) map[int][][]int {
	trainMap := make(map[int][][]int, len(trainTrees))
	for _, trainTree := range trainTrees {
		trainMap[trainTree.root.val] = trainTree.longestTrains()
	}
	return trainMap
}

func main() {
	filename := getFileNameFromArgs()
	availablePlayValues, dominoes, err := readTrainFile(filename)
	if err != nil {
		printError(err)
		os.Exit(1)
	}
	fmt.Println(availablePlayValues, dominoes)
	trainTrees := getTrainTrees(availablePlayValues, dominoes)
	fmt.Println(trainTrees)
	longestTrains := getLongestTrains(trainTrees)
	fmt.Println(longestTrains)
	// TODO: pretty print
}
