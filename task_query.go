package qencode

// StartTaskQuery represents the query form parameter
// as specified in https://docs.qencode.com/api-reference/transcoding/#starting-a-task.
type StartTaskQuery struct {
	Query Query `json:"query"`
}

// Query to start a task.
type Query struct {
	Format         []Format `json:"format"`
	EncoderVersion string   `json:"encoder_version"`
	Source         string   `json:"source"`
	CallbackURL    string   `json:"callback_url"`
}

// Format of the encoded outputs.
type Format struct {
	Output          string      `json:"output"`
	SeparateAudio   int         `json:"separate_audio"`
	VideoCodec      string      `json:"video_codec,omitempty"`
	AudioBitrate    string      `json:"audio_bitrate,omitempty"`
	SegmentDuration int         `json:"segment_duration"`
	Destination     Destination `json:"destination"`
	Stream          []Stream    `json:"stream"`
}

// Destination to where the encoded files will be sent
// after the task completion.
type Destination struct {
	URL    string `json:"url"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

// Stream represents an adaptive stream format.
type Stream struct {
	VideoCodec      string `json:"video_codec"`
	Height          int    `json:"height"`
	AudioBitrate    int    `json:"audio_bitrate"`
	OptimizeBitrate int    `json:"optimize_bitrate"`
	ChunkListName   string `json:"chunklist_name,omitempty"`
}

const (
	// OutputHLS represents the hls file format encoding.
	OutputHLS = "advanced_hls"
)

const (
	//SeparateAudioDisabled disables the audio separation.
	SeparateAudioDisabled = iota
	//SeparateAudioEnabled enables the audio separation.
	SeparateAudioEnabled
)

const (
	//OptimizeBitRateDisabled disables the bitrate optmization.
	OptimizeBitRateDisabled = iota
	//OptimizeBitRateEnabled enables the bitrate optimization.
	OptimizeBitRateEnabled
)
