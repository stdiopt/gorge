package m32

import (
	"math"
)

// Perspective generates a Perspective Matrix.
//
func Perspective(fovy, aspect, near, far float32) Mat4 {
	fovy = (fovy * math.Pi) / 180.0 // convert from degrees to radians
	// fovy = (fovy * math.Pi) / 180.0 // convert from degrees to radians
	nmf, f := near-far, float32(1./math.Tan(float64(fovy)/2.0))

	return Mat4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (near + far) / nmf, -1,
		0, 0, (2. * far * near) / nmf, 0,
	}
}

// Ortho creates an orthographic matrix
// https://stackoverflow.com/questions/31839119/projection-and-orthographic-matrices-for-opengl
func Ortho(left, right, bottom, top, near, far float32) Mat4 {
	rml, tmb, fmn := (right - left), (top - bottom), (far - near)

	return Mat4{
		2. / rml, 0, 0, 0,
		0, 2. / tmb, 0, 0,
		0, 0, -2. / fmn, 0,
		-(right + left) / rml, -(top + bottom) / tmb, -(far + near) / fmn, 1,
	}
}
