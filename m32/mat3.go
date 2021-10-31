package m32

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

// Mat3 is a 3x3 matrix in row major order.
//
// m[3*r + c] is the element in the r'th row and c'th column.
type Mat3 [9]float32

// MulV3 return a the Vec3 multiplied by matrix.
func (m Mat3) MulV3(v Vec3) Vec3 {
	return Vec3{
		m[0]*v[0] + m[3]*v[1] + m[6]*v[2],
		m[1]*v[0] + m[4]*v[1] + m[7]*v[2],
		m[2]*v[0] + m[5]*v[1] + m[8]*v[2],
	}
}

// Mat4 returns a Mat4 from Mat3 by Adding a row and a column with last element
// being 1.
func (m Mat3) Mat4() Mat4 {
	return Mat4{
		m[0], m[1], m[2], 0,
		m[3], m[4], m[5], 0,
		m[6], m[7], m[8], 0,
		0, 0, 0, 1,
	}
}

// Diag is a basic operation on a square matrix that simply
// returns main diagonal (meaning all elements such that row==col).
func (m Mat3) Diag() Vec3 {
	return Vec3{m[0], m[4], m[8]}
}

// M3Ident returns the 3x3 identity matrix.
// The identity matrix is a square matrix with the value 1 on its
// diagonals. The characteristic property of the identity matrix is that
// any matrix multiplied by it is itself. (MI = M; IN = N)
func M3Ident() Mat3 {
	return Mat3{1, 0, 0, 0, 1, 0, 0, 0, 1}
}

// Add performs an element-wise addition of two matrices, this is
// equivalent to iterating over every element of m and adding the corresponding value of m2.
func (m Mat3) Add(m2 Mat3) Mat3 {
	return Mat3{
		m[0] + m2[0], m[1] + m2[1], m[2] + m2[2],
		m[3] + m2[3], m[4] + m2[4], m[5] + m2[5],
		m[6] + m2[6], m[7] + m2[7], m[8] + m2[8],
	}
}

// Sub performs an element-wise subtraction of two matrices, this is
// equivalent to iterating over every element of m and subtracting the corresponding value of m2.
func (m Mat3) Sub(m2 Mat3) Mat3 {
	return Mat3{
		m[0] - m2[0], m[1] - m2[1], m[2] - m2[2],
		m[3] - m2[3], m[4] - m2[4], m[5] - m2[5],
		m[6] - m2[6], m[7] - m2[7], m[8] - m2[8],
	}
}

// MulS performs a scalar multiplcation of the matrix. This is equivalent to iterating
// over every element of the matrix and multiply it by c.
func (m Mat3) MulS(c float32) Mat3 {
	return Mat3{m[0] * c, m[1] * c, m[2] * c, m[3] * c, m[4] * c, m[5] * c, m[6] * c, m[7] * c, m[8] * c}
}

// Mul performs a "matrix product" between this matrix
// and another of the given dimension. For any two matrices of dimensionality
// MxN and NxO, the result will be MxO. For instance, Mat4 multiplied using
// Mul4x2 will result in a Mat4x2.
func (m Mat3) Mul(m2 Mat3) Mat3 {
	return Mat3{
		m[0]*m2[0] + m[3]*m2[1] + m[6]*m2[2],
		m[1]*m2[0] + m[4]*m2[1] + m[7]*m2[2],
		m[2]*m2[0] + m[5]*m2[1] + m[8]*m2[2],
		m[0]*m2[3] + m[3]*m2[4] + m[6]*m2[5],
		m[1]*m2[3] + m[4]*m2[4] + m[7]*m2[5],
		m[2]*m2[3] + m[5]*m2[4] + m[8]*m2[5],
		m[0]*m2[6] + m[3]*m2[7] + m[6]*m2[8],
		m[1]*m2[6] + m[4]*m2[7] + m[7]*m2[8],
		m[2]*m2[6] + m[5]*m2[7] + m[8]*m2[8],
	}
}

// Transpose produces the transpose of this matrix. For any MxN matrix
// the transpose is an NxM matrix with the rows swapped with the columns. For instance
// the transpose of the Mat3x2 is a Mat2x3 like so:
//
//    [[a b]]    [[a c e]]
//    [[c d]] =  [[b d f]]
//    [[e f]]
func (m Mat3) Transpose() Mat3 {
	return Mat3{
		m[0], m[3], m[6],
		m[1], m[4], m[7],
		m[2], m[5], m[8],
	}
}

