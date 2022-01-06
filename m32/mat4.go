package m32

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

// Mat4 is a 4x4 matrix in row major order.
//
// m[4*r + c] is the element in the r'th row and c'th column.
type Mat4 [16]float32

// Col returns a vector representing the corresponding column (starting at col 0).
// This package makes no distinction between row and column vectors, so it
// will be a normal VecN for a MxN matrix.
func (m Mat4) Col(n int) Vec4 {
	return Vec4{m[n*4+0], m[n*4+1], m[n*4+2], m[n*4+3]}
}

// Abs returns the element-wise absolute value of this matrix
func (m Mat4) Abs() Mat4 {
	return Mat4{
		Abs(m[0]), Abs(m[1]), Abs(m[2]), Abs(m[3]),
		Abs(m[4]), Abs(m[5]), Abs(m[6]), Abs(m[7]),
		Abs(m[8]), Abs(m[9]), Abs(m[10]), Abs(m[11]),
		Abs(m[12]), Abs(m[13]), Abs(m[14]), Abs(m[15]),
	}
}

// Add performs an element-wise addition of two matrices, this is
// equivalent to iterating over every element of m1 and adding the corresponding value of m2.
func (m Mat4) Add(m2 Mat4) Mat4 {
	return Mat4{
		m[0] + m2[0], m[1] + m2[1], m[2] + m2[2], m[3] + m2[3],
		m[4] + m2[4], m[5] + m2[5], m[6] + m2[6], m[7] + m2[7],
		m[8] + m2[8], m[9] + m2[9], m[10] + m2[10], m[11] + m2[11],
		m[12] + m2[12], m[13] + m2[13], m[14] + m2[14], m[15] + m2[15],
	}
}

// MulS performs a scalar multiplcation of the matrix. This is equivalent to iterating
// over every element of the matrix and multiply it by c.
func (m Mat4) MulS(v float32) Mat4 {
	return Mat4{
		m[0] * v, m[1] * v, m[2] * v, m[3] * v,
		m[4] * v, m[5] * v, m[6] * v, m[7] * v,
		m[8] * v, m[9] * v, m[10] * v, m[11] * v,
		m[12] * v, m[13] * v, m[14] * v, m[15] * v,
	}
}

// Mul performs a "matrix product" between this matrix
// and another of the given dimension. For any two matrices of dimensionality
// MxN and NxO, the result will be MxO. For instance, Mat4 multiplied using
// Mul4x2 will result in a Mat4x2.
func (m Mat4) Mul(m2 Mat4) Mat4 {
	return Mat4{
		m[0]*m2[0] + m[4]*m2[1] + m[8]*m2[2] + m[12]*m2[3],
		m[1]*m2[0] + m[5]*m2[1] + m[9]*m2[2] + m[13]*m2[3],
		m[2]*m2[0] + m[6]*m2[1] + m[10]*m2[2] + m[14]*m2[3],
		m[3]*m2[0] + m[7]*m2[1] + m[11]*m2[2] + m[15]*m2[3],
		m[0]*m2[4] + m[4]*m2[5] + m[8]*m2[6] + m[12]*m2[7],
		m[1]*m2[4] + m[5]*m2[5] + m[9]*m2[6] + m[13]*m2[7],
		m[2]*m2[4] + m[6]*m2[5] + m[10]*m2[6] + m[14]*m2[7],
		m[3]*m2[4] + m[7]*m2[5] + m[11]*m2[6] + m[15]*m2[7],
		m[0]*m2[8] + m[4]*m2[9] + m[8]*m2[10] + m[12]*m2[11],
		m[1]*m2[8] + m[5]*m2[9] + m[9]*m2[10] + m[13]*m2[11],
		m[2]*m2[8] + m[6]*m2[9] + m[10]*m2[10] + m[14]*m2[11],
		m[3]*m2[8] + m[7]*m2[9] + m[11]*m2[10] + m[15]*m2[11],
		m[0]*m2[12] + m[4]*m2[13] + m[8]*m2[14] + m[12]*m2[15],
		m[1]*m2[12] + m[5]*m2[13] + m[9]*m2[14] + m[13]*m2[15],
		m[2]*m2[12] + m[6]*m2[13] + m[10]*m2[14] + m[14]*m2[15],
		m[3]*m2[12] + m[7]*m2[13] + m[11]*m2[14] + m[15]*m2[15],
	}
}

