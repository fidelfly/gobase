package containerx

type ListNode struct {
	Val  interface{}
	Next *ListNode
}

type LinkedList struct {
	length int
	head   *ListNode
}

func NewLinkedList() *LinkedList {
	return &LinkedList{
		length: 0,
		head: &ListNode{
			Val:  nil,
			Next: nil,
		},
	}
}
