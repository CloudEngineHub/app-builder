// Copyright (c) 2024 Baidu, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package appbuilder

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

const (
	CodeContentType              = "code"
	TextContentType              = "text"
	ImageContentType             = "image"
	RAGContentType               = "rag"
	FunctionCallContentType      = "function_call"
	AudioContentType             = "audio"
	VideoContentType             = "video"
	StatusContentType            = "status"
	ChatflowInterruptContentType = "chatflow_interrupt"
	PublishMessageContentType    = "publish_message"
	JsonContentType              = "json"
	ChatReasoningContentType     = "chat_reasoning"
)

const (
	ChatflowEventType      = "chatflow"
	FollowUpQueryEventType = "FollowUpQuery"
)

var TypeToStruct = map[string]reflect.Type{
	CodeContentType:              reflect.TypeOf(CodeDetail{}),
	TextContentType:              reflect.TypeOf(TextDetail{}),
	ImageContentType:             reflect.TypeOf(ImageDetail{}),
	RAGContentType:               reflect.TypeOf(RAGDetail{}),
	FunctionCallContentType:      reflect.TypeOf(FunctionCallDetail{}),
	AudioContentType:             reflect.TypeOf(AudioDetail{}),
	VideoContentType:             reflect.TypeOf(VideoDetail{}),
	StatusContentType:            reflect.TypeOf(StatusDetail{}),
	ChatflowInterruptContentType: reflect.TypeOf(ChatflowInterruptDetail{}),
	PublishMessageContentType:    reflect.TypeOf(PublishMessageDetail{}),
	JsonContentType:              reflect.TypeOf(JsonDetail{}),
	ChatReasoningContentType:     reflect.TypeOf(ChatReasoningDetail{}),
}

type AppBuilderClientRunRequest struct {
	AppID            string                    `json:"app_id"`
	Query            string                    `json:"query"`
	Stream           bool                      `json:"stream"`
	EndUserID        *string                   `json:"end_user_id"`
	ConversationID   string                    `json:"conversation_id"`
	FileIDs          []string                  `json:"file_ids"`
	Tools            []Tool                    `json:"tools"`
	ToolOutputs      []ToolOutput              `json:"tool_outputs"`
	ToolChoice       *ToolChoice               `json:"tool_choice"`
	Action           *Action                   `json:"action"`
	McpAuthorization *[]map[string]interface{} `json:"mcp_authorization,omitempty"`
}

type AppBuilderClientUploadFileRequest struct {
	AppID          string `json:"app_id"`
	ConversationID string `json:"conversation_id"`
	FilePath       string `json:"file_path"`
	FileURL        string `json:"file_url"`
}

