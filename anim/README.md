# Anim

Usage:

```go
posChannel := anim.NewChannel(anim.Vec3(&elem.Position))
posChannel.SetKey(0, m32.Vec3{0, 0, 0})
posChannel.SetKey(5, m32.Vec3{1, 0, 0}) // 5 segonds after

a := anim.New()
a.AddChannel(posChannel)


a.UpdateDelta(...)
```
