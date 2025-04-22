# Binary Search Trees

## Implementation examples

### Go

```go
package datastructures

import (
	"cmp"
	"reflect"
	"slices"
	"testing"
)

// NewBSTree returns new binary search tree.
func NewBSTree[T cmp.Ordered](values []T) *BSTree[T] {
	if len(values) == 0 {
		return nil
	}

	valsCopy := make([]T, len(values))
	copy(valsCopy, values)
	slices.Sort(valsCopy)

	var tree *BSTree[T]
	tree = tree.insertBalanced(valsCopy)

	return tree
}

// BSTree represents binary tree.
type BSTree[T cmp.Ordered] struct {
	Value T
	Left  *BSTree[T]
	Right *BSTree[T]
}

// Insert adds single tree node with given value to binary tree. This method doesn't
// rebalance the tree.
func (b *BSTree[T]) Insert(value T) *BSTree[T] {
	if b == nil {
		t := &BSTree[T]{
			Value: value,
			Left:  nil,
			Right: nil,
		}
		b = t
		return b
	}

	if b.Value == value {
		return b
	}

	if value < b.Value {
		b.Left = b.Left.Insert(value)
		return b
	}
	b.Right = b.Right.Insert(value)
	return b
}

// Search performs search in binary tree.
func (b *BSTree[T]) Search(value T) *BSTree[T] {
	if b == nil {
		return nil
	}

	if b.Value == value {
		return b
	}

	if value < b.Value {
		found := b.Left.Search(value)
		return found
	}
	found := b.Right.Search(value)
	return found
}

// insertBalanced is in charge for insert tree nodes with proper selection of base element
// for node value to make it balanced.
func (b *BSTree[T]) insertBalanced(values []T) *BSTree[T] {
	valNum := len(values) / 2
	baseElement := values[valNum]
	if b == nil {
		t := &BSTree[T]{
			Value: baseElement,
		}
		b = t
	}

	left := values[:valNum]
	if len(left) != 0 {
		b.Left = b.Left.insertBalanced(left)
	}
	right := values[valNum+1:]
	if len(right) != 0 {
		b.Right = b.Right.insertBalanced(right)
	}
	return b
}

func TestNewBSTree(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   *BSTree[int]
	}{
		{
			name:   "should successfully crete empty tree with nil vlues",
			values: nil,
			want:   nil,
		},
		{
			name:   "should successfully crete empty tree with empty values",
			values: []int{},
			want:   nil,
		},
		{
			name:   "should successfully crete tree with single value",
			values: []int{1},
			want: &BSTree[int]{
				Value: 1,
				Left:  nil,
				Right: nil,
			},
		},
		{
			name:   "should successfully crete tree with values",
			values: []int{1, 2, 3, 4, 5, 6},
			want: &BSTree[int]{
				Value: 4,
				Left: &BSTree[int]{
					Value: 2,
					Left: &BSTree[int]{
						Value: 1,
					},
					Right: &BSTree[int]{
						Value: 3,
					},
				},
				Right: &BSTree[int]{
					Value: 6,
					Left: &BSTree[int]{
						Value: 5,
						Left:  nil,
						Right: nil,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBSTree(tt.values)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("NewBalancedBSTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBSTree_Insert(t *testing.T) {
	t.Run("should successfully insert value to the tree", func(t *testing.T) {
		values := []int{1, 3, 5}
		tree := NewBSTree(values)
		want := &BSTree[int]{
			Value: 3,
			Left: &BSTree[int]{
				Value: 1,
				Right: &BSTree[int]{
					Value: 2,
				},
			},
			Right: &BSTree[int]{
				Value: 5,
			},
		}

		got := tree.Insert(2)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Insert() = %v, want %v", got, want)
		}
	})
}

func TestBSTree_Search(t *testing.T) {
	t.Run("should successfully find node by value in the tree", func(t *testing.T) {
		values := []int{1, 2, 3, 4, 5, 6, 7, 8}
		tree := NewBSTree(values)
		want := 8

		got := tree.Search(want)
		if got == nil {
			t.Fatal("empty btree node but should find some")
		}
		if !reflect.DeepEqual(got.Value, want) {
			t.Fatalf("Search() = %v, want %v", got.Value, want)
		}
	})

	t.Run("should return nil node when value is not in the tree", func(t *testing.T) {
		values := []string{"value1", "value2", "value3"}
		tree := NewBSTree(values)

		got := tree.Search("non_existing")
		if got != nil {
			t.Fatalf("Search() = '%v' but should be nil", got.Value)
		}
	})

	t.Run("should return nil on empty tree", func(t *testing.T) {
		tree := BSTree[int]{}

		got := tree.Search(123)
		if got != nil {
			t.Fatalf("Search() = '%v' but should be nil", got.Value)
		}
	})
}
```
