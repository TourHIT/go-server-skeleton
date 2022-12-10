package middleware

import (
    "fmt"
    "gin-server-skeleton/util"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
    "net/http"
    "os"
    "time"
)


// authApi 中间件
func AuthApiMiddle() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenHeader := c.GetHeader("Authorization")
        if tokenHeader == "" {
            c.JSON(http.StatusOK, gin.H{
                "code":  -1,
                "error": "No token",
            })
            c.Abort()
            return
        }

        tokenStr := getTokenRedis(tokenHeader)
        if tokenStr == "" {
            c.JSON(http.StatusOK, gin.H{
                "code":  -1,
                "error": "Invalid token",
            })
            c.Abort()
            return
        }
        // 解析拿到有效的 token
        jwtToken, err := parseToken(tokenStr)
        if err != nil {
            c.JSON(http.StatusOK, gin.H{
                "code":  -1,
                "error": "Invalid token parseToken",
            })
            c.Abort()
            return
        }
        // 获取 token 中的 claims
        claims, ok := jwtToken.Claims.(*AuthClaims)
        if !ok {
            c.JSON(http.StatusOK, gin.H{
                "code":  -1,
                "error": "Invalid token claims",
            })
            c.Abort()
            return
        }

        c.Set("userId", claims.UserId)
        c.Next()
        return
    }
}

// 解析token 获取jwt.Token
func parseToken(tokenStr string) (*jwt.Token, error) {
    return jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(tk *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
}

var secretKey = []byte(os.Getenv("SECRET_KEY"))
var expireTime int64 = 3600
type AuthClaims struct {
    UserId uint64 `json:"userId"`
    jwt.StandardClaims
}

// 生成token
func GenerateToken(userId uint64) (string, error) {
    claim := AuthClaims{
        UserId: userId,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Unix() + expireTime,
            IssuedAt: time.Now().Unix(),
            Issuer: "server",
            Subject: "auth",
        },
    }
    noSignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
    // 使用 secretKey 密钥进行加密处理后拿到最终 token string
    token, err := noSignedToken.SignedString(secretKey)
    if err != nil {
        fmt.Println(err.Error())
    }
    saveTokenRedis(token)
    return token, nil
}

// 将token存入redis
func saveTokenRedis(token string)  {
    util.InitRedis()
    defer util.Close()
    util.SetRedis(token, token, expireTime)
}

// 从redis中取token
func getTokenRedis(token string) string {
    util.InitRedis()
    defer util.Close()
    return util.GetRedis(token)
}






















