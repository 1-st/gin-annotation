package lib

import "log"

type GenContext struct {
	Root       *GroupTree
	MiddlePool map[string]*MiddlewareFunc
	GroupPool  map[string]*GroupTree
}

// InitGroupTree with middlewares and handlers
func (ctx *GenContext) InitGroupTree(middles map[string]*MiddlewareFunc, handles *[]HandlerFunc) {
	ctx.MiddlePool = middles
	ctx.GroupPool = make(map[string]*GroupTree)
	ctx.Root = new(GroupTree)
	ctx.Root.Children = make(map[string]*GroupTree)
	// append group handlers
	for i, _ := range *handles {
		ctx.PutHandler(&(*handles)[i])
	}
	//append group middleware
	for k, _ := range middles {
		ctx.PutMiddleware(middles[k])
	}

}

func (ctx *GenContext) PutMiddleware(m *MiddlewareFunc) {
	for p, w := range m.Group {
		grp := ctx.GroupPool[p]
		if grp == nil {
			log.Fatalf("no such group:(%s)", p)
		}
		if len(grp.Middles) == 0 {
			grp.Middles = append(grp.Middles, MWTuple{M: m, W: w})
		} else {
			for i, v := range grp.Middles {
				if v.W == w {
					log.Fatalf("same weight for two middlewares: [1]:%s [2]:%s", m.PackagePath, v.M.PackagePath)
				}
				if v.W < w {
					//right of i
					if i != len(grp.Middles)-1 {
						grp.Middles = append(grp.Middles[:i], append([]MWTuple{{M: m, W: w}}, grp.Middles[i+1:]...)...)
					} else {
						grp.Middles = append(grp.Middles[:i+1], MWTuple{M: m, W: w})
					}
				} else {
					//left of i
					if i != 0 {
						grp.Middles = append(grp.Middles[:i-1], append([]MWTuple{{M: m, W: w}}, grp.Middles[i:]...)...)
					} else {
						grp.Middles = append([]MWTuple{{M: m, W: w}}, grp.Middles...)
					}
				}
			}
		}
	}
}

func (ctx *GenContext) PutHandler(h *HandlerFunc) {
	//middlewares
	for _, id := range h.Need {
		if m, ok := ctx.MiddlePool[id]; ok {
			h.Middles = append(h.Middles, m)
		} else {
			log.Fatal("middleware not exist: ", m)
		}
	}
	// grow group tree
	var grp = ctx.Root
	var absPath = "" // absolute path
	if len(h.GroupArray) == 0 {
		grp.Handlers = append(grp.Handlers, h)
	} else {
		// find one's group
		for _, step := range h.GroupArray {
			var find = false
			if _, ok := grp.Children[step]; ok {
				find = true
			}
			absPath += step
			if !find {
				node := NewGroupTreeNode(step, absPath, ctx, grp)
				grp.Children[step] = node
				ctx.GroupPool[absPath] = node
			}
			grp = grp.Children[step]
		}
		grp.Handlers = append(grp.Handlers, h)
	}
	h.Group = grp
}
