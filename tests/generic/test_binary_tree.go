package generic

type NodeValue interface {
	~int | ~string
}

type Node[T NodeValue] struct {
	value     T
	nodeRight *Node[T]
	nodeLeft  *Node[T]
}

func (bt *Node[T]) Add(T) {

}
