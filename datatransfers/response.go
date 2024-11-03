package datatransfers

type Response struct {
	Code    int         `json:"-"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}
