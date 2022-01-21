# Anim

Usage:

```go
a := anim.New()
// Basic channel func we pass an interpolator anim.Vec3
posChannel := anim.AddChannel(a, anim.Vec3)
// Triggers when UpdateDelta is called with an update.
posChannel.On(func(v gm.Vec3) {
	elem.Position = v
})
posChannel.SetKey(0, gm.Vec3{0, 0, 0})
posChannel.SetKey(5, gm.Vec3{1, 0, 0}) // 5 segonds after

a.Start()

event.HandleFunc(g, func(e gorge.EventUpdate) {
	a.UpdateDelta(e.DeltaTime())
})
```
