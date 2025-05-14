package main

// statusMapping maps HTTP status codes to their corresponding text and emoji.
//
// example:
//   - 200: {"OK", "😃"}
var statusMapping = map[int]struct {
	Text  string
	Emoji string
}{
	100: {"Continue", "⏩"},
	101: {"Switching Protocols", "🔄"},
	102: {"Processing", "🔍"},
	200: {"OK", "😃"},
	201: {"Created", "🆕"},
	202: {"Accepted", "👌"},
	203: {"Non-Authoritative Information", "🔎"},
	204: {"No Content", "🙅‍♂️"},
	205: {"Reset Content", "🔄"},
	206: {"Partial Content", "🔍"},
	207: {"Multi-Status", "🔎"},
	208: {"Already Reported", "🙅‍♂️"},
	226: {"IM Used", "🔎"},
	300: {"Multiple Choices", "🔎"},
	301: {"Moved Permanently", "🔀"},
	302: {"Found", "🔍"},
	303: {"See Other", "🔎"},
	304: {"Not Modified", "🙅‍♂️"},
	305: {"Use Proxy", "🔄"},
	307: {"Temporary Redirect", "🔄"},
	308: {"Permanent Redirect", "🔄"},
	400: {"Bad Request", "📄"},
	401: {"Unauthorized", "🔒"},
	402: {"Payment Required", "💰"},
	403: {"Forbidden", "⛔️"},
	404: {"Not Found", "👀"},
	405: {"Method Not Allowed", "❌"},
	406: {"Not Acceptable", "🚫"},
	407: {"Proxy Authentication Required", "🛂"},
	408: {"Request Timeout", "⏳"},
	409: {"Conflict", "⚔️"},
	410: {"Gone", "👋"},
	411: {"Length Required", "📏"},
	412: {"Precondition Failed", "⚠️"},
	413: {"Payload Too Large", "📦"},
	414: {"URI Too Long", "🔗"},
	415: {"Unsupported Media Type", "🎥"},
	416: {"Range Not Satisfiable", "🔍"},
	417: {"Expectation Failed", "❗️"},
	418: {"I'm a teapot", "💻"},
	421: {"Misdirected Request", "↪️"},
	422: {"Unprocessable Entity", "🧩"},
	423: {"Locked", "🔐"},
	424: {"Failed Dependency", "🔗"},
	425: {"Too Early", "⏰"},
	426: {"Upgrade Required", "🔄"},
	428: {"Precondition Required", "⚙️"},
	429: {"Too Many Requests", "🚨"},
	431: {"Request Header Fields Too Large", "📋"},
	451: {"Unavailable For Legal Reasons", "⚖️"},
	500: {"Internal Server Error", "💥"},
	501: {"Not Implemented", "🚧"},
	502: {"Bad Gateway", "😳"},
	503: {"Service Unavailable", "🔧"},
	504: {"Gateway Timeout", "⌛️"},
	505: {"HTTP Version Not Supported", "🔄"},
	506: {"Variant Also Negotiates", "🌀"},
	507: {"Insufficient Storage", "💾"},
	508: {"Loop Detected", "🔁"},
	510: {"Not Extended", "➕"},
	511: {"Network Authentication Required", "🔐"},
}

const defaultText = "Something went wrong... Please try again later."
const defaultEmoji = "😓"

// getStatusInfo returns the text and emoji corresponding to the HTTP status code.
//
// params:
//   - status: the HTTP status code
//
// returns:
//   - the text and emoji corresponding to the HTTP status code
func getStatusInfo(status int) (string, string) {
	if info, ok := statusMapping[status]; ok {
		return info.Text, info.Emoji
	}
	return defaultText, defaultEmoji
}

// sanitizeStatusCode ensures the provided HTTP status code is within the valid range (100–599).
// If the code is outside the valid range, it defaults to 418 (I'm a teapot).
//
// params:
//   - code: the HTTP status code to validate
//
// returns:
//   - a valid HTTP status code (between 100 and 599 inclusive, or 418 as fallback)
func sanitizeStatusCode(code int) int {
	if code < 100 || code > 599 {
		return 418 // I'm a teapot
	}
	return code
}
