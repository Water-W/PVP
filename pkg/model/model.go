package model

type Link interface {
	Endpoints() (Node, Node)
}

type Node interface {
	Links() []Link
	Network() Network
}

type Network interface {
	Nodes() []Node
}
