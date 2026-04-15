package services

import (
	"fmt"
	"net/url"
	"strings"

	"gorm.io/gorm"

	"simon/apps/api/internal/preview"
	"simon/apps/api/internal/repositories"
)

type ResolvedHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RequestPreviewDTO struct {
	Method           string           `json:"method"`
	ResolvedURL      string           `json:"resolved_url"`
	ResolvedHeaders  []ResolvedHeader `json:"resolved_headers"`
	ResolvedBody     string           `json:"resolved_body"`
	Curl             string           `json:"curl"`
	MissingVariables []string         `json:"missing_variables"`
}

type RequestPreviewService struct {
	requests *repositories.RequestRepository
	headers  *repositories.RequestHeaderRepository
	params   *repositories.RequestQueryParamRepository
	envs     *repositories.EnvironmentRepository
	envVars  *repositories.EnvironmentVariableRepository
}

func NewRequestPreviewService(db *gorm.DB) *RequestPreviewService {
	return &RequestPreviewService{
		requests: repositories.NewRequestRepository(db),
		headers:  repositories.NewRequestHeaderRepository(db),
		params:   repositories.NewRequestQueryParamRepository(db),
		envs:     repositories.NewEnvironmentRepository(db),
		envVars:  repositories.NewEnvironmentVariableRepository(db),
	}
}

func (s *RequestPreviewService) Build(userID, requestID, environmentID uint) (*RequestPreviewDTO, error) {
	req, err := s.requests.FindByIDForUser(requestID, userID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, ErrNotFound
	}

	env, err := s.envs.FindByIDForUser(environmentID, userID)
	if err != nil {
		return nil, err
	}
	if env == nil {
		return nil, ErrNotFound
	}

	vars := map[string]string{}
	varList, err := s.envVars.ListByEnvironmentID(environmentID)
	if err != nil {
		return nil, err
	}
	for _, v := range varList {
		if !v.Enabled {
			continue
		}
		vars[v.Key] = v.Value
	}

	var missing [][]string

	resolvedURL, mURL := preview.Resolve(req.URL, vars)
	missing = append(missing, mURL)

	headerRows, err := s.headers.ListByRequestID(requestID)
	if err != nil {
		return nil, err
	}
	resolvedHeaders := make([]ResolvedHeader, 0)
	headerLines := make([]preview.HeaderLine, 0)
	for _, h := range headerRows {
		if !h.Enabled {
			continue
		}
		rk, mk := preview.Resolve(h.Key, vars)
		rv, mv := preview.Resolve(h.Value, vars)
		missing = append(missing, mk, mv)
		resolvedHeaders = append(resolvedHeaders, ResolvedHeader{Key: rk, Value: rv})
		headerLines = append(headerLines, preview.HeaderLine{Key: rk, Value: rv})
	}

	paramRows, err := s.params.ListByRequestID(requestID)
	if err != nil {
		return nil, err
	}
	var enabledParams []queryPair
	for _, p := range paramRows {
		if !p.Enabled {
			continue
		}
		rk, mk := preview.Resolve(p.Key, vars)
		rv, mv := preview.Resolve(p.Value, vars)
		missing = append(missing, mk, mv)
		enabledParams = append(enabledParams, queryPair{key: rk, val: rv})
	}

	finalURL, err := mergeURLQuery(resolvedURL, enabledParams)
	if err != nil {
		return nil, fmt.Errorf("url após resolução: %w", err)
	}

	resolvedBody, mBody := preview.Resolve(req.Body, vars)
	missing = append(missing, mBody)

	method := strings.ToUpper(strings.TrimSpace(req.Method))
	curlStr := preview.BuildCurl(method, finalURL, headerLines, resolvedBody)

	merged := preview.MergeMissing(missing...)

	return &RequestPreviewDTO{
		Method:           method,
		ResolvedURL:      finalURL,
		ResolvedHeaders:  resolvedHeaders,
		ResolvedBody:     resolvedBody,
		Curl:             curlStr,
		MissingVariables: merged,
	}, nil
}

type queryPair struct {
	key, val string
}

func mergeURLQuery(base string, params []queryPair) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	var parts []string
	if u.RawQuery != "" {
		parts = append(parts, u.RawQuery)
	}
	for _, p := range params {
		parts = append(parts, url.QueryEscape(p.key)+"="+url.QueryEscape(p.val))
	}
	if len(parts) > 0 {
		u.RawQuery = strings.Join(parts, "&")
	}
	return u.String(), nil
}
