module github.com/stdiopt/gorge

go 1.13

require (
	github.com/go-gl/gl v0.0.0-20190320180904-bf2b1f2f34d7
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20191125211704-12ad95a8df72
	github.com/go-gl/mathgl v0.0.0-20190713194549-592312d8590a
	github.com/gohxs/folder2go v0.0.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8
	golang.org/x/mobile v0.0.0-20191210151939-1a1fef82734d
)

replace golang.org/x/mobile => ./vendor/golang.org/x/mobile
