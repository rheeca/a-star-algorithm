package astar

import (
	"a-star/src/utils"
	"container/heap"
	"math"

	"github.com/lafriks/go-tiled"
)

type Cell struct {
	X          int
	Y          int
	Cost       float64
	IsWalkable bool
}

type Path struct {
	Cells       []*Cell
	CurrentCell int
}

type Node struct {
	Cell   *Cell
	Parent *Node
	f      float64 // total cost
	g      float64 // distance between current node and origin node
	h      float64 // heuristic
}

func (p *Path) GetCurrentCell() *Cell {
	if p.CurrentCell >= len(p.Cells) {
		return nil
	}
	return p.Cells[p.CurrentCell]
}

func (p *Path) Next() {
	p.CurrentCell += 1
}

func AStar(m *GridMap, originCell, destCell *Cell) (path *Path) {
	var open PriorityQueue
	closed := []*Node{}
	openList := []*Node{}

	// create origin node and add to the open queuee
	originNode := &Node{
		Cell: originCell,
		g:    originCell.Cost,
	}
	originNode.Parent = originNode
	heap.Push(&open, originNode)
	openList = append(openList, originNode)

	// repeat until there are no more open nodes to check
	for len(open) > 0 {
		q := heap.Pop(&open).(*Node)
		openList = Remove(openList, q.Cell)

		// if destination has been reached, reconstruct path and return
		if q.Cell.X == destCell.X && q.Cell.Y == destCell.Y {
			path := &Path{}
			for q.Cell != originCell {
				path.Cells = append(path.Cells, q.Cell)
				q = q.Parent
			}
			path.Reverse()
			return path
		}

		var cell *Cell

		// add neighboring nodes to the open queue
		// cell to the left of q
		cell = m.GetGridCell(q.Cell.X-1, q.Cell.Y)
		if cell != nil && cell.IsWalkable {
			n := &Node{
				Cell:   cell,
				Parent: q,
				g:      q.g + cell.Cost,
				h:      Heuristic(cell, destCell),
			}
			n.f = n.g + n.h

			if !Contains(closed, n.Cell) && !Contains(openList, n.Cell) {
				heap.Push(&open, n)
				openList = append(openList, n)
			} else if Contains(openList, n.Cell) && n.g > q.g+n.Cell.Cost {
				// if current cost is better than previous cost, change to current path
				openNode := GetFromList(openList, n.Cell)
				openNode.Parent = n.Parent
				openNode.f = n.f
				openNode.g = n.g
				openNode.h = n.h
			} else if Contains(closed, n.Cell) && n.g > q.g+n.Cell.Cost {
				// if current cost is better than previous cost, revisit node
				closed = Remove(closed, n.Cell)

				// add to open list
				heap.Push(&open, n)
				openList = append(openList, n)
			}
		}

		// cell to the right of q
		cell = m.GetGridCell(q.Cell.X+1, q.Cell.Y)
		if cell != nil && cell.IsWalkable {
			n := &Node{
				Cell:   cell,
				Parent: q,
				g:      q.g + cell.Cost,
				h:      Heuristic(cell, destCell),
			}
			n.f = n.g + n.h

			if !Contains(closed, n.Cell) && !Contains(openList, n.Cell) {
				heap.Push(&open, n)
				openList = append(openList, n)
			} else if Contains(openList, n.Cell) && n.g > q.g+n.Cell.Cost {
				// if current cost is better than previous cost, change to current path
				openNode := GetFromList(openList, n.Cell)
				openNode.Parent = n.Parent
				openNode.f = n.f
				openNode.g = n.g
				openNode.h = n.h
			} else if Contains(closed, n.Cell) && n.g > q.g+n.Cell.Cost {
				// if current cost is better than previous cost, revisit node
				closed = Remove(closed, n.Cell)

				// add to open list
				heap.Push(&open, n)
				openList = append(openList, n)
			}
		}

		// cell below q
		cell = m.GetGridCell(q.Cell.X, q.Cell.Y-1)
		if cell != nil && cell.IsWalkable {
			n := &Node{
				Cell:   cell,
				Parent: q,
				g:      q.g + cell.Cost,
				h:      Heuristic(cell, destCell),
			}
			n.f = n.g + n.h

			if !Contains(closed, n.Cell) && !Contains(openList, n.Cell) {
				heap.Push(&open, n)
				openList = append(openList, n)
			} else if Contains(openList, n.Cell) && n.g > q.g+n.Cell.Cost {
				// if current cost is better than previous cost, change to current path
				openNode := GetFromList(openList, n.Cell)
				openNode.Parent = n.Parent
				openNode.f = n.f
				openNode.g = n.g
				openNode.h = n.h
			} else if Contains(closed, n.Cell) && n.g > q.g+n.Cell.Cost {
				// if current cost is better than previous cost, revisit node
				closed = Remove(closed, n.Cell)

				// add to open list
				heap.Push(&open, n)
				openList = append(openList, n)
			}
		}

		// cell above q
		cell = m.GetGridCell(q.Cell.X, q.Cell.Y+1)
		if cell != nil && cell.IsWalkable {
			n := &Node{
				Cell:   cell,
				Parent: q,
				g:      q.g + cell.Cost,
				h:      Heuristic(cell, destCell),
			}
			n.f = n.g + n.h

			if !Contains(closed, n.Cell) && !Contains(openList, n.Cell) {
				heap.Push(&open, n)
				openList = append(openList, n)
			} else if Contains(openList, n.Cell) && n.g > q.g+n.Cell.Cost {
				// if current cost is better than previous cost, change to current path
				openNode := GetFromList(openList, n.Cell)
				openNode.Parent = n.Parent
				openNode.f = n.f
				openNode.g = n.g
				openNode.h = n.h
			} else if Contains(closed, n.Cell) && n.g > q.g+n.Cell.Cost {
				// if current cost is better than previous cost, revisit node
				closed = Remove(closed, n.Cell)

				// add to open list
				heap.Push(&open, n)
				openList = append(openList, n)
			}
		}

		closed = append(closed, q)
	}

	return nil
}

