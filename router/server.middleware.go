package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/a5932016/go-ddd-example/config"
	"github.com/a5932016/go-ddd-example/customerror"
	"github.com/a5932016/go-ddd-example/usecase"
	"github.com/a5932016/go-ddd-example/util/mGin"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	limiterMiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
)

func (rH Handler) permissionMiddleware(pair allowancePair) mGin.HandlerFunc {
	return func(ctx *mGin.Context) {
		sid := ctx.GetHeader("Authorization")

		// No permission needed
		if pair == (allowancePair{}) {
			ctx.Set(usecase.SID, sid)
			ctx.Next()
			return
		}

		// Require Authorization
		if !(len(sid) > 0) {
			ctx.Response(http.StatusUnauthorized, "Require Authorization")
			return
		}

		// Fetch request user
		requestUser, err := rH.handler.GetRequestUserFromSID(sid)
		if err != nil {
			ctx.WithError(err).Response(http.StatusInternalServerError, "GetRequestUserFromSID")
			return
		}

		// Self Interdict Check: The account owner can not change themself
		if pair.SelfInterdictFilter {
			aimingUserID, err := getUserIDFromParam(ctx)
			if err != nil {
				ctx.WithError(err).Response(http.StatusBadRequest, "Invalid URI")
				return
			}
			if requestUser.ID == uint(aimingUserID) {
				ctx.ResponseWithCustomError(customerror.NoSelfUpdatePermission)
				return
			}
		}

		// Hierarchy Check: Regular account owner can not change root account
		if pair.HierarchyFilter {
			aimingUserID, err := getUserIDFromParam(ctx)
			if err != nil {
				ctx.WithError(err).Response(http.StatusBadRequest, "Invalid URI")
				return
			}
			user, err := rH.handler.GetUser(ctx, uint(aimingUserID))
			if err != nil {
				ctx.WithError(err).Response(http.StatusInternalServerError, "handler.GetUser")
				return
			}
			if !requestUser.IsRoot && user.IsRoot {
				ctx.ResponseWithCustomError(customerror.NoHierarchyPermission)
				return
			}
		}

		// Root user
		if requestUser.IsRoot {
			ctx.Set(usecase.SID, sid)
			ctx.Next()
			return
		}

		// Root Only Check: Only root account owner can change
		if pair.RootOnly {
			ctx.ResponseWithCustomError(customerror.NoPermission)
			return
		}

		// // Regular user permission check
		// prefixedObj := pair.Resource.Prefix()
		// prefixedAct := pair.Action.Prefix()
		// for _, division := range requestUser.Divisions {
		// 	ok, err := rH.perRepo.Enforce(division.GetPrefixedNameID(), prefixedObj, prefixedAct)
		// 	if err != nil {
		// 		ctx.WithError(err).Response(http.StatusInternalServerError,
		// 			fmt.Sprintf("enforcer.Enforce(%s, %s, %s)", division.GetPrefixedNameID(), prefixedObj, prefixedAct))
		// 		return
		// 	}
		// 	if ok {
		// 		ctx.Set(usecase.SID, sid)
		// 		ctx.Next()
		// 		return
		// 	}
		// }

		// Self Allowed Check: The account owner can change themself even without permission
		if pair.SelfPrivilege {
			aimingUserID, err := getUserIDFromParam(ctx)
			if err != nil {
				ctx.WithError(err).Response(http.StatusBadRequest, "Invalid URI")
				return
			}
			if requestUser.ID == uint(aimingUserID) {
				ctx.Set(usecase.SID, sid)
				ctx.Next()
				return
			}
		}

		ctx.ResponseWithCustomError(customerror.NoPermission)
		return
	}
}

func getUserIDFromParam(ctx *mGin.Context) (uint, error) {
	aimingUserID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(aimingUserID), nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// github.com/gin-contrib/cors
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func RateLimitMiddleware(rH Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request should be rate limited
		if shouldRateLimit(c.Request) {
			limiterMiddleware.NewMiddleware(limiter.New(rH.memRepo.GetAPILimiter(), limiter.Rate{
				Period: 1 * time.Hour,
				Limit:  1000,
			}))(c)
		}
		c.Next()
	}
}

func shouldRateLimit(req *http.Request) bool {
	if req.Header.Get("Skip-Rate-Limit") != config.Env.Core.SkipRateLimitKey {
		return true
	}
	return false
}