type AppBuilderClientFeedbackRequest struct {
	AppID          string   `json:"app_id"`
	ConversationID string   `json:"conversation_id"`
	MessageID      string   `json:"message_id"`
	Type           string   `json:"type"`
	Flag           []string `json:"flag,omitempty"`
	Reason         string   `json:"reason,omitempty"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type ToolOutput struct {
	ToolCallID string `json:"tool_call_id" description:"工具调用ID"`
	Output     string `json:"output" description:"工具输出"`
}

type ToolChoice struct {
	Type     string             `json:"type"`
	Function ToolChoiceFunction `json:"function"`
}

type ToolChoiceFunction struct {
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

type Action struct {
	ActionType string           `json:"action_type"`
	Paramters  *ActionParamters `json:"parameters"`
}

type ActionParamters struct {
	InterruptEvent *ActionInterruptEvent `json:"interrupt_event"`
}

type ActionInterruptEvent struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func NewResumeAction(eventId string) *Action {
	return NewAction("resume", eventId, "chat")
}

func NewAction(actionType string, eventId string, eventType string) *Action {
	return &Action{
		ActionType: actionType,
		Paramters: &ActionParamters{
			InterruptEvent: &ActionInterruptEvent{
				ID:   eventId,
				Type: eventType,
			},
		},
	}
}

type AgentBuilderRawResponse struct {
	RequestID      string           `json:"request_id"`
	Date           string           `json:"date"`
	Answer         string           `json:"answer"`
	ConversationID string           `json:"conversation_id"`
	MessageID      string           `json:"message_id"`
	IsCompletion   bool             `json:"is_completion"`
	Content        []RawEventDetail `json:"content"`
}

type RawEventDetail struct {
	EventCode    int             `json:"event_code"`
	EventMessage string          `json:"event_message"`
	EventType    string          `json:"event_type"`
	EventID      string          `json:"event_id"`
	EventStatus  string          `json:"event_status"`
	ContentType  string          `json:"content_type"`
	Outputs      json.RawMessage `json:"outputs"`
	Usage        Usage           `json:"usage"`
	ToolCalls    []ToolCall      `json:"tool_calls"`
}

type Usage struct {
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	TotalTokens      int    `json:"total_tokens"`
	Name             string `json:"name"`
}

type AgentBuilderAnswer struct {
	Answer string
	Events []Event
}

type Event struct {
	Code        int
	Message     string
	Status      string
	EventType   string
	ContentType string
	Usage       Usage
	Detail      any
	ToolCalls   []ToolCall
}

type ToolCall struct {
	ID       string             `json:"id"`       // 工具调用ID
	Type     string             `json:"type"`     // 需要输出的工具调用的类型。就目前而言，这始终是function
	Function FunctionCallOption `json:"function"` // 函数定义
}

type FunctionCallOption struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type TextDetail struct {
	Text string `json:"text"`
}

type CodeDetail struct {
	Text  string   `json:"text"`
	Code  string   `json:"code"`
	Files []string `json:"files"`
}

type RAGDetail struct {
	Text       string      `json:"text"`
	References []Reference `json:"references"`
}

type Reference struct {
	ID              string `json:"id"`
	From            string `json:"from"`
	URL             string `json:"url"`
	Content         string `json:"content"`
	SegmentID       string `json:"segment_id"`
	DocumentID      string `json:"document_id"`
	DatasetID       string `json:"dataset_id"`
	DocumentName    string `json:"document_name"`
	KnowledgeBaseID string `json:"knowledgebase_id"`
}

type FunctionCallDetail struct {
	Text  any    `json:"text"`
	Image string `json:"image"`
	Audio string `json:"audio"`
	Video string `json:"video"`
}

type ImageDetail struct {
	Image string `json:"image"`
}

type AudioDetail struct {
	Audio string `json:"audio"`
}

type VideoDetail struct {
	Video string `json:"video"`
}

type StatusDetail struct{}

type ChatflowInterruptDetail struct {
	InterruptEventID   string `json:"interrupt_event_id"`
	InterruptEventType string `json:"interrupt_event_type"`
}

type PublishMessageDetail struct {
	Message   string `json:"message"`
	MessageID string `json:"message_id"`
}

type ChatReasoningDetail struct {
	Text string `json:"text"`
}

type JsonDetail struct {
	Json FollowUpQueries `json:"json"`
}

type FollowUpQueries struct {
	FollowUpQueries []string `json:"follow_up_querys"`
}

type DefaultDetail struct {
	Text  string   `json:"text,omitempty"`
	URLS  []string `json:"urls,omitempty"`
	Files []string `json:"files,omitempty"`
	Image string   `json:"image,omitempty"`
	Video string   `json:"video,omitempty"`
	Audio string   `json:"audio,omitempty"`
}

type AppBuilderClientRawResponse struct {
	RequestID      string           `json:"request_id"`
	Date           string           `json:"date"`
	Answer         string           `json:"answer"`
	ConversationID string           `json:"conversation_id"`
	MessageID      string           `json:"message_id"`
	IsCompletion   bool             `json:"is_completion"`
	Content        []RawEventDetail `json:"content"`
	Code           string           `json:"code,omitempty"`
	Message        string           `json:"message,omitempty"`
}

type GetAppListRequest struct {
	Limit  int    `json:"limit"`
	After  string `json:"after"`
	Before string `json:"before"`
}

type GetAppListResponse struct {
	RequestID string `json:"request_id"`
	Data      []App  `json:"data"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

type DescribeAppsRequest struct {
	Marker  *string `json:"marker,omitempty"`
	MaxKeys *int    `json:"maxKeys,omitempty"`
}

type DescribeAppsResponse struct {
	RequestID   string `json:"requestId"`
	Marker      string `json:"marker"`
	IsTruncated bool   `json:"isTruncated"`
	NextMarker  string `json:"nextMarker"`
	MaxKeys     int    `json:"maxKeys"`
	Data        []App  `json:"data"`
}

type DescribeAppRequest struct {
	ID string `json:"id"`
}

type DescribeAppResponse struct {
	RequestID           string              `json:"requestId"`
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	Instruction         string              `json:"instruction"`
	Prologue            string              `json:"prologue"`
	ExampleQueries      []string            `json:"exampleQueries"`
	FollowUpQueries     AppFollowUpQueries  `json:"followUpQueries"`
	Components          []Component         `json:"components"`
	KnowledgeBaseConfig KnowledgeBaseConfig `json:"knowledgeBaseConfig"`
	ModelConfig         ModelConfig         `json:"modelConfig"`
	Background          *Background         `json:"background,omitempty"`
}
type AppFollowUpQueries struct {
	Type   string `json:"type"`
	Prompt string `json:"prompt"`
	Round  string `json:"round"`
}
type Component struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	CustomDesc  string `json:"customDesc,omitempty"`
}
type KnowledgeBaseConfig struct {
	KnowledgeBases []AppKnowledgeBase `json:"knowledgeBases"`
	Retrieval      RetrievalConfig    `json:"retrieval"`
}
type AppKnowledgeBase struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
type RetrievalConfig struct {
	EnableWebSearch bool    `json:"enableWebSearch,omitempty"`
	Order           string  `json:"order,omitempty"`
	Strategy        string  `json:"strategy,omitempty"`
	TopK            int     `json:"topK,omitempty"`
	Threshold       float64 `json:"threshold,omitempty"`
}
type ModelConfig struct {
	Plan PlanConfig `json:"plan"`
	Chat ChatConfig `json:"chat"`
}
type PlanConfig struct {
	ModelID   string      `json:"modelId"`
	Model     string      `json:"model"`
	MaxRounds int         `json:"maxRounds"`
	Config    ModelParams `json:"config"`
}
type ChatConfig struct {
	ModelID           string      `json:"modelId"`
	Model             string      `json:"model"`
	HistoryChatRounds int         `json:"historyChatRounds"`
	Config            ModelParams `json:"config"`
}
type ModelParams struct {
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"topP"`
}
type Background struct {
	ID           string        `json:"id,omitempty"`
	Path         string        `json:"path,omitempty"`
	MobileConfig *MobileConfig `json:"mobile_config,omitempty"`
	PCConfig     *PCConfig     `json:"pc_config,omitempty"`
}
type MobileConfig struct {
	Left   string `json:"left,omitempty"`
	Top    string `json:"top,omitempty"`
	Height string `json:"height,omitempty"`
	Color  string `json:"color,omitempty"`
}
type PCConfig struct {
	Left   string `json:"left,omitempty"`
	Top    string `json:"top,omitempty"`
	Height string `json:"height,omitempty"`
	Color  string `json:"color,omitempty"`
}

type App struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	AppType     string `json:"appType,omitempty"`
	IsPublished bool   `json:"isPublished,omitempty"`
	UpdateTime  int64  `json:"updateTime,omitempty"`
}

type AppBuilderClientAnswer struct {
	MessageID string
	Answer    string
	Events    []Event
	Code      string
	Message   string
	RequestID string
}

func (t *AppBuilderClientAnswer) transform(inp *AppBuilderClientRawResponse) {
	t.Answer = inp.Answer
	t.MessageID = inp.MessageID
	t.Code = inp.Code
	t.Message = inp.Message
	t.RequestID = inp.RequestID
	for _, c := range inp.Content {
		ev := Event{Code: c.EventCode,
			Message:     c.EventMessage,
			Status:      c.EventStatus,
			EventType:   c.EventType,
			ContentType: c.ContentType,
			Usage:       c.Usage,
			Detail:      c.Outputs,
			ToolCalls:   c.ToolCalls}
		// 这部分新改的
		tp, ok := TypeToStruct[ev.ContentType]
		if !ok {
			tp = reflect.TypeOf(DefaultDetail{})
		}
		v := reflect.New(tp)
		_ = json.Unmarshal(c.Outputs, v.Interface())
		ev.Detail = v.Elem().Interface()
		// 这部分新改的
		t.Events = append(t.Events, ev)
	}
}

// AppBuilderClientIterator 定义AppBuilderClient流式/非流式迭代器接口
// 初始状态可迭代,如果返回error不为空则代表迭代结束，
// error为io.EOF，则代表迭代正常结束，其它则为异常结束
type AppBuilderClientIterator interface {
	// Next 获取处理结果，如果返回error不为空，迭代器自动失效，不允许再调用此方法
	Next() (*AppBuilderClientAnswer, error)
}

type AppBuilderClientStreamIterator struct {
	requestID string
	r         *sseReader
	body      io.ReadCloser
}

func (t *AppBuilderClientStreamIterator) Next() (*AppBuilderClientAnswer, error) {
	data, err := t.r.ReadMessageLine()
	if err != nil && !(err == io.EOF) {
		t.body.Close()
		return nil, fmt.Errorf("requestID=%s, err=%v", t.requestID, err)
	}
	if err != nil && err == io.EOF {
		t.body.Close()
		return nil, err
	}
	if strings.HasPrefix(string(data), "data:") {
		var resp AppBuilderClientRawResponse
		if err := json.Unmarshal(data[5:], &resp); err != nil {
			t.body.Close()
			return nil, fmt.Errorf("requestID=%s, err=%v", t.requestID, err)
		}
		answer := &AppBuilderClientAnswer{}
		answer.transform(&resp)
		return answer, nil
	}
	// 非SSE格式关闭连接，并返回数据
	t.body.Close()
	return nil, fmt.Errorf("requestID=%s, body=%s", t.requestID, string(data))
}

// AppBuilderClientOnceIterator 非流式返回时对应的迭代器，只可迭代一次
type AppBuilderClientOnceIterator struct {
	body      io.ReadCloser
	requestID string
}

func (t *AppBuilderClientOnceIterator) Next() (*AppBuilderClientAnswer, error) {
	data, err := io.ReadAll(t.body)
	if err != nil {
		return nil, fmt.Errorf("requestID=%s, err=%v", t.requestID, err)
	}
	defer t.body.Close()
	var resp AppBuilderClientRawResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("requestID=%s, err=%v", t.requestID, err)
	}
	answer := &AppBuilderClientAnswer{}
	answer.transform(&resp)
	return answer, nil
}
