// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"io"
	"path/filepath"

	"github.com/disintegration/imaging"
	agent "github.com/saferwall/saferwall/internal/agent"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/utils"
	"github.com/saferwall/saferwall/services/config"
)

const (
	defaultGRPCPort        = ":50051"
	defaultFileScanTimeout = 30
	defaultVPNCountry      = "USA"
	defaultOS              = "Windows 7 64-bit"
)

// DetonationResults represents the results for a detonation.
type DetonationResult struct {
	// The API trace results. This consists of a list of all API calls made by
	// the sample.
	APITrace []byte `json:"api_trace,omitempty"`
	// The logs produced by the agent running inside the VM.
	AgentLog []byte `json:"agent_log,omitempty"`
	// The logs produced by the sandbox.
	SandboxLog []byte `json:"sandbox_log,omitempty"`
	// The config used to scan dynamically the sample.
	ScanCfg config.DynFileScanCfg `json:"scan_config,omitempty"`
	// Environment represents the environment used to scan the file.
	Environment `json:"env,omitempty"`
}

// Environment represents the environment used to scan the file.
// This include OS versions, installed software, analyzer version.
type Environment struct {
	// The sandbox version.
	SandboxVersion string `json:"sandbox_version,omitempty"`
	// The agent version.
	AgentVersion string `json:"agent_version,omitempty"`
}

func (s *Service) detonate(logger log.Logger, vm *VM,
	cfg config.FileScanCfg) (DetonationResult, error) {

	detRes := DetonationResult{}
	ctx := context.Background()

	// Establish a gRPC connection to the agent server running
	// inside the guest.
	client, err := agent.New(vm.IP + defaultGRPCPort)
	if err != nil {
		logger.Errorf("failed to establish connection to server: %v", err)
		return detRes, err
	}

	// Deploy the sandbox component files inside the guest.
	ver, err := client.Deploy(ctx, s.cfg.Agent.AgentDestDir, s.sandbox)
	if err != nil {
		return detRes, err
	}
	logger.Infof("sandbox version %s has been deployed", ver)

	detRes.SandboxVersion = ver
	detRes.AgentVersion = vm.AgentVersion

	src := filepath.Join(s.cfg.SharedVolume, cfg.SHA256)
	sampleContent, err := utils.ReadAll(src)
	if err != nil {
		return detRes, err
	}

	// Analyze the sample. This call will block until results
	// are ready.
	scanCfg := toJSON(cfg.DynFileScanCfg)
	res, err := client.Analyze(ctx, scanCfg, sampleContent)
	if err != nil {
		return detRes, err
	}

	var traceLog []interface{}
	err = Decode(bytes.NewReader(res.TraceLog), &traceLog)
	if err != nil {
		return detRes, err
	}

	var sandboxLog []interface{}
	err = Decode(bytes.NewReader(res.SandboxLog), &sandboxLog)
	if err != nil {
		return detRes, err
	}

	var agentLog []interface{}
	err = Decode(bytes.NewReader(res.AgentLog), &agentLog)
	if err != nil {
		return detRes, err
	}

	detRes.ScanCfg = cfg.DynFileScanCfg
	detRes.APITrace = toJSON(traceLog)
	detRes.AgentLog = toJSON(agentLog)
	detRes.SandboxLog = toJSON(sandboxLog)

	return detRes, nil
}

// generate thumbnails for the sandbox desktop screenshots.
func (s *Service) generateThumbnail(r io.Reader) (io.Writer, error) {

	buf := new(bytes.Buffer)

	// load images and make 100x100 thumbnails of them
	img, err := imaging.Decode(r, nil)
	if err != nil {
		return nil, err
	}

	x := 730
	y := 450
	thumbnail := imaging.Thumbnail(img, x, y, imaging.CatmullRom)

	// create a new blank image
	dst := imaging.New(x, y, color.NRGBA{0, 0, 0, 0})

	// paste thumbnails into the new image side by side
	dst = imaging.Paste(dst, thumbnail, image.Pt(0, 0))

	// write the combined image to an io writer.
	opts := imaging.JPEGQuality(80)
	err = imaging.Encode(buf, dst, imaging.JPEG, opts)
	if err != nil {
		return nil, err

	}

	return buf, nil
}
