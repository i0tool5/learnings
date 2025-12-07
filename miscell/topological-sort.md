# Topological sort

```go
import (
	"slices"
	"testing"
)

type Node struct {
	ID       string
	ParentID string
}

// dfs implements depth-first search
func dfs(current Node, children map[string][]Node, visited map[string]bool) []Node {

	if visited[current.ID] {
		return nil
	}
	visited[current.ID] = true

	result := []Node{current}
	for _, childID := range children[current.ID] {
		result = append(result, dfs(childID, children, visited)...)
	}

	return result
}

func TopologicalSort(nodes []Node) []Node {
	children := make(map[string][]Node)
	var roots []Node

	for _, node := range nodes {
		if node.ParentID == "" {
			roots = append(roots, node)
			continue
		}
		if node.ParentID != "" {
			children[node.ParentID] = append(children[node.ParentID], node)
		}
	}

	var result []Node
	visited := make(map[string]bool)
	for _, root := range roots {
		result = append(result, dfs(root, children, visited)...)
	}

	return result
}

func TestTopologicalSort(t *testing.T) {
	tests := []struct {
		name  string
		nodes []Node
		want  []Node
	}{
		{
			name: "do not sort one root node",
			nodes: []Node{
				{
					ID:       "node1",
					ParentID: "",
				},
			},
			want: []Node{
				{
					ID:       "node1",
					ParentID: "",
				},
			},
		},
		{
			name: "do not sort with two root nodes",
			nodes: []Node{
				{
					ID:       "node1",
					ParentID: "",
				},
				{
					ID:       "node2",
					ParentID: "",
				},
			},
			want: []Node{
				{
					ID:       "node1",
					ParentID: "",
				},
				{
					ID:       "node2",
					ParentID: "",
				},
			},
		},
		{
			name: "sort two nodes: one root and one child",
			nodes: []Node{
				{
					ID:       "node2",
					ParentID: "node1",
				},
				{
					ID:       "node1",
					ParentID: "",
				},
			},
			want: []Node{
				{
					ID:       "node1",
					ParentID: "",
				},
				{
					ID:       "node2",
					ParentID: "node1",
				},
			},
		},
		{
			name: "sort one root with many children",
			nodes: []Node{
				{
					ID:       "node4",
					ParentID: "node3",
				},
				{
					ID:       "node2",
					ParentID: "node1",
				},
				{
					ID:       "node3",
					ParentID: "node5",
				},
				{
					ID:       "node1",
					ParentID: "",
				},
				{
					ID:       "node5",
					ParentID: "node2",
				},
			},
			want: []Node{
				{
					ID:       "node1",
					ParentID: "",
				},
				{
					ID:       "node2",
					ParentID: "node1",
				},
				{
					ID:       "node5",
					ParentID: "node2",
				},
				{
					ID:       "node3",
					ParentID: "node5",
				},
				{
					ID:       "node4",
					ParentID: "node3",
				},
			},
		},
		{
			name: "sort two roots with many children",
			nodes: []Node{
				{
					ID:       "node4_r1",
					ParentID: "node3_r1",
				},
				{
					ID:       "node2_r1",
					ParentID: "root1",
				},
				{
					ID:       "node3_r1",
					ParentID: "node5_r1",
				},
				{
					ID:       "root1",
					ParentID: "",
				},
				{
					ID:       "node5_r1",
					ParentID: "node2_r1",
				},
				{
					ID:       "node2_r2",
					ParentID: "node1_r2",
				},
				{
					ID:       "node1_r2",
					ParentID: "root2",
				},
				{
					ID:       "root2",
					ParentID: "",
				},
				{
					ID:       "node3_r2",
					ParentID: "node2_r2",
				},
			},
			want: []Node{
				{
					ID:       "root1",
					ParentID: "",
				},
				{
					ID:       "node2_r1",
					ParentID: "root1",
				},
				{
					ID:       "node5_r1",
					ParentID: "node2_r1",
				},
				{
					ID:       "node3_r1",
					ParentID: "node5_r1",
				},
				{
					ID:       "node4_r1",
					ParentID: "node3_r1",
				},
				{
					ID:       "root2",
					ParentID: "",
				},
				{
					ID:       "node1_r2",
					ParentID: "root2",
				},
				{
					ID:       "node2_r2",
					ParentID: "node1_r2",
				},
				{
					ID:       "node3_r2",
					ParentID: "node2_r2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TopologicalSort(tt.nodes)
			eq := slices.Equal(tt.want, got)
			if !eq {
				t.Fatalf("slices arent equal\nwant: %#v\ngot:  %#v\n", tt.want, got)
			}
		})
	}
}
```
