package gorlet

import "github.com/stdiopt/gorge"

func calcMaskOn(l int) *gorge.Stencil {
	sbit := 1 << l
	ref := sbit | (sbit - 1)
	rMask := uint32(sbit - 1)
	wMask := uint32(sbit | (sbit - 1))
	// log.Printf("Lvl: %v id: %d wMask: %d rMask: %d", l, ref, wMask, rMask)
	if sbit == 1 {
		return &gorge.Stencil{
			WriteMask: 0xFF,
			Func:      gorge.StencilFuncAlways, Ref: 1, ReadMask: 0xFF,
			Fail: gorge.StencilOpKeep, ZFail: gorge.StencilOpReplace, ZPass: gorge.StencilOpReplace,
		}
	}
	return &gorge.Stencil{
		WriteMask: wMask,
		Func:      gorge.StencilFuncEqual, Ref: ref, ReadMask: rMask,
		Fail: gorge.StencilOpKeep, ZFail: gorge.StencilOpReplace, ZPass: gorge.StencilOpReplace,
	}
}

func calcMaskOff(l int) *gorge.Stencil {
	sbit := 1 << l
	ref := sbit - 1
	rMask := uint32(sbit - 1)
	wMask := uint32(sbit | (sbit - 1))
	if sbit == 1 {
		return &gorge.Stencil{
			WriteMask: 0xFF,
			Func:      gorge.StencilFuncAlways, Ref: 1, ReadMask: 0xFF,
			Fail: gorge.StencilOpKeep, ZFail: gorge.StencilOpZero, ZPass: gorge.StencilOpZero,
		}
	}
	return &gorge.Stencil{
		WriteMask: wMask,
		Func:      gorge.StencilFuncEqual, Ref: ref, ReadMask: rMask,
		Fail: gorge.StencilOpKeep, ZFail: gorge.StencilOpReplace, ZPass: gorge.StencilOpReplace,
	}
}

/*
func maskLevel(e *Entity) int {
	lvl := -1
	var next *Entity
	for {
		p, ok := e.Parent().(*Entity)
		if !ok {
			return lvl
		}
		if p.masked {
			lvl++
		}
		next = p
	}
}
*/
