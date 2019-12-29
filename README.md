# Gorge

gorge - go r? game engine

(reduced? rudimentary?)

Just a personal project that I started by exploring some ECS patterns

## Examples

- [gophers](https://stdiopt.github.io/gorge/wasm.html?t=gophers)
- [simple](https://stdiopt.github.io/gorge/wasm.html?t=simple)
- [scene](https://stdiopt.github.io/gorge/wasm.html?t=scene)
- [boxes](https://stdiopt.github.io/gorge/wasm.html?t=boxes)
- [text](https://stdiopt.github.io/gorge/wasm.html?t=text)

## Platforms

It was developed with wasm in mind others were added later

- wasm
- glfw (linux, windows?, osx?)
- mobile (golang.org/x/mobile - WIP)

## Example

```go
func main() {
	opt := platform.Options{}

	platform.Start(opt, func(g *gorge.Gorge) {
		gorgeutils.TrackballCamera(g)

		light := gorgeutils.NewLight()
		light.SetPosition(0, 10, -4)

		cube := primitive.Cube()

		g.Handle(func(gorge.StartEvent) {
			g.AddEntity(light)
			g.AddEntity(cube)
		})
		g.Handle(func(dt gorge.UpdateEvent) {
			cube.Rotate(0, 1*float32(dt), 0)
		})
	})
}
```

## Notes

It contains a slightly modified version of `golang.org/x/mobile` in vendor
folder as a couple of gles3 bindings were missing
