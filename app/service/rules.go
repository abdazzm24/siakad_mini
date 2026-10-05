package service

import "math"

// File ini berisi business rules MURNI: tidak menyentuh Fiber maupun database,
// sehingga dapat diuji hanya dengan memanggil fungsinya.

// MaxSKS menentukan batas SKS berdasarkan IPK terakhir.
//
//	IPK >= 3.00        -> 24 SKS
//	2.50 <= IPK < 3.00 -> 21 SKS
//	IPK < 2.50         -> 18 SKS
func MaxSKS(ipk float64) int {
	// Dibulatkan ke sen agar perbandingan tidak terganggu pecahan biner.
	cents := int(math.Round(ipk * 100))
	switch {
	case cents >= 300:
		return 24
	case cents >= 250:
		return 21
	default:
		return 18
	}
}

// RemainingSKS menghitung sisa SKS yang masih boleh diambil (tidak pernah negatif).
func RemainingSKS(maxSKS, current int) int {
	return max(maxSKS-current, 0)
}

// ExceedsSKSLimit true bila menambah newSKS membuat total melebihi batas.
func ExceedsSKSLimit(maxSKS, current, newSKS int) bool {
	return current+newSKS > maxSKS
}

// IsCourseFull true bila kuota sudah terpenuhi (terisi >= kuota).
func IsCourseFull(terisi, kuota int) bool {
	return terisi >= kuota
}