// MulV4 returns a Vec4 based on the multiplication of m and m2.
func (m Mat4) MulV4(m2 Vec4) Vec4 {
	return Vec4{
		m[0]*m2[0] + m[4]*m2[1] + m[8]*m2[2] + m[12]*m2[3],
		m[1]*m2[0] + m[5]*m2[1] + m[9]*m2[2] + m[13]*m2[3],
		m[2]*m2[0] + m[6]*m2[1] + m[10]*m2[2] + m[14]*m2[3],
		m[3]*m2[0] + m[7]*m2[1] + m[11]*m2[2] + m[15]*m2[3],
	}
}

// Det returns the determinant of a matrix. It is a measure of a square matrix's
// singularity and invertability, among other things. In this library, the
// determinant is hard coded based on pre-computed cofactor expansion, and uses
// no loops. Of course, the addition and multiplication must still be done.
func (m Mat4) Det() float32 {
	return m[0]*m[5]*m[10]*m[15] -
		m[0]*m[5]*m[11]*m[14] -
		m[0]*m[6]*m[9]*m[15] +
		m[0]*m[6]*m[11]*m[13] +
		m[0]*m[7]*m[9]*m[14] -
		m[0]*m[7]*m[10]*m[13] -
		m[1]*m[4]*m[10]*m[15] +
		m[1]*m[4]*m[11]*m[14] +
		m[1]*m[6]*m[8]*m[15] -
		m[1]*m[6]*m[11]*m[12] -
		m[1]*m[7]*m[8]*m[14] +
		m[1]*m[7]*m[10]*m[12] +
		m[2]*m[4]*m[9]*m[15] -
		m[2]*m[4]*m[11]*m[13] -
		m[2]*m[5]*m[8]*m[15] +
		m[2]*m[5]*m[11]*m[12] +
		m[2]*m[7]*m[8]*m[13] -
		m[2]*m[7]*m[9]*m[12] -
		m[3]*m[4]*m[9]*m[14] +
		m[3]*m[4]*m[10]*m[13] +
		m[3]*m[5]*m[8]*m[14] -
		m[3]*m[5]*m[10]*m[12] -
		m[3]*m[6]*m[8]*m[13] +
		m[3]*m[6]*m[9]*m[12]
}

