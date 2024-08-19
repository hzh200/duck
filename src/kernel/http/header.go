package http

import netUrl "net/url"

const (
	Accept                        = "Accept"
	AcceptCharset                 = "Accept-Charset"
	AcceptEncoding                = "Accept-Encoding"
	AcceptLanguage                = "Accept-Language"
	Authorization                 = "Authorization"
	CacheControl                  = "Cache-Control"
	ContentLength                 = "Content-Length"
	ContentMD5                    = "Content-MD5"
	ContentType                   = "Content-Type"
	Connection                    = "Connection"
	DoNotTrack                    = "DNT"
	IfMatch                       = "If-Match"
	IfModifiedSince               = "If-Modified-Since"
	IfNoneMatch                   = "If-None-Match"
	IfRange                       = "If-Range"
	IfUnmodifiedSince             = "If-Unmodified-Since"
	MaxForwards                   = "Max-Forwards"
	ProxyAuthorization            = "Proxy-Authorization"
	Pragma                        = "Pragma"
	Range                         = "Range"
	Referer                       = "Referer"
	Host                          = "Host"
	UserAgent                     = "User-Agent"
	TE                            = "TE"
	Via                           = "Via"
	Warning                       = "Warning"
	Cookie                        = "Cookie"
	Origin                        = "Origin"
	AcceptDatetime                = "Accept-Datetime"
	XRequestedWith                = "X-Requested-With"
	AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	AccessControlAllowMethods     = "Access-Control-Allow-Methods"
	AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	AccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	AccessControlMaxAge           = "Access-Control-Max-Age"
	AccessControlRequestMethod    = "Access-Control-Request-Method"
	AccessControlRequestHeaders   = "Access-Control-Request-Headers"
	AcceptPatch                   = "Accept-Patch"
	AcceptRanges                  = "Accept-Ranges"
	Allow                         = "Allow"
	ContentEncoding               = "Content-Encoding"
	ContentLanguage               = "Content-Language"
	ContentLocation               = "Content-Location"
	ContentDisposition            = "Content-Disposition"
	ContentRange                  = "Content-Range"
	ETag                          = "ETag"
	Expires                       = "Expires"
	LastModified                  = "Last-Modified"
	Link                          = "Link"
	Location                      = "Location"
	P3P                           = "P3P"
	ProxyAuthenticate             = "Proxy-Authenticate"
	Refresh                       = "Refresh"
	RetryAfter                    = "Retry-After"
	Server                        = "Server"
	SetCookie                     = "Set-Cookie"
	StrictTransportSecurity       = "Strict-Transport-Security"
	TransferEncoding              = "Transfer-Encoding"
	Upgrade                       = "Upgrade"
	Vary                          = "Vary"
	WWWAuthenticate               = "WWW-Authenticate"

	// Non-Standard
	XFrameOptions          = "X-Frame-Options"
	XXSSProtection         = "X-XSS-Protection"
	ContentSecurityPolicy  = "Content-Security-Policy"
	XContentSecurityPolicy = "X-Content-Security-Policy"
	XWebKitCSP             = "X-WebKit-CSP"
	XContentTypeOptions    = "X-Content-Type-Options"
	XPoweredBy             = "X-Powered-By"
	XUACompatible          = "X-UA-Compatible"
	XForwardedProto        = "X-Forwarded-Proto"
	XHTTPMethodOverride    = "X-HTTP-Method-Override"
	XForwardedFor          = "X-Forwarded-For"
	XRealIP                = "X-Real-IP"
	XCSRFToken             = "X-CSRF-Token"
	XRatelimitLimit        = "X-Ratelimit-Limit"
	XRatelimitRemaining    = "X-Ratelimit-Remaining"
	XRatelimitReset        = "X-Ratelimit-Reset"
)

// text/html; charset=utf-8
const MimetypeReStr string = `([a-zA-Z]+)\/([a-zA-Z\-]+)(?:; charset=(.+))?`
// Last-Modified: Wed, 21 Oct 2015 07:28:00 GMT
// Last-Modified: <day-name>, <day> <month> <year> <hour>:<minute>:<second> GMT
const LastModifiedReStr string = `([a-zA-Z]+), ([0-9]{1, 2}) ([a-zA-Z]+) ([0-9]{4}) ([0-9]{2}):([0-9]{2}):([0-9]{2}) ([a-zA-Z]+)`

func ParseHost(url string) (string, error) {
	parsedUrlInfo, err := netUrl.Parse(url)
	if err != nil {
		return "", err
	}
	return parsedUrlInfo.Host, nil
}

func RequestHeaders(url string) (map[string]string, error) {
	headers := map[string]string{
		Accept: "*/*",
		AcceptEncoding: "gzip, deflate, br",
		Connection: "keep-alive",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.148 Safari/537.36",
	}

	host, err := ParseHost(url)
	if err != nil {
		return nil, err
	}

	headers[Host] = host
	
	return headers, nil
}
