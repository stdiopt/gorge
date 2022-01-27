package gorlet

import (
	"github.com/stdiopt/gorge"
	"github.com/stdiopt/gorge/math/gm"
)

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

// defaultObservers experiment it will be attached on Create so every gorlet will have these.
func defaultObservers(e *Entity) {
	Observe(e, "anchor", func(v gm.Vec4) {
		e.SetAnchor(v[:]...)
	})
	Observe(e, "rect", func(v gm.Vec4) {
		e.SetRect(v[:]...)
	})
	Observe(e, "margin", func(v gm.Vec4) {
		e.SetMargin(v[:]...)
	})
	Observe(e, "width", e.SetWidth)
	Observe(e, "height", e.SetHeight)
	Observe(e, "pivot", func(v gm.Vec2) {
		e.SetPivot(v[:]...)
	})
	Observe(e, "layout", e.SetLayout)
	Observe(e, "border", func(v gm.Vec4) {
		e.SetBorder(v[:]...)
	})
}