// Inv computes the inverse of a square matrix. An inverse is a square matrix
// such that when multiplied by the original, yields the identity.
//
// M_inv * M = M * M_inv = I
//
// In this library, the math is precomputed, and uses no loops, though the
// multiplications, additions, determinant calculation, and scaling are still
// done. This can still be (relatively) expensive for a 4x4.
//
// This function checks the determinant to see if the matrix is invertible.
// If the determinant is 0.0, this function returns the zero matrix. However, due to floating point errors, it is
// entirely plausible to get a false positive or negative.
// In the future, an alternate function may be written which takes in a pre-computed determinant.
func (m Mat4) Inv() Mat4 {
	det := m.Det()
	if FloatEqual(det, float32(0.0)) {
		return Mat4{}
	}

	retMat := Mat4{
		-m[7]*m[10]*m[13] + m[6]*m[11]*m[13] + m[7]*m[9]*m[14] - m[5]*m[11]*m[14] - m[6]*m[9]*m[15] + m[5]*m[10]*m[15],
		m[3]*m[10]*m[13] - m[2]*m[11]*m[13] - m[3]*m[9]*m[14] + m[1]*m[11]*m[14] + m[2]*m[9]*m[15] - m[1]*m[10]*m[15],
		-m[3]*m[6]*m[13] + m[2]*m[7]*m[13] + m[3]*m[5]*m[14] - m[1]*m[7]*m[14] - m[2]*m[5]*m[15] + m[1]*m[6]*m[15],
		m[3]*m[6]*m[9] - m[2]*m[7]*m[9] - m[3]*m[5]*m[10] + m[1]*m[7]*m[10] + m[2]*m[5]*m[11] - m[1]*m[6]*m[11],
		m[7]*m[10]*m[12] - m[6]*m[11]*m[12] - m[7]*m[8]*m[14] + m[4]*m[11]*m[14] + m[6]*m[8]*m[15] - m[4]*m[10]*m[15],
		-m[3]*m[10]*m[12] + m[2]*m[11]*m[12] + m[3]*m[8]*m[14] - m[0]*m[11]*m[14] - m[2]*m[8]*m[15] + m[0]*m[10]*m[15],
		m[3]*m[6]*m[12] - m[2]*m[7]*m[12] - m[3]*m[4]*m[14] + m[0]*m[7]*m[14] + m[2]*m[4]*m[15] - m[0]*m[6]*m[15],
		-m[3]*m[6]*m[8] + m[2]*m[7]*m[8] + m[3]*m[4]*m[10] - m[0]*m[7]*m[10] - m[2]*m[4]*m[11] + m[0]*m[6]*m[11],
		-m[7]*m[9]*m[12] + m[5]*m[11]*m[12] + m[7]*m[8]*m[13] - m[4]*m[11]*m[13] - m[5]*m[8]*m[15] + m[4]*m[9]*m[15],
		m[3]*m[9]*m[12] - m[1]*m[11]*m[12] - m[3]*m[8]*m[13] + m[0]*m[11]*m[13] + m[1]*m[8]*m[15] - m[0]*m[9]*m[15],
		-m[3]*m[5]*m[12] + m[1]*m[7]*m[12] + m[3]*m[4]*m[13] - m[0]*m[7]*m[13] - m[1]*m[4]*m[15] + m[0]*m[5]*m[15],
		m[3]*m[5]*m[8] - m[1]*m[7]*m[8] - m[3]*m[4]*m[9] + m[0]*m[7]*m[9] + m[1]*m[4]*m[11] - m[0]*m[5]*m[11],
		m[6]*m[9]*m[12] - m[5]*m[10]*m[12] - m[6]*m[8]*m[13] + m[4]*m[10]*m[13] + m[5]*m[8]*m[14] - m[4]*m[9]*m[14],
		-m[2]*m[9]*m[12] + m[1]*m[10]*m[12] + m[2]*m[8]*m[13] - m[0]*m[10]*m[13] - m[1]*m[8]*m[14] + m[0]*m[9]*m[14],
		m[2]*m[5]*m[12] - m[1]*m[6]*m[12] - m[2]*m[4]*m[13] + m[0]*m[6]*m[13] + m[1]*m[4]*m[14] - m[0]*m[5]*m[14],
		-m[2]*m[5]*m[8] + m[1]*m[6]*m[8] + m[2]*m[4]*m[9] - m[0]*m[6]*m[9] - m[1]*m[4]*m[10] + m[0]*m[5]*m[10],
	}

	return retMat.MulS(1 / det)
}

// Transpose produces the transpose of this matrix. For any MxN matrix
// the transpose is an NxM matrix with the rows swapped with the columns. For instance
// the transpose of the Mat3x2 is a Mat2x3 like so:
//
//    [[a b]]    [[a c e]]
//    [[c d]] =  [[b d f]]
//    [[e f]]
func (m Mat4) Transpose() Mat4 {
	return Mat4{
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15],
	}
}

// Mat3 returns a truncated mat3 removing last col and row.
func (m Mat4) Mat3() Mat3 {
	return Mat3{
		m[0], m[1], m[2],
		m[4], m[5], m[6],
		m[8], m[9], m[10],
	}
}

// Diag is a basic operation on a square matrix that simply
// returns main diagonal (meaning all elements such that row==col).
func (m Mat4) Diag() Vec4 {
	return Vec4{m[0], m[5], m[10], m[15]}
}

