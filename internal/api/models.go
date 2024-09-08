package api

type EllipticArgs struct {
	// Coefficients of the curve equation y^2 = x^3 + Ax + B
	A int `json:"a" form:"a" binding:"required"`
	B int `json:"b" form:"b" binding:"required"`
	P int `json:"p" form:"p" binding:"required"` // Prime modulus
}

type EncryptData struct {
	Text string `json:"text"  form:"text"`
}

type DecryptData struct {
	Text    []byte `json:"text"  form:"text"`
	Private string `json:"private"  form:"private"`
}

type ResultData struct {
	Text string `json:"text" binding:"required" form:"text"`
}

type Keys struct {
	Public  string `json:"public"  form:"public"`
	Private string `json:"private"  form:"private"`
}

type PublicKey struct {
	X int64 `json:"x" form:"x"`
	Y int64 `json:"y" form:"y"`
}