// Det returns the determinant of a matrix. It is a measure of a square matrix's
// singularity and invertability, among other things. In this library, the
// determinant is hard coded based on pre-computed cofactor expansion, and uses
// no loops. Of course, the addition and multiplication must still be done.
func (m Mat3) Det() float32 {
	return m[0]*m[4]*m[8] + m[3]*m[7]*m[2] + m[6]*m[1]*m[5] - m[6]*m[4]*m[2] - m[3]*m[1]*m[8] - m[0]*m[7]*m[5]
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
// If the determinant is 0.0, this function returns the zero matrix. However,
// due to floating point errors, it is entirely plausible to get a false
// positive or negative.
// In the future, an alternate function may be written which takes in a
// pre-computed determinant.
func (m Mat3) Inv() Mat3 {
	det := m.Det()
	if FloatEqual(det, float32(0.0)) {
		return Mat3{}
	}

	retMat := Mat3{
		m[4]*m[8] - m[5]*m[7],
		m[2]*m[7] - m[1]*m[8],
		m[1]*m[5] - m[2]*m[4],
		m[5]*m[6] - m[3]*m[8],
		m[0]*m[8] - m[2]*m[6],
		m[2]*m[3] - m[0]*m[5],
		m[3]*m[7] - m[4]*m[6],
		m[1]*m[6] - m[0]*m[7],
		m[0]*m[4] - m[1]*m[3],
	}

	return retMat.MulS(1 / det)
}

// ApproxEqual performs an element-wise approximate equality test between two
// matrices, as if FloatEqual had been used.
func (m Mat3) ApproxEqual(m2 Mat3) bool {
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
func (m Mat3) ApproxEqualThreshold(m2 Mat3, threshold float32) bool {
	for i := range m {
		if !FloatEqualThreshold(m[i], m2[i], threshold) {
			return false
		}
	}
	return true
}

// ApproxFuncEqual performs an element-wise approximate equality test between two matrices
// with a given equality functions, intended to be used with FloatEqualFunc; although and comparison
// function may be used in practice.
func (m Mat3) ApproxFuncEqual(m2 Mat3, eq func(float32, float32) bool) bool {
	for i := range m {
		if !eq(m[i], m2[i]) {
			return false
		}
	}
	return true
}

// At returns the matrix element at the given row and column.
// This is equivalent to mat[col * numRow + row] where numRow is constant
// (E.G. for a Mat3x2 it's equal to 3)
//
// This method is garbage-in garbage-out. For instance, on a Mat4 asking for
// At(5,0) will work just like At(1,1). Or it may panic if it's out of bounds.
func (m Mat3) At(row, col int) float32 {
	return m[col*3+row]
}

// Set sets the corresponding matrix element at the given row and column.
// This has a pointer receiver because it mutates the matrix.
//
// This method is garbage-in garbage-out. For instance, on a Mat4 asking for
// Set(5,0,val) will work just like Set(1,1,val). Or it may panic if it's out of bounds.
func (m *Mat3) Set(row, col int, value float32) {
	m[col*3+row] = value
}

// Index returns the index of the given row and column, to be used with direct
// access. E.G. Index(0,0) = 0.
//
// This is a garbage-in garbage-out method. For instance, on a Mat4 asking for the index of
// (5,0) will work the same as asking for (1,1). Or it may give you a value that will cause
// a panic if you try to access the array with it if it's truly out of bounds.
func (m Mat3) Index(row, col int) int {
	return col*3 + row
}

// Row returns a vector representing the corresponding row (starting at row 0).
// This package makes no distinction between row and column vectors, so it
// will be a normal VecM for a MxN matrix.
func (m Mat3) Row(row int) Vec3 {
	return Vec3{m[row+0], m[row+3], m[row+6]}
}

// Rows decomposes a matrix into its corresponding row vectors.
// This is equivalent to calling mat.Row for each row.
func (m Mat3) Rows() (row0, row1, row2 Vec3) {
	return m.Row(0), m.Row(1), m.Row(2)
}

// Col returns a vector representing the corresponding column (starting at col 0).
// This package makes no distinction between row and column vectors, so it
// will be a normal VecN for a MxN matrix.
func (m Mat3) Col(col int) Vec3 {
	return Vec3{m[col*3+0], m[col*3+1], m[col*3+2]}
}

// Cols decomposes a matrix into its corresponding column vectors.
// This is equivalent to calling mat.Col for each column.
func (m Mat3) Cols() (col0, col1, col2 Vec3) {
	return m.Col(0), m.Col(1), m.Col(2)
}

// Trace is a basic operation on a square matrix that simply
// sums up all elements on the main diagonal (meaning all elements such that row==col).
func (m Mat3) Trace() float32 {
	return m[0] + m[4] + m[8]
}

// Abs returns the element-wise absolute value of this matrix
func (m Mat3) Abs() Mat3 {
	return Mat3{Abs(m[0]), Abs(m[1]), Abs(m[2]), Abs(m[3]), Abs(m[4]), Abs(m[5]), Abs(m[6]), Abs(m[7]), Abs(m[8])}
}

// String pretty prints matrix
func (m Mat3) String() string {
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 4, 4, 1, ' ', tabwriter.AlignRight)
	for i := 0; i < 3; i++ {
		for _, col := range m.Row(i) {
			fmt.Fprintf(w, "%f\t", col)
		}

		fmt.Fprintln(w, "")
	}
	w.Flush() // nolint: errcheck

	return buf.String()
}

// M3Scale returns a Scale matrix.
func M3Scale(sx, sy float32) Mat3 {
	return Mat3{
		sx, 0, 0,
		0, sy, 0,
		0, 0, 1,
	}
}

// M3Rotate returns a 2D rotation matrix based on radians 'angle'.
func M3Rotate(angle float32) Mat3 {
	// angle = (angle * math.Pi) / 180.0
	sin, cos := Sin(angle), Cos(angle)
	return Mat3{cos, sin, 0, -sin, cos, 0, 0, 0, 1}
}
