package protocols

import "gopkg.in/dedis/onet.v2"

type NEW_ROUND struct {
	*onet.TreeNode
	roundID int32
}
