package treeops

import "container/list"

func AssignOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.getMatchingNodes(nodeToMap(candidate), pathNode.Rhs)

		if err != nil {
			return nil, err
		}

		// grab the first value
		first := rhs.Front()

		if first != nil {
			candidate.UpdateFrom(first.Value.(*CandidateNode))
		}
	}
	return matchingNodes, nil
}

// does not update content or values
func AssignAttributesOperator(d *dataTreeNavigator, matchingNodes *list.List, pathNode *PathTreeNode) (*list.List, error) {
	lhs, err := d.getMatchingNodes(matchingNodes, pathNode.Lhs)
	if err != nil {
		return nil, err
	}
	for el := lhs.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		rhs, err := d.getMatchingNodes(nodeToMap(candidate), pathNode.Rhs)

		if err != nil {
			return nil, err
		}

		// grab the first value
		first := rhs.Front()

		if first != nil {
			candidate.UpdateAttributesFrom(first.Value.(*CandidateNode))
		}
	}
	return matchingNodes, nil
}
