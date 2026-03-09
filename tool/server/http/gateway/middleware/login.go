package middleware

//// Login 校验登录
//func Login() gateway.ServerInterceptor {
//	return func(ctx context.Context, requestInfo map[string]any, request any, handler gateway.Handler) (resp any, err error) {
//		if _, err = authorization.NewClaimsSessionContext[authorization.AccessClaims, authorization.UserSession]().GetClaims(ctx); err != nil {
//			return nil, err
//		}
//		return handler(ctx, request)
//	}
//}
