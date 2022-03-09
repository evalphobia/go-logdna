package logdna

import (
	"fmt"
	"time"

	"github.com/evalphobia/httpwrapper/request"
)

// Client handles logging to LogDNA.
type Client struct {
	Config
	daemon *Daemon
}

// New creates an initialized *Client.
func New(conf Config) (*Client, error) {
	if err := conf.Init(); err != nil {
		return nil, err
	}

	cli := &Client{
		Config: conf,
	}

	if !conf.Sync {
		cli.RunDaemon(conf.getCheckpointSize(), conf.getCheckpointInterval())
	}
	return cli, nil
}

// Debug sends a debug level log.
func (c *Client) Debug(line string, opt ...map[string]interface{}) error {
	return c.EmitWithLevel(LogLevelDebug, line, opt...)
}

// Trace sends a trace level log.
func (c *Client) Trace(line string, opt ...map[string]interface{}) error {
	return c.EmitWithLevel(LogLevelTrace, line, opt...)
}

// Info sends a info level log.
func (c *Client) Info(line string, opt ...map[string]interface{}) error {
	return c.EmitWithLevel(LogLevelInfo, line, opt...)
}

// Warn sends a warning level log.
func (c *Client) Warn(line string, opt ...map[string]interface{}) error {
	return c.EmitWithLevel(LogLevelWarn, line, opt...)
}

// Err sends a error level log.
func (c *Client) Err(line string, opt ...map[string]interface{}) error {
	return c.EmitWithLevel(LogLevelError, line, opt...)
}

// Fatal sends a fatal level log.
func (c *Client) Fatal(line string, opt ...map[string]interface{}) error {
	return c.EmitWithLevel(LogLevelFatal, line, opt...)
}

// Emit sends a log of a given level.
func (c *Client) EmitWithLevel(level, line string, opt ...map[string]interface{}) error {
	if !isMoreLevel(c.MinimumLevel, level) {
		return nil
	}

	d := LogData{
		Message: line,
		Level:   level,
		App:     c.App,
		Env:     c.Env,
	}
	if len(opt) != 0 {
		d.Meta = opt[0]
	}
	return c.Emit(d)
}

func (c *Client) Emit(d LogData) error {
	if c.Sync {
		return c.send([]*logPayload{d.toPayload()})
	}

	c.daemon.Add(d.toPayload())
	return nil
}

// RunDaemon runs a Daemon in background.
func (c *Client) RunDaemon(size int, interval time.Duration) {
	c.daemon = NewDaemon(size, interval, c.send)
	c.daemon.Run()
}

func (c *Client) send(logs []*logPayload) error {
	if len(logs) == 0 {
		return nil
	}

	jsonData, err := logsToJSON(logs)
	if err != nil {
		return err
	}

	return c.callAPI(jsonData)
}

// callAPI sends a POST request to endpoint.
func (c *Client) callAPI(params interface{}) error {
	conf := c.Config
	now := time.Now().UnixNano() / int64(time.Millisecond)
	url := fmt.Sprintf("%s?hostname=%s&mac=%s&ip=%s&now=%d", conf.endpoint, conf.hostname, conf.macaddr, conf.ip, now)

	resp, err := request.POST(url, request.Option{
		Payload:     params,
		PayloadType: request.PayloadTypeJSON,
		User:        conf.apikey,
		Retry:       !conf.NoRetry,
		Debug:       conf.Debug,
		UserAgent:   "go-logdna/v0.0.1",
		Timeout:     conf.timeout,
	})
	if err != nil {
		return err
	}
	return resp.Close()
}