// Sub performs an element-wise subtraction of two matrices, this is equivalent
// to iterating over every element of m1 and subtracting the corresponding
// value of m2.
func (m Mat4) Sub(m2 Mat4) Mat4 {
	return Mat4{
		m[0] - m2[0], m[1] - m2[1], m[2] - m2[2], m[3] - m2[3],
		m[4] - m2[4], m[5] - m2[5], m[6] - m2[6], m[7] - m2[7],
		m[8] - m2[8], m[9] - m2[9], m[10] - m2[10], m[11] - m2[11],
		m[12] - m2[12], m[13] - m2[13], m[14] - m2[14], m[15] - m2[15],
	}
}

// Quat returns the quaternion representation of this matrix.
func (m Mat4) Quat() Quat {
	return QLookAt(m.Col(2).Vec3(), m.Col(1).Vec3())
}

// ApproxEqual performs an element-wise approximate equality test between two
// matrices, as if FloatEqual had been used.
func (m Mat4) ApproxEqual(m2 Mat4) bool {
	for i := range m {
		if !FloatEqual(m[i], m2[i]) {
			return false
		}
	}
	return true
}

// ApproxEqualThreshold performs an element-wise approximate equality test
// between two matrices with a given epsilon threshold, as if
// FloatEqualThreshold had been used.
func (m Mat4) ApproxEqualThreshold(m2 Mat4, threshold float32) bool {
	for i := range m {
		if !FloatEqualThreshold(m[i], m2[i], threshold) {
			return false
		}
	}
	return true
}

// ApproxFuncEqual performs an element-wise approximate equality test between
// two matrices with a given equality functions, intended to be used with
// FloatEqualFunc; although and comparison
// function may be used in practice.
func (m Mat4) ApproxFuncEqual(m2 Mat4, eq func(float32, float32) bool) bool {
	for i := range m {
		if !eq(m[i], m2[i]) {
			return false
		}
	}
	return true
}

// Row returns a vector representing the corresponding row (starting at row 0).
// This package makes no distinction between row and column vectors, so it
// will be a normal VecM for a MxN matrix.
func (m Mat4) Row(row int) Vec4 {
	return Vec4{m[row+0], m[row+4], m[row+8], m[row+12]}
}

// Trace is a basic operation on a square matrix that simply
// sums up all elements on the main diagonal (meaning all elements such that row==col).
func (m Mat4) Trace() float32 {
	return m[0] + m[5] + m[10] + m[15]
}

func (m Mat4) String() string {
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 4, 4, 1, ' ', tabwriter.AlignRight)
	for i := 0; i < 4; i++ {
		for _, col := range m.Row(i) {
			fmt.Fprintf(w, "%f\t", col)
		}

		fmt.Fprintln(w, "")
	}
	w.Flush() // nolint: errcheck

	return buf.String()
}

// M4Ident returns a mat4 identity.
func M4Ident() Mat4 {
	return Mat4{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// M4Mul multiply sequence of matrices
func M4Mul(ms ...Mat4) Mat4 {
	r := ms[0]
	for _, m := range ms[1:] {
		r = r.Mul(m)
	}
	return r
}

// Translate3D returns a translation matrix
func Translate3D(x, y, z float32) Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		x, y, z, 1,
	}
}

// Scale3D returns a scale matrix
func Scale3D(sx, sy, sz float32) Mat4 {
	return Mat4{
		sx, 0, 0, 0,
		0, sy, 0, 0,
		0, 0, sz, 0,
		0, 0, 0, 1,
	}
}

// LookAt based on f32
func LookAt(eye, center, up Vec3) Mat4 {
	f := center.Sub(eye).Normalize()
	s := f.Cross(up.Normalize()).Normalize()
	u := s.Cross(f)

	m := Mat4{
		s[0], u[0], -f[0], 0,
		s[1], u[1], -f[1], 0,
		s[2], u[2], -f[2], 0,
		0, 0, 0, 1,
	}

	return m.Mul(Translate3D(-eye[0], -eye[1], -eye[2]))
}
