package authorization

//func Test_ClaimsContext(t *testing.T) {
//	claimsContext := NewClaimsSessionContext[AccessClaims, UserSession]()
//
//	ctx := claimsContext.SetSession(context.Background(), &UserSession{Id: 1, Name: "test"})
//	session, err := claimsContext.GetSession(ctx)
//	t.Log(session, err)
//
//	ctx = claimsContext.SetClaims(context.Background(), &AccessClaims{UserInfo: &UserInfo{Id: 1, Name: "test"}})
//	claims, err := claimsContext.GetClaims(ctx)
//	t.Log(claims, err)
//
//}
