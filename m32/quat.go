package m32

import (
	"math"
)

// RotationOrder is the order in which rotations will be transformed for the
// purposes of QFromAngles.
type RotationOrder int

// The RotationOrder constants represent a series of rotations along the given
// axes for the use of QFromAngles.
const (
	XYX RotationOrder = iota
	XYZ
	XZX
	XZY
	YXY
	YXZ
	YZY
	YZX
	ZYZ
	ZYX
	ZXZ
	ZXY
)

// Quat short for quaternion:
// 	https://answers.unity.com/questions/467614/what-is-the-source-code-of-quaternionlookrotation.html
// W is [3]
type Quat [4]float32

// Len computes quaternion length.
func (q Quat) Len() float32 {
	return float32(math.Sqrt(float64(q[0]*q[0] + q[1]*q[1] + q[2]*q[2] + q[3]*q[3])))
}

// Normalize returns a normalized copy of the quaternion.
func (q Quat) Normalize() Quat {
	length := q.Len()

	if FloatEqual(1, length) {
		return q
	}
	if length == 0 {
		return QIdent()
	}
	if length == float32(math.Inf(1)) {
		length = MaxValue
	}

	invLen := 1 / length
	return Quat{
		q[0] * invLen,
		q[1] * invLen,
		q[2] * invLen,
		q[3] * invLen,
	}
}

// W returns the W part of the quaternion.
func (q Quat) W() float32 {
	return q[3]
}

// V returns the Vector (0,1,2) part of the quaternion.
func (q Quat) V() Vec3 {
	return Vec3{q[0], q[1], q[2]}
}

// Mul returns a new quaternion based on the multiplication with q2.
func (q Quat) Mul(q2 Quat) Quat {
	v1 := q.V()
	v2 := q2.V()

	qv := v1.Cross(v2).
		Add(v2.Mul(q[3])).
		Add(v1.Mul(q2[3]))
	return Quat{
		qv[0],
		qv[1],
		qv[2],
		q[3]*q2[3] - v1.Dot(v2),
	}
}

// Add returns a new quaternion based on the sum of q and q2.
func (q Quat) Add(q2 Quat) Quat {
	return Quat{
		q[0] + q2[0],
		q[1] + q2[1],
		q[2] + q2[2],
		q[3] + q2[3],
	}
}

// ToEuler quaternion to euler angles.
func (q Quat) ToEuler() Vec3 {
	ret := Vec3{}

	q = q.Normalize()

	sinRCosP := 2 * (q[3]*q[0] + q[1]*q[2])
	cosRCosP := 1 - 2*(q[0]*q[0]+q[1]*q[1])
	ret[0] = Atan2(sinRCosP, cosRCosP)

	sinP := 2 * (q[3]*q[1] - q[2]*q[0])
	if Abs(sinP) >= 1 {
		ret[1] = Copysign(math.Pi/2, sinP)
	} else {
		ret[1] = Asin(sinP)
	}

	sinYCosP := 2 * (q[3]*q[2] + q[0]*q[1])
	cosYCosP := 1 - 2*(q[1]*q[1]+q[2]*q[2])
	ret[2] = Atan2(sinYCosP, cosYCosP)

	return ret
}

// Mat4 returns a 4x4 matrix from the quaternion.
func (q Quat) Mat4() Mat4 {
	x, y, z, w := q[0], q[1], q[2], q[3]
	return Mat4{
		1 - 2*y*y - 2*z*z, 2*x*y + 2*w*z, 2*x*z - 2*w*y, 0,
		2*x*y - 2*w*z, 1 - 2*x*x - 2*z*z, 2*y*z + 2*w*x, 0,
		2*x*z + 2*w*y, 2*y*z - 2*w*x, 1 - 2*x*x - 2*y*y, 0,
		0, 0, 0, 1,
	}
}

// Slerp spherical linear interpolation between two quat
// https://github.com/toji/gl-matrix/blob/6c0268c89f30090b17bcadade9e7feb7205b85c5/src/quat.js#L296
func (q Quat) Slerp(b Quat, t float32) Quat {
	ax := q[0]
	ay := q[1]
	az := q[2]
	aw := q[3]

	bx := b[0]
	by := b[1]
	bz := b[2]
	bw := b[3]

	var omega, cosom, sinom, scale0, scale1 float32

	cosom = ax*bx + ay*by + az*bz + aw*bw
	if cosom < 0.0 {
		cosom = -cosom
		bx = -bx
		by = -by
		bz = -bz
		bw = -bw
	}

	if 1.0-cosom > Epsilon {
		omega = Cos(cosom)
		sinom = Sin(omega)
		scale0 = Sin((1-t)*omega) / sinom
		scale1 = Sin(t*omega) / sinom
	} else {
		scale0 = 1.0 - t
		scale1 = t
	}
	return Quat{
		scale0*ax + scale1*bx,
		scale0*ay + scale1*by,
		scale0*az + scale1*bz,
		scale0*aw + scale1*bw,
	}
}

