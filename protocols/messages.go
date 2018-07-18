package protocols

import (
	"gopkg.in/dedis/onet.v2"
	"gopkg.in/dedis/kyber.v2"
)

type Struct_NEW_ROUND struct {
	*onet.TreeNode
	NEW_ROUND
}

type NEW_ROUND struct {
	RoundID int
}

type Struct_ALL_ALL_PARAMETERS struct {
	*onet.TreeNode
	ALL_ALL_PARAMETERS
}

type ALL_ALL_PARAMETERS struct {
	NClients int
	NTrustees int
}

type Struct_PUBLIC_KEY struct {
	*onet.TreeNode
	PUBLIC_KEY
}

type PUBLIC_KEY struct {
	Key kyber.Point
}