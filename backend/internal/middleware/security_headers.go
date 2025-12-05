package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersConfig configures security headers.
type SecurityHeadersConfig struct {
	// Content-Security-Policy configuration
	CSPDefaultSrc string
	CSPScriptSrc  string
	CSPStyleSrc   string
	CSPImgSrc     string
	CSPConnectSrc string
	CSPFontSrc    string
	CSPFrameSrc   string

	// Whether to enable HSTS (only for production with HTTPS)
	EnableHSTS bool
	// HSTS max-age in seconds (default 1 year)
	HSTSMaxAge int
	// Whether to include subdomains in HSTS
	HSTSIncludeSubdomains bool
	// Whether to enable HSTS preload
	HSTSPreload bool

	// X-Frame-Options value (DENY, SAMEORIGIN, or ALLOW-FROM uri)
	XFrameOptions string

	// X-Content-Type-Options (nosniff)
	XContentTypeOptions string

	// Referrer-Policy value
	ReferrerPolicy string

	// Permissions-Policy configuration
	PermissionsPolicy string
}

// DefaultSecurityHeadersConfig returns default security headers configuration.
func DefaultSecurityHeadersConfig() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		CSPDefaultSrc: "'self'",
		CSPScriptSrc:  "'self' 'unsafe-inline'",
		CSPStyleSrc:   "'self' 'unsafe-inline'",
		CSPImgSrc:     "'self' data: https:",
		CSPConnectSrc: "'self'",
		CSPFontSrc:    "'self'",
		CSPFrameSrc:   "'none'",

		EnableHSTS:            false, // Only enable in production with HTTPS
		HSTSMaxAge:            31536000,
		HSTSIncludeSubdomains: true,
		HSTSPreload:           false,

		XFrameOptions:       "DENY",
		XContentTypeOptions: "nosniff",
		ReferrerPolicy:      "strict-origin-when-cross-origin",
		PermissionsPolicy:   "geolocation=(), microphone=(), camera=()",
	}
}

// ProductionSecurityHeadersConfig returns security headers for production.
func ProductionSecurityHeadersConfig() SecurityHeadersConfig {
	cfg := DefaultSecurityHeadersConfig()
	cfg.EnableHSTS = true
	cfg.HSTSPreload = true
	return cfg
}

// SecurityHeadersMiddleware adds security headers to all responses.
func SecurityHeadersMiddleware(config SecurityHeadersConfig) gin.HandlerFunc {
	// Build CSP header
	csp := buildCSP(config)

	return func(c *gin.Context) {
		// Content-Security-Policy
		c.Header("Content-Security-Policy", csp)

		// Strict-Transport-Security (HSTS)
		if config.EnableHSTS {
			hstsValue := buildHSTS(config)
			c.Header("Strict-Transport-Security", hstsValue)
		}

		// X-Frame-Options - prevent clickjacking
		if config.XFrameOptions != "" {
			c.Header("X-Frame-Options", config.XFrameOptions)
		}

		// X-Content-Type-Options - prevent MIME type sniffing
		if config.XContentTypeOptions != "" {
			c.Header("X-Content-Type-Options", config.XContentTypeOptions)
		}

		// X-XSS-Protection - legacy XSS protection (disabled in favor of CSP)
		c.Header("X-XSS-Protection", "0")

		// Referrer-Policy
		if config.ReferrerPolicy != "" {
			c.Header("Referrer-Policy", config.ReferrerPolicy)
		}

		// Permissions-Policy (formerly Feature-Policy)
		if config.PermissionsPolicy != "" {
			c.Header("Permissions-Policy", config.PermissionsPolicy)
		}

		// Cache-Control for API responses
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Next()
	}
}

// buildCSP builds the Content-Security-Policy header value.
func buildCSP(config SecurityHeadersConfig) string {
	csp := ""

	if config.CSPDefaultSrc != "" {
		csp += "default-src " + config.CSPDefaultSrc + "; "
	}
	if config.CSPScriptSrc != "" {
		csp += "script-src " + config.CSPScriptSrc + "; "
	}
	if config.CSPStyleSrc != "" {
		csp += "style-src " + config.CSPStyleSrc + "; "
	}
	if config.CSPImgSrc != "" {
		csp += "img-src " + config.CSPImgSrc + "; "
	}
	if config.CSPConnectSrc != "" {
		csp += "connect-src " + config.CSPConnectSrc + "; "
	}
	if config.CSPFontSrc != "" {
		csp += "font-src " + config.CSPFontSrc + "; "
	}
	if config.CSPFrameSrc != "" {
		csp += "frame-src " + config.CSPFrameSrc + "; "
	}

	// Additional recommended CSP directives
	csp += "base-uri 'self'; "
	csp += "form-action 'self'; "
	csp += "frame-ancestors 'none'; "
	csp += "upgrade-insecure-requests"

	return csp
}

// buildHSTS builds the Strict-Transport-Security header value.
func buildHSTS(config SecurityHeadersConfig) string {
	hsts := "max-age=" + itoa(config.HSTSMaxAge)
	if config.HSTSIncludeSubdomains {
		hsts += "; includeSubDomains"
	}
	if config.HSTSPreload {
		hsts += "; preload"
	}
	return hsts
}

// itoa converts int to string without importing strconv
func itoa(i int) string {
	if i == 0 {
		return "0"
	}

	var b [20]byte
	pos := len(b)
	negative := i < 0
	if negative {
		i = -i
	}

	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}

	if negative {
		pos--
		b[pos] = '-'
	}

	return string(b[pos:])
}