// QFromAngles returns a rotation quaternion based on the angles.
func QFromAngles(a1, a2, a3 float32, order RotationOrder) Quat {
	var s [3]float64
	var c [3]float64

	s[0], c[0] = math.Sincos(float64(a1 / 2))
	s[1], c[1] = math.Sincos(float64(a2 / 2))
	s[2], c[2] = math.Sincos(float64(a3 / 2))

	switch order {
	case ZYX:
		return Quat{
			float32(c[0]*c[1]*s[2] - s[0]*s[1]*c[2]),
			float32(c[0]*s[1]*c[2] + s[0]*c[1]*s[2]),
			float32(s[0]*c[1]*c[2] - c[0]*s[1]*s[2]),
			float32(c[0]*c[1]*c[2] + s[0]*s[1]*s[2]),
		}
	case ZYZ:
		return Quat{
			float32(c[0]*s[1]*s[2] - s[0]*s[1]*c[2]),
			float32(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
			float32(s[0]*c[1]*c[2] + c[0]*c[1]*s[2]),
			float32(c[0]*c[1]*c[2] - s[0]*c[1]*s[2]),
		}
	case ZXY:
		return Quat{
			float32(c[0]*s[1]*c[2] - s[0]*c[1]*s[2]),
			float32(c[0]*c[1]*s[2] + s[0]*s[1]*c[2]),
			float32(c[0]*s[1]*s[2] + s[0]*c[1]*c[2]),
			float32(c[0]*c[1]*c[2] - s[0]*s[1]*s[2]),
		}
	case ZXZ:
		return Quat{
			float32(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
			float32(s[0]*s[1]*c[2] - c[0]*s[1]*s[2]),
			float32(c[0]*c[1]*s[2] + s[0]*c[1]*c[2]),
			float32(c[0]*c[1]*c[2] - s[0]*c[1]*s[2]),
		}
	case YXZ:
		return Quat{
			float32(c[0]*s[1]*c[2] + s[0]*c[1]*s[2]),
			float32(s[0]*c[1]*c[2] - c[0]*s[1]*s[2]),
			float32(c[0]*c[1]*s[2] - s[0]*s[1]*c[2]),
			float32(c[0]*c[1]*c[2] + s[0]*s[1]*s[2]),
		}
	case YXY:
		return Quat{
			float32(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
			float32(s[0]*c[1]*c[2] + c[0]*c[1]*s[2]),
			float32(c[0]*s[1]*s[2] - s[0]*s[1]*c[2]),
			float32(c[0]*c[1]*c[2] - s[0]*c[1]*s[2]),
		}
	case YZX:
		return Quat{
			float32(c[0]*c[1]*s[2] + s[0]*s[1]*c[2]),
			float32(c[0]*s[1]*s[2] + s[0]*c[1]*c[2]),
			float32(c[0]*s[1]*c[2] - s[0]*c[1]*s[2]),
			float32(c[0]*c[1]*c[2] - s[0]*s[1]*s[2]),
		}
	case YZY:
		return Quat{
			float32(s[0]*s[1]*c[2] - c[0]*s[1]*s[2]),
			float32(c[0]*c[1]*s[2] + s[0]*c[1]*c[2]),
			float32(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
			float32(c[0]*c[1]*c[2] - s[0]*c[1]*s[2]),
		}
	case XYZ:
		return Quat{
			float32(c[0]*s[1]*s[2] + s[0]*c[1]*c[2]),
			float32(c[0]*s[1]*c[2] - s[0]*c[1]*s[2]),
			float32(c[0]*c[1]*s[2] + s[0]*s[1]*c[2]),
			float32(c[0]*c[1]*c[2] - s[0]*s[1]*s[2]),
		}
	case XYX:
		return Quat{
			float32(c[0]*c[1]*s[2] + s[0]*c[1]*c[2]),
			float32(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
			float32(s[0]*s[1]*c[2] - c[0]*s[1]*s[2]),
			float32(c[0]*c[1]*c[2] - s[0]*c[1]*s[2]),
		}
	case XZY:
		return Quat{
			float32(s[0]*c[1]*c[2] - c[0]*s[1]*s[2]),
			float32(c[0]*c[1]*s[2] - s[0]*s[1]*c[2]),
			float32(c[0]*s[1]*c[2] + s[0]*c[1]*s[2]),
			float32(c[0]*c[1]*c[2] + s[0]*s[1]*s[2]),
		}
	case XZX:
		return Quat{
			float32(c[0]*c[1]*s[2] + s[0]*c[1]*c[2]),
			float32(c[0]*s[1]*s[2] - s[0]*s[1]*c[2]),
			float32(c[0]*s[1]*c[2] + s[0]*s[1]*s[2]),
			float32(c[0]*c[1]*c[2] - s[0]*c[1]*s[2]),
		}
	default:
		panic("Unsupported rotation order")
	}
}

// QLookAt creates a LookAt Matrix and extract the quaternion
func QLookAt(dir, up Vec3) Quat {
	mat := LookAt(dir, Vec3{0, 0, 0}, up)
	var (
		m00, m01, m02 = mat[0], mat[4], mat[8]
		m10, m11, m12 = mat[1], mat[5], mat[9]
		m20, m21, m22 = mat[2], mat[6], mat[10]
	)

	fourXSquaredMinus1 := m00 - m11 - m22
	fourYSquaredMinus1 := m11 - m00 - m22
	fourZSquaredMinus1 := m22 - m00 - m11
	fourWSquaredMinus1 := m00 + m11 + m22

	biggestIndex := 0
	fourBiggestSquaredMinus1 := fourWSquaredMinus1
	if fourXSquaredMinus1 > fourBiggestSquaredMinus1 {
		fourBiggestSquaredMinus1 = fourXSquaredMinus1
		biggestIndex = 1
	}
	if fourYSquaredMinus1 > fourBiggestSquaredMinus1 {
		fourBiggestSquaredMinus1 = fourYSquaredMinus1
		biggestIndex = 2
	}
	if fourZSquaredMinus1 > fourBiggestSquaredMinus1 {
		fourBiggestSquaredMinus1 = fourZSquaredMinus1
		biggestIndex = 3
	}

	biggestVal := Sqrt(fourBiggestSquaredMinus1+1) * 0.5
	mult := 0.25 / biggestVal

	switch biggestIndex {
	case 0:
		return Quat{
			(m12 - m21) * mult, (m20 - m02) * mult, (m01 - m10) * mult,
			biggestVal,
		}
	case 1:
		return Quat{
			biggestVal, (m01 + m10) * mult, (m20 + m02) * mult,
			(m12 - m21) * mult,
		}
	case 2:
		return Quat{
			(m01 + m10) * mult, biggestVal, (m12 + m21) * mult,
			(m20 - m02) * mult,
		}
	case 3:
		return Quat{
			(m20 + m02) * mult, (m12 + m21) * mult, biggestVal,
			(m01 - m10) * mult,
		}
	}
	return QIdent()
}

// QIdent returns a quaternion identity.
func QIdent() Quat {
	return Quat{0, 0, 0, 1}
}

// QEuler returns a quaternion based on those euler angles
// https://en.wikipedia.org/wiki/Conversion_between_quaternions_and_Euler_angles
func QEuler(x, y, z float32) Quat {
	cy, sy := Sincos(z * 0.5)
	cp, sp := Sincos(y * 0.5)
	cr, sr := Sincos(x * 0.5)

	return Quat{
		cy*cp*sr - sy*sp*cr,
		sy*cp*sr + cy*sp*cr,
		sy*cp*cr - cy*sp*sr,
		cy*cp*cr + sy*sp*sr,
	}
}

// QAxisAngle returns a axis angle quaternion.
func QAxisAngle(v3 Vec3, rad float32) Quat {
	s, c := Sincos(rad * 0.5)

	return Quat{
		v3[0] * s,
		v3[1] * s,
		v3[2] * s,
		c,
	}.Normalize()
}

// QBetweenV3 returns a quaternion between two vectors.
func QBetweenV3(v1, v2 Vec3) Quat {
	var q Quat
	a := v1.Normalize().Cross(v2.Normalize())
	q[0] = a[0]
	q[1] = a[1]
	q[2] = a[2]
	q[3] = Sqrt(1 + v1.Dot(v2))
	return q.Normalize()
	// return QAxisAngle(v1.Cross(v2), v1.Dot(v2))
}
