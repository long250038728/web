package authorization

////================================ ctx ================================
//
//type ClaimsSessionContext[C Claims, S any] interface {
//	SetClaims(ctx context.Context, claims *C) context.Context
//	GetClaims(ctx context.Context) (*C, error)
//	SetSession(ctx context.Context, session *S) context.Context
//	GetSession(ctx context.Context) (*S, error)
//}
//
//type claimsKey struct{}
//type sessionKey struct{}
//
//func NewClaimsSessionContext[C Claims, S any]() ClaimsSessionContext[C, S] {
//	return &claimsSessionContext[C, S]{}
//}
//
//type claimsSessionContext[C Claims, S any] struct {
//	claims  *C
//	session *S
//}
//
//func (c *claimsSessionContext[C, S]) SetClaims(ctx context.Context, claims *C) context.Context {
//	return context.WithValue(ctx, claimsKey{}, claims)
//}
//
//func (c *claimsSessionContext[C, S]) GetClaims(ctx context.Context) (*C, error) {
//	if val, ok := ctx.Value(claimsKey{}).(*C); ok {
//		return val, nil
//	}
//	return nil, app_error.ClaimsNull
//}
//
//func (c *claimsSessionContext[C, S]) SetSession(ctx context.Context, session *S) context.Context {
//	return context.WithValue(ctx, sessionKey{}, session)
//}
//
//func (c *claimsSessionContext[C, S]) GetSession(ctx context.Context) (*S, error) {
//	if val, ok := ctx.Value(sessionKey{}).(*S); ok {
//		return val, nil
//	}
//	return nil, app_error.ClaimsNull
//}
