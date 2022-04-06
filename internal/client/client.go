package client

import (
	"errors"
	"fmt"
	"net"
	guidecontracts "word-of-wisom/api/guide"
	srvcontracts "word-of-wisom/api/server"
	"word-of-wisom/pkg/gtp"

	"github.com/sirupsen/logrus"
	"github.com/tinylib/msgp/msgp"
)

var errServiceRestricted = errors.New("service restricted error")

type Client struct {
	server string
	guides []string
}

func NewClient(server string, guides []string) *Client {
	return &Client{
		server: server,
		guides: guides,
	}
}

func (c *Client) RequestQuote(retryCount int) (string, error) {
	for i := 0; i < retryCount; i++ {
		response, err := c.initialRequest()
		if err != nil {
			return "", fmt.Errorf("initial request: %w", err)
		}

		switch srvcontracts.ResponseType(response.Type) {
		case srvcontracts.ServiceGranted:
			serviceGrantedMsg := srvcontracts.ServiceGrantedPayload{}
			if _, err = serviceGrantedMsg.UnmarshalMsg(response.Payload); err != nil {
				return "", fmt.Errorf("unmarshal service granted response: %w", err)
			}

			logrus.Debug("service granted msg")
			logrus.Debugf("%v+", serviceGrantedMsg)

			return serviceGrantedMsg.Quote, nil
		case srvcontracts.ServiceRestricted:
			serviceRestrictedMsg := srvcontracts.ServiceRestrictedPayload{}
			if _, err = serviceRestrictedMsg.UnmarshalMsg(response.Payload); err != nil {
				return "", fmt.Errorf("unmarshal service restricted response: %w", err)
			}

			logrus.Debug("service restricted msg")
			logrus.Debugf("%v+", serviceRestrictedMsg)

			lastHash, err := c.guidedTourRequest(serviceRestrictedMsg)
			if err != nil {
				return "", fmt.Errorf("guided tour request: %w", err)
			}

			quote, err := c.tourCompleteRequest(serviceRestrictedMsg.InitialHash, lastHash)

			if err != nil {
				if errors.Is(err, errServiceRestricted) {
					logrus.Debug("service restricted again")
					continue
				}

				return "", fmt.Errorf("tour complete request: %w", err)
			}

			return quote, nil

		case srvcontracts.UnsupportedRequest:
			return "", fmt.Errorf("unsupported request: %w", err)
		default:
			return "", fmt.Errorf("unknown response type: %v", response.Type)
		}
	}

	return "", fmt.Errorf("max retries")
}

func (c *Client) initialRequest() (srvcontracts.ResponseMsg, error) {
	request := srvcontracts.RequestMsg{Type: byte(srvcontracts.InitialRequest)}
	response := srvcontracts.ResponseMsg{}

	conn, err := net.Dial("tcp", c.server)
	if err != nil {
		return response, fmt.Errorf("connect to server: %w", err)
	}

	writer := msgp.NewWriter(conn)

	err = request.EncodeMsg(writer)
	if err != nil {
		return response, fmt.Errorf("encode msg: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return response, fmt.Errorf("flush writer: %w", err)
	}

	logrus.Debugf("Server request sent: %v", request)

	if err := response.DecodeMsg(msgp.NewReader(conn)); err != nil {
		return response, fmt.Errorf("decode server response: %w", err)
	}

	return response, nil
}

func (c *Client) guidedTourRequest(serviceRestrictedMsg srvcontracts.ServiceRestrictedPayload) ([32]byte, error) {
	prevHash := serviceRestrictedMsg.InitialHash

	for i := 1; i < int(serviceRestrictedMsg.TourLength)+1; i++ {
		logrus.Debugf("tour number: %v", i)

		request := guidecontracts.RequestMsg{
			PreviousHash: prevHash,
			TourNumber:   byte(i),
			TourLength:   serviceRestrictedMsg.TourLength,
		}
		guide := c.guides[gtp.GuideIndex(prevHash, len(c.guides))]

		conn, err := net.Dial("tcp", guide)
		if err != nil {
			return prevHash, fmt.Errorf("connect to guide: %w", err)
		}

		writer := msgp.NewWriter(conn)

		err = request.EncodeMsg(writer)
		if err != nil {
			return prevHash, fmt.Errorf("encode msg: %w", err)
		}

		if err := writer.Flush(); err != nil {
			return prevHash, fmt.Errorf("flush writer: %w", err)
		}

		logrus.Debugf("Guide request sent: %v", request)

		response := guidecontracts.ResponseMsg{}
		if err := response.DecodeMsg(msgp.NewReader(conn)); err != nil {
			return prevHash, fmt.Errorf("decode guide response: %w", err)
		}

		prevHash = response.Hash
	}

	return prevHash, nil
}

func (c *Client) tourCompleteRequest(initialHash, lastHash [32]byte) (string, error) {
	tourCompletePayload := srvcontracts.TourCompletePayload{
		InitialHash: initialHash, LastHash: lastHash}

	requestPayload, err := tourCompletePayload.MarshalMsg(nil)
	if err != nil {
		logrus.Debugf("%v+", tourCompletePayload)
		return "", fmt.Errorf("marshal tour complete payload: %w", err)
	}

	request := srvcontracts.RequestMsg{Type: byte(srvcontracts.TourCompleteRequest), Payload: requestPayload}

	conn, err := net.Dial("tcp", c.server)
	if err != nil {
		return "", fmt.Errorf("connect to server: %w", err)
	}

	writer := msgp.NewWriter(conn)

	err = request.EncodeMsg(writer)
	if err != nil {
		return "", fmt.Errorf("encode msg: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("flush writer: %w", err)
	}

	logrus.Debugf("Server request sent: %v", request)

	response := srvcontracts.ResponseMsg{}
	if err := response.DecodeMsg(msgp.NewReader(conn)); err != nil {
		return "", fmt.Errorf("decode server response: %w", err)
	}

	switch srvcontracts.ResponseType(response.Type) {
	case srvcontracts.ServiceGranted:
		serviceGrantedMsg := srvcontracts.ServiceGrantedPayload{}
		if _, err := serviceGrantedMsg.UnmarshalMsg(response.Payload); err != nil {
			return "", fmt.Errorf("unmarshal service granted response: %w", err)
		}

		return serviceGrantedMsg.Quote, nil
	case srvcontracts.ServiceRestricted:
		return "", errServiceRestricted
	case srvcontracts.UnsupportedRequest:
		return "", fmt.Errorf("unsupported request")
	default:
		return "", fmt.Errorf("unknown response type: %v", response.Type)
	}
}
