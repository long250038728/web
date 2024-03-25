package tool

//=========================== app ================================

var Address = NewError("10000001", "IP / Address Not Find")
var JobSubTime = NewError("10000002", "SubTime Is Error")

// ===================== limiter / Cache ==========================

var LimiterErr = NewError("10001001", "API Limiter")

//=========================== auth ================================

var Token = NewError("10002001", "Token Disabled")
