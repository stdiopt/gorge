# resource package revamp

## Things to care about

- Resource tracking
- Avoid dup asset memory
- Avoid double loading
- Resource release

### Trigger message directly to gpu

```go
	func MySystem(g *gorge.Gorge, r *resource.Context) {
		texData := gorge.TextureData{}
		// Load right here
		r.Load(&texData, "file.jpg")
		tex :=gorge.NewTexture(texData)

		g.Handle(func(gorge.EventUpdate) {
			// Bad will block the renderer and possibly lock stuff in wasm
			r.Load(...,"somefile")
		})
	}
```

- If we trigger the resource to load, we might not want it on GPU yet, so we
  need something to choose if we load from
- Resource loading delegation

Remove reference counter from renderer,
resource should be the one that loads, unloads, cache, tracks

This way we will do the ref counter of the resource in the resource system
it will be responsible for:

- Loading stuff
- Tracking resources
- Unloading

## Current:

- (good) each system tracks it's own hardware resources
- (good) avoid reloading hardware resources
- (bad) depends on bundle loading and unloading
- (bad) It has a complex back and forth system
- (bad) async resource loading which might get trickier on things that we need to
  load right away

### Process:

- [app] get bundle loader from resources package that implement gorge.Bundler
- [app] get resource from bundle (bundle registers it)
- [resource] bundle as a message
- [any...] receives the bundle message
- [any...] sends back load intent with required resources
  (each system signs the intention, so if the resource is already in the
  system it might be avoided)
- [resource] loads the necessary resources
- [resource] send data via message for each individual resource
- [any...] receives data and track the loader to the gpu data
- [...][]byte data is no longer needed
- [app] unload bundle
- [any...] discard any hardware resource related to bundled loaders

## Planned

Try to achieve garbage collected resources

- (good) resource tracks all resources triggering load and release messages
- (good) systems don't need to track resources
- (good) resource can implement a []byte cache system for specific resource
- (good) we don't depend on scenes, for tracking resources, and resources will
  be cross scenes
- (bad) bundle will be harder (not impossible)
- (bad) resource might reload if dynamically loaded

### Process:

- [app] resource provider is declared in function params
- [app] Load resource on init
- [resource] send message with resource data and increases the resource
  reference
- [any...] receive the asset data, and registers data locally
- [app] unreference the resource somehow
- [gc,resource] uppon GC we decrease reference, when reference is lost we
  trigger the message to release the resource
- [any...] unload the hardware resource