func AddNeighboringCell(cell, destCell *Cell, q *Node, openList, closed []*Node, open PriorityQueue) {
	n := &Node{
		Cell:   cell,
		Parent: q,
		g:      q.g + cell.Cost,
		h:      Heuristic(cell, destCell),
	}
	n.f = n.g + n.h

	if !Contains(closed, n.Cell) && !Contains(openList, n.Cell) {
		heap.Push(&open, n)
		openList = append(openList, n)
	} else if Contains(openList, n.Cell) && n.g > q.g+n.Cell.Cost {
		// if current cost is better than previous cost, change to current path
		openNode := GetFromList(openList, n.Cell)
		openNode.Parent = n.Parent
		openNode.f = n.f
		openNode.g = n.g
		openNode.h = n.h
	} else if Contains(closed, n.Cell) && n.g > q.g+n.Cell.Cost {
		// if current cost is better than previous cost, revisit node
		closed = Remove(closed, n.Cell)

		// add to open list
		heap.Push(&open, n)
		openList = append(openList, n)
	}
}

type GridMap struct {
	Cells      [][]*Cell
	Width      int
	Height     int
	CellWidth  int
	CellHeight int
}

func NewGridMap(gameMap *tiled.Map) *GridMap {
	gridMap := &GridMap{CellWidth: gameMap.TileWidth, CellHeight: gameMap.TileHeight}

	mapTiles := gameMap.Layers[utils.CollisionLayer].Tiles
	for tileY := 0; tileY < gameMap.Height; tileY++ {
		cellRow := []*Cell{}
		for tileX := 0; tileX < gameMap.Width; tileX++ {
			tile := mapTiles[tileY*gameMap.Width+tileX]
			if tile.IsNil() {
				cellRow = append(cellRow, &Cell{X: tileX, Y: tileY, Cost: 1, IsWalkable: true})
			} else {
				cellRow = append(cellRow, &Cell{X: tileX, Y: tileY, Cost: 1, IsWalkable: false})
			}
		}
		gridMap.Cells = append(gridMap.Cells, cellRow)
	}
	gridMap.Width = len(gridMap.Cells[0])
	gridMap.Height = len(gridMap.Cells)

	return gridMap
}

func (m *GridMap) GetGridCell(x, y int) *Cell {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return nil
	}
	return m.Cells[y][x]
}

func Heuristic(cell, destCell *Cell) float64 {
	// Manhattan distance
	return math.Abs(float64(destCell.X-cell.X)) + math.Abs(float64(destCell.Y-cell.Y))
}

func GetCell(x, y int) (cell *Cell) {
	return &Cell{
		X: x / utils.UnitSize,
		Y: y / utils.UnitSize,
	}
}

func (p *Path) Reverse() {
	cells := []*Cell{}
	for i := len(p.Cells) - 1; i >= 0; i-- {
		cells = append(cells, p.Cells[i])
	}
	p.Cells = cells
}

// Priority Queue
// - implement the priority queue as a min heap
type PriorityQueue []*Node

func (q PriorityQueue) Len() int { return len(q) }

func (p PriorityQueue) Less(i, j int) bool {
	// Pop returns the lowest
	return p[i].f < p[j].f
}

func (p PriorityQueue) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (q *PriorityQueue) Push(x any) {
	*q = append(*q, x.(*Node))
}

func (q *PriorityQueue) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*q = old[0 : n-1]
	return item
}

func Contains(nodeList []*Node, cell *Cell) bool {
	for _, n := range nodeList {
		if n.Cell == cell {
			return true
		}
	}
	return false
}

func Remove(nodeList []*Node, cell *Cell) []*Node {
	var index int
	for i, n := range nodeList {
		if n.Cell == cell {
			index = i
			break
		}
	}
	return append(nodeList[:index], nodeList[index+1:]...)
}

func GetFromList(nodeList []*Node, cell *Cell) *Node {
	for i, n := range nodeList {
		if n.Cell == cell {
			return nodeList[i]
		}
	}
	return nil
}
